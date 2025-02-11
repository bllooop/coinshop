build:
		docker-compose up --build
migrate:	
		goose -dir ./backend/migrations postgres "postgres://postgres:54321@localhost:5436/postgres?sslmode=disable" up
migrate-down:	
		goose -dir ./backend/migrations postgres "postgres://postgres:54321@localhost:5436/postgres?sslmode=disable" down