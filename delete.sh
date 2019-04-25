#!/bin/bash

docker rm demomicroservice_auth_1 -f
docker rm demomicroservice_employee_1 -f
docker rm demomicroservice_role_1 -f
docker rm demomicroservice_schedule_1 -f
docker rm demomicroservice_mongo_seed_1 -f
docker image rm v/auth
docker image rm v/employee
docker image rm v/role
docker image rm v/schedule
docker image rm demomicroservice_mongo_seed