# quick exit on error
set -o errexit
set -o nounset

echo "Updating backend"
cd backend
go get -t -u ./...
go mod tidy
cd ../

echo "Updating frontend"
cd frontend
npx npm-check-updates -u
npm install
cd ../

echo "Updating UI tests"
cd tests/ui
npx npm-check-updates -u
npm install
cd ../../

echo "Updating main project"
npx npm-check-updates -u
npm install