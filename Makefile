.PHONY: proto run down

proto:
		protoc --go_out=. --go_opt=paths=source_relative \
            --go-grpc_out=. --go-grpc_opt=paths=source_relative \
            salesservice/proto/sales.proto

run:
	docker-compose up -d

down:
	docker-compose down --remove-orphans
	docker-compose down --volumes