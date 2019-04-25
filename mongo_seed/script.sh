#!/bin/bash
mongoimport -h mongo:27017 -d userdb -c user user.json
mongoimport -h mongo:27017 -d employeedb -c employee employee.json
mongoimport -h mongo:27017 -d roledb -c role role.json
mongoimport -h mongo:27017 -d scheduledb -c schedule schedule.json