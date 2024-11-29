.PHONY: build up down run local docker clean wrk-docker vegeta-docker test-local test-docker

# Local targets
run:
	cd ./ && go run cmd/app/main.go

local: run

# Docker targets
build:
	APP_ENV=docker docker-compose build

up:
	APP_ENV=docker docker-compose up -d

down:
	docker-compose down

clean: down
	docker-compose rm -f


docker: build up

# Goals for integration tests in Docker
test-integration-docker:
	@echo "Running integration tests in Docker"
	docker-compose run app go test ./... -v

# Docker stress-tests with WRK and Vegeta
wrk-docker:
	@echo "Running wrk test on /order endpoint (Docker)"
	wrk -t12 -c400 -d30s http://localhost:8080/order?id=b563feb7b2b84b6test

vegeta-docker:
	@echo "Running vegeta test on /order endpoint (Docker)"
	vegeta attack -duration=30s -rate=10/s -targets=targets.txt > results.bin


test-docker: wrk-docker vegeta-docker
	@echo "Completed Docker tests with wrk and vegeta"
	@echo "Docker vegeta metrics saved to metrics-docker.json and plot to plot-docker.html"
	open plot-docker.html