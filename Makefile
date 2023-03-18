include .env

migrcreate:
	migrate create -ext sql -dir ./schema -seq init

migrup:
	migrate -path ./schema -database 'postgres://marketplace:${DB_PASSWORD}@localhost:5432/marketplace?sslmode=disable' up

migrdown:
	migrate -path ./schema -database 'postgres://marketplace:${DB_PASSWORD}@localhost:5432/marketplace?sslmode=disable' down

build: clean
	go mod download
	go build -tags=jsoniter -o marketplace-api cmd/app/main.go

clean:
	rm -f marketplace-api

run:
	go run -tags=jsoniter cmd/app/main.go

lint:
	gofumpt -w . && gci write --skip-generated -s standard,default . && golangci-lint run

fields:
	fieldalignment -fix ./internal/usecase/user

swagger:
	swag init --parseDependency --generalInfo ./internal/delivery/http/delivery.go --output ./docs/

json:
	easyjson -all ./internal/domain/category/type.go

imagebuild: imageremove
	docker build -t marketplace-api .

imageremove:
	docker image rm -f marketplace-api

imageprune:
	docker image prune

contrun:
	docker run -dit --rm --net=host --name marketplace-api marketplace-api

contstop:
	docker stop marketplace-api

protos: # execute from root directory
	protoc -I proto/ \
			--go_out=internal/protos \
			--go_opt=paths=source_relative \
			--go-grpc_out=internal/protos \
			--go-grpc_opt=paths=source_relative \
			protos/item.proto

test:
	grpcurl --plaintext 0.0.0.0:1111 protos.ItemService.FindAll