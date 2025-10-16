## environment setup

### 1. Create .env file
```
cp .env.example .env
```

### 2. Modify .env file

* Make sure your build platform is linux/arm64 or linux/amd64 or else ..

### 3. docker compose
```bash
docker-compose up -d
```

### 4. Check environment setup success
```bash
# Get response {"message":"health check: PORT 9443"}
curl http://localhost:9443/hc
```

## Database setup

### 1. Create database

* This project use [goose](https://github.com/pressly/goose) and [Go migrations](https://github.com/pressly/goose/tree/master/examples/go-migrations) to manage database schema.

```bash
# use the following command to create a new migration
./migrator goose up
```
## 2. Import data

```bash
# use the following command to import data
./migrator goose -action=seed up
```

### 3. If create new migration
```bash
# use the following command to create a new migration
go build -o migrator ./cmd/dbmigrate
```

## Build & Run
```bash
gin -a 3000 -p 9443 run main.go
```
