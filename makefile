# Initial step is set this project name.
APPLICATION_NAME=myfirst_microservice_project
# Step 2: Example: 'wms', short name for this application normally for route name use.
APPLICATION_SHORT_NAME=mfmp
#
.PHONY: init
init:
	sudo rm -r .git && mkdir api api/$(APPLICATION_SHORT_NAME) api/$(APPLICATION_SHORT_NAME)/v1 api/$(APPLICATION_SHORT_NAME)/v1/proto \
	cmd cmd/$(APPLICATION_NAME) \
	internal internal/biz internal/client internal/config \
	internal/db internal/db/query internal/db/sqlc \
	internal/server internal/service \
	migration && touch .gitignore app.env dockerfile sqlc.yaml wait-for.sh \
	cmd/$(APPLICATION_NAME)/main.go api/$(APPLICATION_SHORT_NAME)/v1/proto/service_$(APPLICATION_SHORT_NAME).proto \
	&& migrate create -ext sql -dir ./migration -seq initial \
	&& go mod init && go mod tidy

.PHONY: go_clean_cache
go_clean_cache:
	go clean -modcache

.PHONY: new_migration
new_migration:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir ./migration -seq $$name

.PHONY: sqlc
sqlc:
	sqlc generate

.PHONY: proto
proto:
	rm -f api/$(APPLICATION_SHORT_NAME)/v1/*.go
	protoc --proto_path=api/$(APPLICATION_SHORT_NAME)/v1/proto --proto_path=./third_party \
		--go_out=api/$(APPLICATION_SHORT_NAME)/v1 \
		--go_opt=paths=source_relative \
		--go-grpc_out=api/$(APPLICATION_SHORT_NAME)/v1 \
		--go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=api/$(APPLICATION_SHORT_NAME)/v1 \
		--grpc-gateway_opt=paths=source_relative \
		--grpc-gateway_opt=allow_delete_body=true \
		--grpc-gateway_opt=allow_delete_body=true \
		api/$(APPLICATION_SHORT_NAME)/v1/proto/*.proto

.PHONY: multiplatform_docker_build
TAG ?=  $(shell bash -c 'read -p "Enter the build image tag: " tag; echo $$tag')
multiplatform_docker_build:
	docker build . --platform=linux/amd64,linux/arm64,linux/arm/v7 \
		--tag rksouthasiait/$(APPLICATION_NAME):$(TAG) \
		--build-arg ACCESS_TOKEN_USR=${RKSOUTASIAIT_GH_USERNAME} \
		--build-arg ACCESS_TOKEN_PWD=${RKSOUTASIAIT_GH_TK} \
		--push
		
.PHONY: spin_local_postgres
spin_local_postgres:
	docker run --name my_first_microservice_postgres \
	--restart unless-stopped \
	-p 5432:5432 \
	-e POSTGRES_USER=myfirst \
	-e POSTGRES_PASSWORD=abc123 \
	-v ${HOME}/pg/v16/my_first_proj/data:/var/lib/postgresql/data \
	-d postgres:16.3