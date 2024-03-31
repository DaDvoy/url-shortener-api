.PHONY: build-http build-grpc all clean up-in-memory up-postgres

all:	build-http up-in-memory

build-http:
		docker-compose -f deploy/docker-compose.yml build

build-grpc:
		PROTOCOL=cmd/url-shortener-grpc/main.go docker-compose -f deploy/docker-compose.yml build

up-in-memory:
		docker-compose -f deploy/docker-compose.yml up

up-postgres:
		STORAGE=-postgres docker-compose -f deploy/docker-compose.yml --profile db up

clean:
		docker stop $$(docker ps -qa); \
		docker rm $$(docker ps -qa); \
		docker rmi -f $$(docker images -qa); \
		docker volume rm $$(docker volume ls -q); \
#		2> /dev/null

