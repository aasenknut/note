redis:
	docker run --name my-redis -p 6379:6379 -d redis

test-db:
	docker run --name redis-test-ctxptrl -p 6370:6379 -d redis

integration-test:
	go test ./client/redis -db-addr="localhost:6543"
