.PHONY: bootstrap

build/yml2json: *.go
	go build -o build/yml2json .

build/yml2json-linux: *.go
	docker run --rm --volume "$(PWD):/go/src/yml2json" --workdir "/go/src/yml2json" golang:1.6-alpine sh -c "go build -o build/yml2json-linux ."

bootstrap:
	git submodule update --init --recursive

clean:
	@rm -r build
