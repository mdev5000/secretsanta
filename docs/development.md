# Development

## Development Environment Setup

### 1. Prerequisites

Requires:

- **Go 1.19**
- **NPM v18**
- Docker
- Docker compose
- [protoc](https://grpc.io/docs/protoc-installation/)

### 2. Install deps

```bash
make deps.install
```

---

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

---

## Testing

### Acceptance tests

```bash
cd tests/ui

# run the ui tests
npm run test

# view the report
npm run report

# report will in playwright-report folder, if you wish to make it available elsewhere
ls playwright-report
```

### Backend

```bash
cd backend

# run all tests
make test

# don't run db tests
NODB=1 make test
```

### Frontend

Will add later.

---

## Creating new request / response objects

All JSON requests and responses should be created as protobuf messages so that the frontend and backend contracts stay
consistent.

You can create new message by:

1. Adding your new message to the `/protos` folder.

Ex.

```protobuf
syntax = "proto3";

package requests;

option go_package = "requests/gen";

message MyRequest {
  string field1 = 1;
  int64 field2 = 2;
}
```

2. Generate the **Go** and **Typescript** files via:

```bash
make schemas.gen
```