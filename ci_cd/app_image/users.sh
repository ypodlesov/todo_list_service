#!/bin/bash

sudo useradd todo_list_service
sudo usermod -aG sudo todo_list_service

# mkdir -p /home/todo_list_service
# useradd -U -d /home/todo_list_service -s /bin/bash todo_list_service

# mkdir -p /etc/todo_list_service/certs
# chown todo_list_service /etc/todo_list_service/certs
# chmod u+rw /etc/todo_list_service/certs