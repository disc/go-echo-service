if [ $1 = 0 ]; then
    /usr/bin/systemctl stop echo-service@{8080..8081}
    /usr/bin/systemctl disable echo-service@{8080..8081}
    /usr/bin/systemctl daemon-reload
fi