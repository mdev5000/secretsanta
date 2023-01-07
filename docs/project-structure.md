
# Project Structure

`backend/` contains backend server code, written in Go.


`frontend/` contains the SvelteKit SPA frontend written in Typescript. The frontend is compiled into a set of SPA files
which are then embedded in the backend Go binary.


`schemas/` contains [json type definitions](https://jsontypedef.com/) for generating Go and Typescript request
datastructures and can possibly be used to validate requests on the frontend in the future.


`scripts/` contains misc development scripts for building the project, etc.


`tests/` contains the UI testing code for acceptance level testing.