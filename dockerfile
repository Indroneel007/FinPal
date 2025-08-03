FROM golang:1.21-alpine

# Install migrate CLI
RUN apk add --no-cache curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz && mv migrate.linux-amd64 /usr/local/bin/migrate

# Set up workdir and copy files
WORKDIR /app
COPY . .

# Build your Go app
RUN go build -o app .

# Run migrations and then start your app
CMD migrate -source file://db/migration -database "postgresql://finpal_postgres_user:wMJTxyATm6dtr2NGq29Vm7Eala082iEZ@dpg-d27efo6uk2gs73e30sh0-a/finpal_postgres" -verbose up && ./app