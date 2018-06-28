deps:
	dep ensure

test:
	go test -v -bench=.

all: deps
	go build -v -o ts-cli

clean:
	rm -f ts-cli

