start:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o input/input-linux-amd64 input/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o stats/stats-linux-amd64 stats/main.go
	docker-compose up -d
stop:
	docker-compose down --volumes
	rm -rf input/input-linux-amd64 stats/stats-linux-amd64