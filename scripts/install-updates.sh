# quick exit on error
set -o errexit
set -o nounset
set -o pipefail

echo "Updating backend"
cd backend
go get -t -u ./...
go mod tidy
cd ../

echo "Updating frontend"
cd frontend
ncu -u
npm install
cd ../

echo "Updating UI tests"
cd tests/ui
ncu -u
npm install
cd ../../

echo "Updating main project"
ncu -u
npm install