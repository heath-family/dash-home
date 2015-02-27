
default: bin/web

bin/web: src
	go build -o bin/web dash-home.go
