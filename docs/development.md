# Development

## Development Environment Setup

### 1. Prerequisites

Requires:
- **Go 1.19**
- **NPM v18**
- Docker
- Docker compose
- [protoc](https://grpc.io/docs/protoc-installation/)

#### 2. Install deps

```bash
make deps.install
```

## Running development environment

Process will be improved soon.

```bash
docker-compose up -d

# in one tab
cd backend
ENV=development go run ./cmd/secretsanta/main.go

# in another tab
cd frontend
npm run dev
```