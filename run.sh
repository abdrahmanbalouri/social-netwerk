#!/bin/bash

cd ./backend
go run main.go &

# run frontend from project root (package.json lives in frontend/my-app)
cd ../frontend/my-app
npm install
# ensure npm uses the same node binary (helps with snap-installed node mismatch)
 npm run dev --scripts-prepend-node-path=






# COMMAND CREATE NEW MIGRATION IS ::
# sql-migrate new create_user_table

# BUT BEFORE INSTALL THIS :
# go install github.com/rubenv/sql-migrate/sql-migrate@latest