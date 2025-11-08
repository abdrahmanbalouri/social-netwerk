#!/bin/bash
clear
cd ./backend
go run ./cmd &

cd ../frontend
npm install
npm run dev 


# IF YOU GET THIS ERROR : command not found: sql-migrate --- FOLLOW THESE STEPS 
# 1- RUN THIS : go env GOPATH
# 2- TAKE THE RESULT WJM3HA M3A HADO : /bin/sql-migrate 
# 3- WRAHA ZID : new cretae_tableName
# flkhr atkoun l command haka : /home/youruser/go/bin/sql-migrate new create_user_table -- ri l path li radi ttbdl 3la hssab chno l path dyalk