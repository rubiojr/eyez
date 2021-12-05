package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rjeczalik/notify"
	"github.com/rubiojr/eyez/internal/db"
)

const MAX_RECORDS = "50"

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
	stack = NewRStack()
)

var m model

type Record struct {
	db.Record
}

func (i Record) Length() int {
	return len(i.Body)
}

func (i Record) Title() string {
	//elapsed := i.DateEnd.UnixMilli() - i.DateStart.UnixMilli()
	elapsed := 0
	return fmt.Sprintf("%s %s %d %db %dms", i.Method, i.Path, i.Status, i.Length, elapsed)
}
func (i Record) Description() string {
	d := i.DateEnd.Format("2006-01-02 15:04:05")
	return fmt.Sprintf("%s connect core", d)
}
func (i Record) FilterValue() string { return i.URL }

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

type model struct {
	list         list.Model
	keys         *listKeyMap
	delegateKeys *delegateKeyMap
}

var records list.Model

func monitorDatabase() error {
	for {
		c := make(chan notify.EventInfo, 1)

		if err := notify.Watch(db.DefaultDatabase, c, notify.InModify); err != nil {
			return err
		}
		defer notify.Stop(c)

		// Block until an event is received.
		switch ei := <-c; ei.Event() {
		case notify.InModify:
			fmt.Println("Database changed")
			fetchItems()
		}
	}
}

func fetchItems() []list.Item {
	items := []list.Item{}
	err := db.Each(func(rec *db.Record) error {
		r := Record{*rec}
		items = append(items, r)
		stack.Push(&r)
		return nil
	})
	if err != nil {
		panic(err)
	}

	return items
}

func NewModel() model {
	var (
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
	)

	go func() {
		err := monitorDatabase()
		if err != nil {
			panic(err)
		}
	}()

	// Setup list
	delegate := newItemDelegate(delegateKeys)
	records = list.NewModel(fetchItems(), delegate, 0, 0)
	records.Title = "Outboud Connections"
	records.Styles.Title = titleStyle
	records.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.insertItem,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	return model{
		list:         records,
		keys:         listKeys,
		delegateKeys: delegateKeys,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, tick(), m.list.StartSpinner())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		topGap, rightGap, bottomGap, leftGap := appStyle.GetPadding()
		m.list.SetSize(msg.Width-leftGap-rightGap, msg.Height-topGap-bottomGap)

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.toggleSpinner):
			cmd := m.list.ToggleSpinner()
			return m, cmd

		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil

		case key.Matches(msg, m.keys.insertItem):
			statusCmd := m.list.NewStatusMessage(statusMessageStyle("Refreshed"))
			return m, tea.Batch(statusCmd)
		}
	case TickMsg:
		records := []list.Item{}
		stack.Each(func(r *Record) {
			records = append(records, r)
		})
		cmds := []tea.Cmd{tick()}
		if len(records) > 0 {
			cmds = append(cmds, m.list.SetItems(records))
		}
		return m, tea.Batch(cmds...)
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return appStyle.Render(m.list.View())
}

type TickMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(1*time.Second, func(_ time.Time) tea.Msg {
		return TickMsg{}
	})
}
