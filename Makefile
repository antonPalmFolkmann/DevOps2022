prod:
	./bootstrap.sh

destroy-prod:
	./src/scripts/terraform/destroy.sh

build:
	./src/scripts/docker/build.sh

build-and-push:
	./src/scripts/docker/build-and-push.sh

start:
	docker-compose down && docker-compose up --build
