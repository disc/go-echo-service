#!/bin/sh
set -e

if [ "$1" = configure ]; then
    /usr/bin/systemctl daemon-reload
    /usr/bin/systemctl enable echo-servicee@8080
    /usr/bin/systemctl enable echo-servicee@8081
fi