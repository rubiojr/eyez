BIN=eyez

bin:
	mkdir -p bin
	go build -o bin/$(BIN) ./cmd/$(BIN)

.PHONY: bin
