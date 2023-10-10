APP = mccs
# Get the latest Git tag, or if not available, get the latest commit hash.
GIT_TAG = $(shell if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ]; then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
# Get the current date and time in UTC in ISO 8601 format.
BUILD_DATE = $(shell TZ=UTC date +%FT%T%z)
# Get the hash of the latest commit.
GIT_COMMIT = $(shell git log --pretty=format:'%H' -n 1)
# Check the status of the Git tree, determining whether it's clean or dirty.
GIT_TREE_STATUS = $(shell if git status | grep -q 'clean'; then echo clean; else echo dirty; fi)

# Production target for starting the production server.
production:
	@echo "============= Starting production server ============="
	GIT_TAG=${GIT_TAG} BUILD_DATE=${BUILD_DATE} GIT_COMMIT=${GIT_COMMIT} GIT_TREE_STATUS=${GIT_TREE_STATUS} \
	docker-compose -f docker-compose.production.yml up --build

# Clean target for removing the application.
clean:
	@echo "============= Removing app ============="
	rm -f ${APP}

# Run target for starting the server using the development Docker Compose configuration.
run:
	@echo "============= Starting server ============="
	docker-compose -f docker-compose.dev.yml up --build

# Test target for running unit tests on the application.
test:
	@echo "============= Running tests ============="
	go test ./...

# Seed target for generating seed data.
seed:
	@echo "============= Generating seed data ============="
	go run cmd/seed/main.go -config="seed"

# es-restore target for restoring Elasticsearch data.
es-restore:
	@echo "============= Restoring Elasticsearch data ============="
	go run cmd/es-restore/main.go -config="seed"

# pg-setup target for setting up PostgreSQL accounts.
pg-setup:
	@echo "============= Setting up PostgreSQL accounts ============="
	go run cmd/pg-setup/main.go -config="seed"
