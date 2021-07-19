package cmonitor

import (
	"log"

	"github.com/go-resty/resty"
	procmeminfo "github.com/guillermo/go.procmeminfo"
	"github.com/speshiy/LinuxMonitorGo/_monitoring/models/mmonitor"
	"github.com/speshiy/LinuxMonitorGo/settings"
)

//PostRAM post info about ram
func PostRAM() {
	if settings.Server == "" {
		return
	}

	var monitor mmonitor.Monitor

	meminfo := &procmeminfo.MemInfo{}
	meminfo.Update()

	monitor.IP = settings.ServerIP
	monitor.Server = settings.Server

	monitor.IsBackupServer = settings.IsBackupServer
	monitor.IsBackupInProgress = settings.IsBackupInProgress
	monitor.BackupingDBName = settings.BackupingDBName

	monitor.RAM = meminfo.Total()
	monitor.RAMUsed = meminfo.Used()

	swapTotal := (*meminfo)["SwapTotal"]
	swapFree := (*meminfo)["SwapFree"]
	monitor.Swap = swapTotal
	monitor.SwapUsed = swapTotal - swapFree

	request := resty.New()
	resp, err := request.R().
		SetHeader("Content-Type", "application/json").
		SetBody(monitor).
		Post("https://admin.tuvis.world/api/app-admin/monitor")

	if err != nil {
		log.Println("PostRAM.request - " + err.Error())
	}

	if !resp.IsSuccess() {
		log.Println("PostRAM - " + resp.String())
	}
}
