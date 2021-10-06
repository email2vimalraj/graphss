run: gen
	@echo "running"
	go run pkg/cmd/graphss/main.go --config ./conf/config.toml

build: gen
	@echo "building a binary"
	go build -o bin/graphss pkg/cmd/graphss/main.go

gen:
	@echo "generating wire files"
	wire ./pkg/http