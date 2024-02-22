
secrets_dir := secrets/
migrations_dir := migrations/
db_pass := $(file < $(secrets_dir)/db_password.txt)
db_user := $(file < $(secrets_dir)/db_user.txt)
db_uri := $(db_user):$(db_pass)@localhost:5432/marketmoogle?sslmode=disable

build:
	go build github.com/level-5-pidgey/MarketMoogleApi

create_migration: 
	migrate -source file://$(migrations_dir) -database postgres://$(db_uri) create -ext .sql -dir $(migrations_dir) -seq -digits 3 mm_migration

migrate_up:
	migrate -source file://$(migrations_dir) -database postgres://$(db_uri) up

migrate_down:
	migrate -source file://$(migrations_dir) -database postgres://$(db_uri) down

run:
	air -- -setup