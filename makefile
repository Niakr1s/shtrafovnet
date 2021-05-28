proto:
	buf generate
.PHONY: proto

docker:
	CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o server main.go
	docker build -t shtrafov .
	rm server

run:
	docker run -p 8081:8081 -p 9000:9000 shtrafov