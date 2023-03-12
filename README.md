 # Secret Santa
 
Work in progress, please check back later.

## Environment Setup

See instructions [here](./docs/dev-setup.md) for setting up your environment.

## Build the project

```bash
make build
# or to fully rebuild
make build.refresh


# run the binary
export MONGO_URI="mongo://user:password@uri-to-mongo"
./backend/_build/secretsanta
```

## Development

See [Development](./docs/development.md) docs for details on setup, running, testing, etc.