swagger:
	swag init --dir ./cmd,./api/http,./internal --output ./docs

run:
	~/go/bin/air

tidy:
	go mod tidy

migrate-up:
	migrate -path migrations -database "postgres://postgres:yourpassword@localhost:5432/edu?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:yourpassword@localhost:5432/edu?sslmode=disable" down

migrate-create:
	migrate create -ext sql -dir migrations $(name)