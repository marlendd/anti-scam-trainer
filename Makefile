CONTAINER_RUNTIME := $(shell which docker || which podman)

$(info using ${CONTAINER_RUNTIME})

.PHONY: up down clean test test-e2e lint lint-frontend tools

up: down
	${CONTAINER_RUNTIME} compose up --build -d

down:
	${CONTAINER_RUNTIME} compose down

clean:
	${CONTAINER_RUNTIME} compose down -v

test:
	cd backend && go test -v -race -cover ./...

test-e2e:
	${CONTAINER_RUNTIME} compose down -v --remove-orphans
	${CONTAINER_RUNTIME} volume rm -f anti-scam-trainer_antiscam_data 2>/dev/null || true
	${CONTAINER_RUNTIME} compose up --build -d
	@echo "wait cluster to start..." && sleep 5
	cd backend && RUN_INTEGRATION_TESTS=1 go test -v -race ./...
	$(MAKE) clean
	@echo "e2e tests finished"

run-tests-in-docker: 
	${CONTAINER_RUNTIME} run --rm --network=host tests:latest

lint:
	cd backend && golangci-lint run

lint-frontend:
	cd frontend && npm run lint

tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $$(go env GOPATH)/bin v2.4.0
