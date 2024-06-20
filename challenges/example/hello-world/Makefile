.PHONY: build
build:
	docker build -t pathwar/level-example .

.PHONY: clean
clean:
	pathwar.land hypervisor prune || true
	docker rm -f pathwar-level-example || true

.PHONY: run-with-pathwar
run-with-pathwar: build
	pathwar.land hypervisor run --web-port=8899 pathwar/level-example

.PHONY: run-with-pathwar
run-with-docker: build
	docker run --name=pathwar-level-example -it --rm -p 8899:80 pathwar/level-example
