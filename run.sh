#!/bin/bash

cd ./backend
go run main.go &

cd ../frontend/my-app/app
npm run dev
