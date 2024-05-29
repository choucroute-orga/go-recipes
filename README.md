# Recipes

This MS is a simple recipes management system. It is a REST API that allows you to manage the recipes.

### Start the server

```bash
cp .env.example .env
export $(cat .env | xargs)
docker-compose up
go run main.go
```
