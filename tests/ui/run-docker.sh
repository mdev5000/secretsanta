#!/bin/bash
cd /tests
npm install
REPORT_PATH=playwright-report-ci npm run test