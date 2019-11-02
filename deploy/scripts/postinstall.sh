# Install
if [ $1 = 1 ]; then
    /usr/bin/systemctl daemon-reload
    /usr/bin/systemctl enable echo-service@{8080..8081}
    /usr/bin/systemctl restart echo-service@{8080..8081}
fi

# Upgrade
if [ $1 = 2 ]; then
    /usr/bin/systemctl daemon-reload
    /usr/bin/systemctl enable echo-service@{8080..8081}
    /usr/bin/systemctl restart echo-service@{8080..8081}
fi