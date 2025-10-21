#!/bin/bash

cd ./backend
go run main.go &

# run frontend from project root (package.json lives in frontend/my-app)
cd ../frontend/my-app
npm install
# ensure npm uses the same node binary (helps with snap-installed node mismatch)
npm run dev --scripts-prepend-node-path=true





# COMMAND CREATE NEW MIGRATION IS ::
# sql-migrate new create_user_table

# BUT BEFORE INSTALL THIS :
# go install github.com/rubenv/sql-migrate/sql-migrate@latest

# IF YOU GET THIS ERROR : command not found: sql-migrate --- FOLLOW THESE STEPS 
# 1- RUN THIS : go env GOPATH
# 2- TAKE THE RESULT WJM3HA M3A HADO : /bin/sql-migrate 
# 3- WRAHA ZID : new cretae_tableName
# flkhr atkoun l command haka : /home/youruser/go/bin/sql-migrate new create_user_table -- ri l path li radi ttbdl 3la hssab chno l path dyalk