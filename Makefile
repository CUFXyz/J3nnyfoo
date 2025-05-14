run: 
	go run ./cmd/main.go 

swag: 
	swag init -g ./cmd/main.go --parseInternal -o ./docs

docker:
	docker build -t j3nnyfoo .
	docker run -p 9090:9090 j3nnyfoo