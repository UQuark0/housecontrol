#!/bin/bash

git reset --hard
git pull
go build
sudo systemctl restart housecontrold