IMAGE_NAME = parking_client

CONTAINER_NAME = parking_client

PORT = 5173

.PHONY: build run clean all stop

build:
	docker build -t $(IMAGE_NAME) .

run:
	docker run -d --name $(CONTAINER_NAME) -p $(PORT):80 $(IMAGE_NAME)

stop:
	docker stop $(CONTAINER_NAME) || true

clean:
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true

initial-start: clean build run

clean-image:
	docker rmi $(IMAGE_NAME) || true

clean-all: clean clean-image
