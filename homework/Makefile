build-all: clear-all
	docker-compose build

clear-all:
	docker-compose down --remove-orphans --rmi all --volumes
	docker system prune --force --all --volumes
	docker network prune --force

run-all:
	docker-compose up --force-recreate --build -d cart loms pg-0 pg-1 testdb kafka0 kafka-init-topics notifier-1 notifier-2 notifier-3

run-monitoring:
	docker-compose up --force-recreate --build -d prometheus grafana jaeger kafka-ui

migrate-stage:
	docker-compose exec -T loms /bin/sh -c "go install github.com/pressly/goose/v3/cmd/goose@latest"
	docker-compose exec -T loms /bin/sh -c "goose -dir ./migrations postgres \"host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable\" up"

run-e2e-tests:
	docker-compose exec -T cart /bin/sh -c "cd e2e/app && go test"

e2e-tests: clear-all run-all migrate-stage run-e2e-tests clear-all

lint-cart:
	cd cart && make lint

lint-loms:
	cd loms && make lint

lint-notifier:
	cd notifier && make lint

lint: lint-cart lint-loms lint-notifier

test-cart:
	cd cart && make test

test-loms:
	cd loms && make test

unit-tests: test-cart test-loms

migrate-test:
	docker-compose exec -T loms /bin/sh -c "go install github.com/pressly/goose/v3/cmd/goose@latest"
	docker-compose exec -T loms /bin/sh -c "goose -dir ./migrations postgres \"host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable\" up"

loms-integration-tests:
	docker-compose exec -T loms /bin/sh -c "cd tests/cmd && go test"

integration-tests: clear-all run-all migrate-test loms-integration-tests clear-all

take-pprof-cart:
	go tool pprof http://localhost:8082/debug/pprof/profile\?seconds\=10

take-pprof-loms:
	go tool pprof http://localhost:8081/debug/pprof/profile\?seconds\=10

take-pprof-profiles: take-pprof-cart take-pprof-loms
