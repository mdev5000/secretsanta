
# Project Structure

`backend/` contains backend server code, written in Go.


`frontend/` contains the SvelteKit SPA frontend written in Typescript. The frontend is compiled into a set of SPA files
which are then embedded in the backend Go binary.


`protos/` protobuf files used for generating json request/response objects.


`scripts/` contains misc development scripts for building the project, etc.


`tests/` contains the UI testing code for acceptance level testing.