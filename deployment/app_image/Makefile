BASE_IMAGE = "ghcr.io/ypodlesov/todo_list_service/base_image:latest"

.PHONY: auth
auth:
	echo $(GITHUB_TOKEN) | docker login ghcr.io -u ypodlesov --password-stdin

.PHONY: docker-build
docker-build: auth
	PG_PORT=$(shell expr $(PG_PORT_STR) + 0);

	docker build --network=host \
		--build-arg BASE_IMAGE=$(BASE_IMAGE) \
		--build-arg CONFIG_PATH=$(CONFIG_PATH) \
		--build-arg PG_HOST=$(PG_HOST) \
		--build-arg PG_USER=$(PG_USER) \
		--build-arg PG_PASSWORD=$(PG_PASSWORD) \
		--build-arg PG_DB_NAME=$(PG_DB_NAME) \
		--build-arg PG_MIGRATIONS_DIR=$(PG_MIGRATIONS_DIR) \
		-t $(APP_IMAGE) \
		-f $(PROJECT_PATH)/deployment/app_image/Dockerfile $(PROJECT_PATH)

.PHONY: docker-push
docker-push: docker-build
	docker push $(APP_IMAGE)
