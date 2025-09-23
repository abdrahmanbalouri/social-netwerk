#!/bin/bash

cd ./backend
go run main.go &

# run frontend from project root (package.json lives in frontend/my-app)
cd ../frontend/my-app
npm install
# ensure npm uses the same node binary (helps with snap-installed node mismatch)
npm run dev --scripts-prepend-node-path=true
