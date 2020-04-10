help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

init:
	@docker-compose down --remove-orphans
	@docker-compose up -d --build --remove-orphans

start:
	@docker-compose up -d --build

stop:
	@docker-compose stop

restart:
	@docker-compose restart

test: ## run tests
	@docker-compose -f docker-compose.test.yml down --remove-orphans
	@docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --remove-orphans
