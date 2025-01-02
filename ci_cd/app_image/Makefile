BASE_IMAGE = "ghcr.io/ypodlesov/todo_list_service/base_image:latest"
APP_IMAGE_DIR = "ci_cd/app_image"

.PHONY: auth
auth:
	echo $(GITHUB_TOKEN) | docker login ghcr.io -u ypodlesov --password-stdin

.PHONY: docker-build
docker-build: auth
	docker build --network=host \
		--build-arg APP_IMAGE_DIR=$(APP_IMAGE_DIR) \
		--build-arg BASE_IMAGE=$(BASE_IMAGE) \
		--build-arg TELEGRAM_API_TOKEN=$(TELEGRAM_API_TOKEN) \
		-t $(APP_IMAGE) \
		-f Dockerfile $(PROJECT_PATH)

.PHONY: docker-push
docker-push: docker-build
	docker push $(APP_IMAGE)
