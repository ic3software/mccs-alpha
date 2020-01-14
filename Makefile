APP=mccs

GIT_TAG = $(shell if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
BUILD_DATE = $(shell TZ=UTC date +%FT%T%z)
GIT_COMMIT = $(shell git log --pretty=format:'%H' -n 1)
GIT_TREE_STATUS = $(shell if git status|grep -q 'clean';then echo clean; else echo dirty; fi)

production:
	@echo "=============starting production server============="
	GIT_TAG=${GIT_TAG} BUILD_DATE=${BUILD_DATE} GIT_COMMIT=${GIT_COMMIT} GIT_TREE_STATUS=${GIT_TREE_STATUS} \
	docker-compose -f docker-compose.production.yml up --build

clean:
	@echo "=============removing app============="
	rm -f ${APP}

run:
	@echo "=============starting server============="
	docker-compose -f docker-compose.dev.yml up --build

test:
	@echo "=============running test============="
	go test ./...

seed:
	@echo "=============generating seed data============="
	go run cmd/seed/main.go -config="seed"

es-restore:
	@echo "=============restoring es data============="
	go run cmd/es-restore/main.go -config="seed"

pg-setup:
	@echo "=============setup pg accounts============="
	go run cmd/pg-setup/main.go -config="seed"
