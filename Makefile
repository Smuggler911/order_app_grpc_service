
gen:
	@protoc -I ./api \
         --go_out=./pkg --go_opt=paths=source_relative \
         --go-grpc_out=./pkg  --go-grpc_opt paths=source_relative\
		 --grpc-gateway_out ./pkg \
         --grpc-gateway_opt paths=source_relative \
         --grpc-gateway_opt grpc_api_configuration=config/orderService.yaml \
         --grpc-gateway_opt standalone=true \
         ./api/orders/orders.proto

migrate:
	@migrate -path ./migration/ -database "postgres://postgres:coffice@213.226.127.170/orders?sslmode=disable" up

drop:
	@migrate -path ./migrations/ -database "postgres://postgres:coffice@213.226.127.170/orders?sslmode=disable" down
