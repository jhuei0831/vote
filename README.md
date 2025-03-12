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