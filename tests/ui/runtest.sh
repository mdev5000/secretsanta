testName=$1
npx playwright test --trace on --browser firefox --retries 0 -g "${testName}"