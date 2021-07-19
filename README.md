# LinuxMonitorGo

This service allow to view RAM usage on a server

# Install
You can use already compiled app LinuxMonitorGo_linux_arm64
For install
1) Make directory /root/linux-monitor-go-service
2) Copy LinuxMonitorGo_linux_arm64 and linux-monitor-go-service.service to your server in /root/linux-monitor-go-service;
3) Create and start the service with following commands:
cp -rf linux-monitor-go-service.service /lib/systemd/system/
systemctl daemon-reload
systemctl enable linux-monitor-go-service
systemctl start linux-monitor-go-service
