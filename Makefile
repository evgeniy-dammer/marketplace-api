include .env

migrcreate:
	migrate create -ext sql -dir ./schema -seq init

migrup:
	migrate -path ./schema -database 'postgres://emenu:${DB_PASSWORD}@localhost:5432/emenu?sslmode=disable' up

migrdown:
	migrate -path ./schema -database 'postgres://emenu:${DB_PASSWORD}@localhost:5432/emenu?sslmode=disable' down





build: clean
	go mod download
	go build -tags=jsoniter -o emenu-api cmd/main.go

clean:
	rm -f emenu-api

run:
	go run -tags=jsoniter cmd/main.go

lint:
	gofumpt -w . && gci write --skip-generated -s standard,default .


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


protos: # execute from root directory
	protoc -I proto/ \
			--go_out=internal/protos \
			--go_opt=paths=source_relative \
			--go-grpc_out=internal/protos \
			--go-grpc_opt=paths=source_relative \
			protos/item.proto

test:
	grpcurl --plaintext 0.0.0.0:1111 protos.ItemService.FindAll