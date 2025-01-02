IMAGE = ghcr.io/ypodlesov/todo_list_service/base_image:latest
BASE_IMAGE_DIR = "ci_cd/base_image"

.PHONY: auth
auth:
	echo $(GITHUB_TOKEN) | docker login ghcr.io -u ypodlesov --password-stdin

.PHONY: docker-build
docker-build: auth
	docker build --network=host \
	--build-arg BASE_IMAGE_DIR=$(BASE_IMAGE_DIR) \
	-t $(IMAGE) -f Dockerfile ../..

.PHONY: docker-push
docker-push: docker-build
	docker push $(IMAGE)