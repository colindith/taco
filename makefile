dev:
		go run httpd/main.go
build:
		go build -o ./app httpd/main.go
run:
		go build -o ./app httpd/main.go
		./app
