systemctl stop linux-monitor-go-service
rm -rf /root/linux-monitor-go-service/LinuxMonitorGo_linux_amd64
cp LinuxMonitorGo_linux_amd64 /root/linux-monitor-go-service/LinuxMonitorGo_linux_amd64
chmod 776 /root/linux-monitor-go-service/LinuxMonitorGo_linux_amd64
systemctl start linux-monitor-go-service