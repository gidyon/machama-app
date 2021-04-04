API_IN_PATH := api/proto
API_OUT_PATH := pkg/api
OPEN_API_V2_OUT_PATH := api/openapiv2

setup_dev: ## Sets up a development environment for the emrs project
	@cd deployments/compose/dev &&\
	docker-compose up -d

setup_redis: ## Starts redis server
	@cd deployments/compose/dev &&\
	docker-compose up -d redis

teardown_dev: ## Tear down development environment for the emrs project
	@cd deployments/compose/dev &&\
	docker-compose down

protoc_chama: ## Compiles chama.proto
	@protoc -I=$(API_IN_PATH) -I=third_party --go-grpc_out=$(API_OUT_PATH)/chama --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/chama --go_opt=paths=source_relative chama.proto
	@protoc -I=$(API_IN_PATH) -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/chama chama.proto
	@protoc -I=$(API_IN_PATH) -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) chama.proto

protoc_loan: ## Compiles chama.proto
	@protoc -I=$(API_IN_PATH) -I=third_party --go-grpc_out=$(API_OUT_PATH)/loan --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/loan --go_opt=paths=source_relative loan.proto
	@protoc -I=$(API_IN_PATH) -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/loan loan.proto
	@protoc -I=$(API_IN_PATH) -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) loan.proto

protoc_transaction: ## Compiles chama.proto
	@protoc -I=$(API_IN_PATH) -I=third_party --go-grpc_out=$(API_OUT_PATH)/transaction --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/transaction --go_opt=paths=source_relative transaction.proto
	@protoc -I=$(API_IN_PATH) -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/transaction transaction.proto
	@protoc -I=$(API_IN_PATH) -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) transaction.proto

copy_documentation: ## Updates documentation
	@cp -r $(OPEN_API_V2_OUT_PATH) cmd/apidoc/dist

protoc_all: protoc_chama protoc_loan protoc_transaction copy_documentation

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
