# LinuxMonitorGo

This service allow to view RAM usage on a server

# Install
You can use already compiled app LinuxMonitorGo_linux_arm64
For install
1) Make directory <strong>/root/linux-monitor-go-service</strong>
2) Copy <strong>LinuxMonitorGo_linux_arm64</strong> and <strong>linux-monitor-go-service.service</strong> to your server in <strong>/root/linux-monitor-go-service</strong>;
3) Create and start the service with following commands:<br>
cp -rf linux-monitor-go-service.service /lib/systemd/system/<br>
systemctl daemon-reload<br>
systemctl enable linux-monitor-go-service<br>
systemctl start linux-monitor-go-service

If you want to change database which collect monitoring data please change flags is <strong>linux-monitor-go-service.service</strong>:<br>
<strong>-dbmhFlag=127.0.0.1</strong> - database host<br>
<strong>-dbmnFlag=monitoring</strong> - database name<br>
<strong>-dbmuFlag=root</strong> - database user<br>
<strong>-dbmpFlag=1</strong> - database password

# P.S
You can reload service by http for this open <strong>/routes/routes.go</strong> and read routes
