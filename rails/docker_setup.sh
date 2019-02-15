#!/bin/bash

# build docker images
docker-compose build

# create docker-compose compatible db config
cp config/database.yml.docker config/database.yml

# start docker services
docker-compose up -d

# set up dbs
docker-compose exec web rake db:create
docker-compose exec web rake db:migrate

