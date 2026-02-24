# KMH Web Profile Redesigned!

This is the backend API for the KMH Web Profile.

- Repository-Service Pattern
- Go + Fiber
- PostgreSQL
- JWT Auth
- RBAC
- sqlc + dbmate
- Scalar docs

# Development Inquiry
### First of all pls make sure you have installed _Go 1.25+_. 
After that, run this in your terminal (this cmd is optional as they may have been generated and migrated by the project coordinator):
```go
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

and this (for Windows):
```bash
scoop install dbmate
```
-- if you're not on Windows, consider checking this repo for correct installation:
https://github.com/amacneil/dbmate?tab=readme-ov-file#installation

After that, you may follow the general steps:
1. Clone the repo

```bash
git clone https://github.com/KMHTelU/KMH-WebProfile-API.git
```
2. Go to the directory/folder

```bash
cd KMH-WebProfile-API-main
```
3. Change the .env.example into .env and fill the values (ask the project coordinator)

```bash
cp .env.example .env
```
4. Run go mod tidy

```go
go mod tidy
```
5. Install air-verse (if you already installed this, then skip)

```go
go install github.com/air-verse/air@latest
```
6. Run air

```bash
air
```

### Swagger Docs
please refer to this: https://swagger.io/docs/specification/v3_0/basic-structure/