#!/bin/bash
export GOOS=linux
export CGO_ENABLED=0

cd AuthService;go get;go build -o auth;echo built `pwd`;cd ..

cd EmployeeService;go get;go build -o employee;echo built `pwd`;cd ..

cd RoleService;go get;go build -o role;echo built `pwd`;cd ..

cd ScheduleService;go get;go build -o schedule;echo built `pwd`;cd ..

export GOOS=darwin