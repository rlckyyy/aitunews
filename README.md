# Link to GitHub repository: https://github.com/rlckyyy/aitunews
# AITU News Application

## 1. Clone application

```bash(right now clonning is not possible because of repo is private)
git clone https://github.com/rlckyyy/snippets.git
```

## 2. Get the necessary drivers
```bash
go get github.com/lib/pq
```
## 3. Up the Docker Container via Terminal in root directory
```bash
docker compose up -d
```

## 4. Connect to DB, create table and index 
```
CREATE TABLE news
(
    id       SERIAL PRIMARY KEY,
    title    VARCHAR(100) NOT NULL,
    content  TEXT         NOT NULL,
    created  TIMESTAMP    NOT NULL,
    category VARCHAR(15)
);

CREATE INDEX idx_news_created ON news (created);
```
## 5. Start the application via
```bash
go run ./cmd/web
```

