# Development

## Development Environment Setup

### 1. Prerequisites

Requires:

- **Go 1.19**
- **NPM v18**
- Docker
- Docker compose
- [protoc](https://grpc.io/docs/protoc-installation/)


For updating dependencies ([npm-check-updates](https://github.com/raineorshine/npm-check-updates)):
```bash
npm install -g npm-check-updates
```

### 2. Install deps

```bash
make deps.install
```

---

## Running development environment

```bash
make dev.docker.up
make dev.run

```

### View and edit mongo data

Go to http://localhost:8081

---

## Testing

### Run all tests

```bash
make test.all
```

### UI tests

```bash
cd tests/ui

# run the ui tests
# by default this assumes the dev environment is running
npm run test

# you can also test again a running instance
BASE_URL="http://localhost:3000" npm run test

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

# run all tests and check for race conditions.
make test.race

# don't run db tests
NODB=1 make test
```

### Frontend

```bash
cd frontend

# Run all tests
make test

# Run test watcher
npm run test
```

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


## Updates

```bash
# List changes
make updates.list

# Install updates
make updates.install
```

You'll also likely need to update playwright version in [run-test-ui.sh](scripts/run-test-ui.sh).