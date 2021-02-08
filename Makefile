BINARY := "log-reloading"

default: $(BINARY)

.PHONY:
clean:
	@git clean -Xf
	@go mod tidy
	@go clean

$(BINARY): main.go types.go
	@go build -o $(BINARY)

.PHONY:
run: $(BINARY)
	@./log-reloading

all: clean run