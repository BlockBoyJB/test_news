run:
	go run cmd/app/main.go

compose-up:
	docker-compose up -d --build

compose-down:
	docker-compose down

mocks:
	mockgen -source=internal/repo/categories.go -destination=internal/mocks/repomocks/categories.go -package=repomocks
	mockgen -source=internal/repo/news.go -destination=internal/mocks/repomocks/news.go -package=repomocks
	mockgen -source=internal/repo/txmanager/tx.go -destination=internal/mocks/txmocks/tx.go -package=txmocks
	mockgen -source=internal/service/service.go -destination=internal/mocks/servicemocks/service.go -package=servicemocks

pg-tests:
	docker run --name postgres --rm -d \
		-p 6000:6000 \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=1234567890 \
		-e POSTGRES_DB=postgres postgres:15 -p 6000

init-test-containers: pg-tests

stop-test-containers:
	docker stop postgres

init-tests:
	go test -v ./...

tests: init-test-containers init-tests stop-test-containers
