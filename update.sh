#Останавливаем службу
systemctl stop linux-monitor-go-service
#Удаляем бинарник
rm -rf /root/linux-monitor-go-service/LinuxMonitorGo_linux_amd64
cp LinuxMonitorGo_linux_amd64 /root/linux-monitor-go-service/LinuxMonitorGo_linux_amd64
#Стартуем службу
systemctl start linux-monitor-go-service