.PHONY: build all clean in-memory postgres

all:	build
		docker-compose -f deploy/docker-compose.yml up

build:
		docker-compose -f deploy/docker-compose.yml build

in-memory:	build
		docker-compose -f deploy/docker-compose.yml up

postgres: build
		STORAGE=-postgres docker-compose -f deploy/docker-compose.yml --profile db up

clean:
		rm $(BINARY_NAME)

