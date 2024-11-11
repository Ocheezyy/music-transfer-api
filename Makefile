# Makefile
.PHONY: api-run consumer-run api-tidy consumer-tidy api-test consumer-test api-migrate api-build consumer-build tidy-all

API_DIR=api
CONSUMER_DIR=consumer

api-run:
	cd $(API_DIR) && go run main.go

api-migrate:
	cd $(API_DIR) && go run migrate/migrate.go

api-tidy:
	cd $(API_DIR) && go mod tidy

api-test:
	cd $(API_DIR) && go test ./controllers -v

api-build:
	cd $(API_DIR) && go build -o main main.go

consumer-run:
	cd $(CONSUMER_DIR) && go run main.go

consumer-tidy:
	cd $(CONSUMER_DIR) && go mod tidy

consumer-test:
	cd $(CONSUMER_DIR) && go test -v

consumer-build:
	cd $(CONSUMER_DIR) && go build -o main main.go

tidy-all:
	make api-tidy && make consumer-tidy