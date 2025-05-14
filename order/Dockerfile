FROM golang:1.23

WORKDIR /app

# Enable CGO for SQLite
ENV CGO_ENABLED=1

# Install SQLite and migrate
RUN apt-get update && apt-get install -y \
    sqlite3 \
    libsqlite3-dev \
    && go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest


COPY . .

RUN go build -o main ./cmd/main.go

EXPOSE 50053

RUN touch /app/storage/database.db && \
    migrate -source file://migrations -database sqlite3:///app/storage/app.database up

CMD ["./main"]