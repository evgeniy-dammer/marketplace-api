.PHONY: protos

protos: # execute from emenu-api directory
	protoc -I pkg/protos/ \
			--go_out=pkg/protos \
			--go_opt=paths=source_relative \
			--go-grpc_out=pkg/protos \
			--go-grpc_opt=paths=source_relative \
			pkg/protos/item.proto

test:
	grpcurl --plaintext 0.0.0.0:1111 protos.ItemService.FindAll

build: clean
	go mod download
	go build -o emenu-api cmd/main.go

clean:
	rm -f emenu-api

run:
	go run cmd/main.go

image: rmimage
	docker build -t emenu-api .

rmimage:
	docker image rm -f emenu-api

cont:
	docker run -dit --rm -p 1111:1111 --name emenu-api emenu-api

stopcont:
	docker stop emenu-api

prune:
	docker image prune