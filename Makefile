all: build copy-config run

build:
	go build -o bin/s0

copy-config:
	cp config.sample.yaml config.yaml

run:
	bin/s0
