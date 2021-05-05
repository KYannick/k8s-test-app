registry = 
image = testapp
tag = latest

default: build

build-go:
	GOOS=linux go build -o ./simple-go-http

build-docker:
	docker build -t ${registry}/${image}:${tag} .

build: build-go build-docker clean

clean:
	rm -f simple-go-http