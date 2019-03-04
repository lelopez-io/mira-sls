build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/list-by-genre	list-by-genre/main.go
