BIN=eyez

bin:
	mkdir -p bin
	go build -o bin/$(BIN) ./cmd/$(BIN)

certinfo:
	openssl x509 -in certs/rootCA.crt -text

certgen:
	openssl req -x509 -newkey rsa:4096 -keyout certs/rootCA.key -out certs/rootCA.crt -sha256 -nodes -subj "/C=WO/ST=World/L=Island/O=eyeZ Ltd/OU=Org/CN=eyez.rubiojr.github.io"

.PHONY: bin