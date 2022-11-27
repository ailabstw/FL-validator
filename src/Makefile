PROJ_ROOT=$(realpath ..)

all: image

image:
	docker run --rm \
		-v $(PROJ_ROOT):/validator \
		-w /validator/src \
		golang:1.15.2 go build -o validator
	docker build --tag federated/validator .