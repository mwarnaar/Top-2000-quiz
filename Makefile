image:
	docker compose build
push: image
	@echo ${PAT} | docker login ghcr.io --username ${DOCKER_LOGIN} --password-stdin
	docker compose push

