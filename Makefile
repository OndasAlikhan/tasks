include .env

up:
	goose -dir ./migrations postgres "$(DATABASE_URL)" up
down:
	goose -dir ./migrations postgres "$(DATABASE_URL)" down
create:
	goose -dir ./migrations create $(name) sql