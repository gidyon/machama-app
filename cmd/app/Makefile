SERVICE := machama-app

JWT_KEY := hDI0eBv11TbuboZ01qpnOuYRYLh6gQUOQhC9Mfagzv9l3gJso7CalTt7wGzJCVwbeDIfOX6fwS79pnisW7udhQ
API_BLOCK_KEY := 9AI8o4ta02gdqWsVhYe0r276z7my6yDwY78rCsrcofT7pCNq4WwnRoW93hn8WFJM0HheZHDYPc4tD+tUXVNEGw
API_HASH_KEY := 73H/I3+27Qp3ZETqMzbYa/EGT826Zxx2821cmHUl7fTX/DmkIWPJatczkxN3p8RHbdOGWT/HDRAf7gqhZcZOow

local_image := machama-app
image := gidyon/machama-app
context := .

ifdef IMAGE
	image=$(IMAGE)
else
	imagex := $(image)
	image_local := $(local_image)
	ifdef tag
		image=$(imagex):$(tag)
		local_image=$(image_local):$(tag)
	else
		image=$(imagex):latest
		local_image=$(image_local):latest
	endif
endif

ifdef BUILD_CONTEXT
	context=$(BUILD_CONTEXT)
endif

run:
	go build -v -o service && \
	SMS_JWT=avc \
	./service -config-file=./config.yaml

compile:
	@GOOS=linux CGO_ENABLED=0 go build -tags netgo -installsuffix netgo -v -o service .

docker_build:
	docker build -t $(local_image) .

docker_tag:
	@docker tag $(local_image) $(image)

docker_push_skaffold:
ifeq ($(PUSH_IMAGE), true)
	@docker push $(image)
endif

commit_master:
	git push origin master

docker_push:
	@docker push $(image)

build_service: compile docker_build docker_tag docker_push

build_service_skaffold: compile docker_build docker_tag docker_push_skaffold


