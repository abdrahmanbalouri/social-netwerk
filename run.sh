#!/bin/bash
clear
cd ./backend
go run main.go &

cd ../frontend/my-app
npm install
npm run dev 
