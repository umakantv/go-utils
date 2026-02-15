module user-service

go 1.19

replace github.com/umakantv/go-utils => ../..

require (
	github.com/gorilla/mux v1.8.0
	github.com/mattn/go-sqlite3 v1.14.16
	github.com/umakantv/go-utils v0.0.0-00010101000000-000000000000
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/jmoiron/sqlx v1.3.5 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
)
