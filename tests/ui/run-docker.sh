#!/bin/bash
cd /tests
npm install
REPORT_PATH=playwright-report-ci BASE_URL="http://secretsanta:3000" npm run test