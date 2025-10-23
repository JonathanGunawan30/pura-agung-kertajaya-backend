# Golang Clean Architecture - Pura Agung Kertajaya

## Description

This repository contains the backend service for the Pura Agung Kertajaya website. It is built using a Clean Architecture in Golang to ensure maintainability and scalability.

The service operates in two main capacities:

Public API: Provides all necessary data for the public-facing website.

Private CMS API: A secure, authenticated API for administrators to manage all site content.

## Architecture

![Clean Architecture](architecture.png)

1. External system perform request (HTTP, gRPC, Messaging, etc)
2. The Delivery creates various Model from request data
3. The Delivery calls Use Case, and execute it using Model data
4. The Use Case create Entity data for the business logic
5. The Use Case calls Repository, and execute it using Entity data
6. The Repository use Entity data to perform database operation
7. The Repository perform database operation to the database
8. The Use Case create various Model for Gateway or from Entity data
9. The Use Case calls Gateway, and execute it using Model data
10. The Gateway using Model data to construct request to external system 
11. The Gateway perform request to external system (HTTP, gRPC, Messaging, etc)

## Tech Stack

- Golang : https://github.com/golang/go
- MySQL (Database) : https://github.com/mysql/mysql-server
- Redis (Cache, Sessions, Rate Limiter) : https://redis.io/
- Cloudflare R2 (S3-Compatible Object Storage) : https://developers.cloudflare.com/r2/

## Framework & Library

- GoFiber (HTTP Framework) : https://github.com/gofiber/fiber
- GORM (ORM) : https://github.com/go-gorm/gorm
- Viper (Configuration) : https://github.com/spf13/viper
- GoDotEnv (Environment Loader) : https://github.com/joho/godotenv
- Golang Migrate (Database Migration) : https://github.com/golang-migrate/migrate
- Go Playground Validator (Validation) : https://github.com/go-playground/validator
- Go-Redis (Redis Client) : https://github.com/redis/go-redis
- AWS SDK for Go (Cloudflare R2 Client) : https://github.com/aws/aws-sdk-go

## Configuration

All configuration is in `config.json` file.

## API Spec

All API Spec is in `api` folder.

## Database Migration

All database migration is in `db/migrations` folder.

### Create Migration

```shell
migrate create -ext sql -dir db/migrations create_table_xxx
```

### Run Migration

```shell
migrate -database "mysql://root:@tcp(localhost:3306)/pura_agung_kertajaya?charset=utf8mb4&parseTime=True&loc=Local" -path db/migrations up
```

## Run Application

### Run unit test

```bash
go test -v ./test/
```

### Run web server

```bash
go run cmd/web/main.go
```