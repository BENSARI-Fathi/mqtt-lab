#!/bin/bash

docker run \
--detach \
--name=mysql-image \
--env="MYSQL_ROOT_PASSWORD=my_password" \
--publish 6603:3306 \
mysql
