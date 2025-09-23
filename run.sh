#!/bin/bash

cd ./backend
go run main.go &
# cd ../frontend/my-app
# rm -rf node_modules package-lock.json
# npm install

# run frontend from project root (package.json lives in frontend/my-app)
cd ../frontend/my-app
# ensure npm uses the same node binary (helps with snap-installed node mismatch)
npm run dev --scripts-prepend-node-path=true





# COMMAND CREATE NEW MIGRATION IS :::
# sql-migrate new create_user_table

# BUT BEFORE INSTALL THIS :
# go install github.com/rubenv/sql-migrate/sql-migrate@latest