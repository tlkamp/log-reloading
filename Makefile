default: build

safe-clean:
	@git clean -Xi

clean:
	@git clean -Xf

build:
	@go build

run: clean build
	@./log-reloading