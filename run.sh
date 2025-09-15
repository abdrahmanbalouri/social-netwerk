#!/bin/bash
#rm -rf node_modules package-lock.json
#npm install

cd ./backend
go run main.go &

cd ../frontend/my-app/app
npm run dev
