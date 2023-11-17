.PHONY:build
build:
	rm -rf build/*
	go build -o build/mahjong-server cmd/main.go

run:
	go run cmd/main.go

