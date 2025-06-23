.PHONY: install build

install: build
	go install


build:
	templ generate
	go build .