#!/bin/sh
set -e
if [ "$1" = remove ]; then
    /usr/bin/systemctl stop echo-service@8080
    /usr/bin/systemctl stop echo-service@8081
    /usr/bin/systemctl disable echo-service@8080
    /usr/bin/systemctl disable echo-service@8081
    /usr/bin/systemctl daemon-reload
fi