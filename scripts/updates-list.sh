# quick exit on error
set -o errexit
set -o nounset

echo "Updates for backend:"
cd backend
# List direct dependencies that need updates.
go list -u -m -f '{{if and (not .Indirect) .Update}}{{.}}{{end}}' all
cd ../
echo ""
echo ""
echo ""

echo "Updates for frontend:"
cd frontend
npx npm-check-updates
cd ../
echo ""
echo ""
echo ""

echo "Updates for UI tests"
cd tests/ui
npx npm-check-updates
cd ../../
echo ""
echo ""
echo ""

echo "Updates for main project"
npx npm-check-updates