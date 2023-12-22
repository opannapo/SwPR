.PHONY: clean all init generate generate_mocks compose-up compose-down

all: build/main

build/main: cmd/main.go generated
	@echo "Building..."
	go build -o $@ $<

clean:
	rm -rf generated

init: generate
	go mod tidy
	go mod vendor

test:
	go test -short -coverprofile coverage.out -v ./...

generate: generated generate_mocks generate_mocks_util

generated: api.yml
	@echo "Generating files..."
	mkdir generated || true
	oapi-codegen --package generated -generate types,server,spec $< > generated/api.gen.go

INTERFACES_GO_FILES := $(shell find repository -name "interfaces.go")
INTERFACES_GEN_GO_FILES := $(INTERFACES_GO_FILES:%.go=%.mock.gen.go)

#generate_mocks: $(INTERFACES_GEN_GO_FILES)
#$(INTERFACES_GEN_GO_FILES): %.mock.gen.go: %.go
#	@echo "Generating mocks $@ for $<"
#	mockgen -source=$< -destination=$@ -package=$(shell basename $(dir $<))

generate_mocks:
	mockgen -source="repository/interfaces.go" -destination="mock/repository/interfaces.go" -package="repository_mock"
	mockgen -source="util/password.go" -destination="mock/util/password.go" -package="util_mock"

compose-up:
	sudo docker-compose up --build -d

compose-down:
	sudo docker-compose down --volumes