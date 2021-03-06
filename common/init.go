package common

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/speshiy/LinuxMonitorGo/settings"
	"gopkg.in/yaml.v2"
)

//PrintBinPath print bin path
func PrintBinPath() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	log.Println("Bin = ", exPath)
}

//InitGlobalVars init global variables
func InitGlobalVars() {
	releaseFlag := flag.Bool("release", false, "release")
	portFlag := flag.String("port", "9090", "Port")
	portServiceFlag := flag.String("port_service", "9091", "PortService")
	dbmhFlag := flag.String("dbmhFlag", "127.0.0.1", "DatabaseMonitorHost")
	dbmnFlag := flag.String("dbmnFlag", "monitoring", "DatabaseMonitorName")
	dbmuFlag := flag.String("dbmuFlag", "root", "DatabaseMonitorUser")
	dbmpFlag := flag.String("dbmpFlag", "1", "DatabaseMonitorPassword")

	flag.Parse()

	// Get server config
	readServerConfig()

	settings.IsRelease = *releaseFlag
	settings.Port = *portFlag
	settings.PortService = *portServiceFlag
	settings.DatabaseMonitorHost = *dbmhFlag
	settings.DatabaseMonitorName = *dbmnFlag
	settings.DatabaseMonitorUser = *dbmuFlag
	settings.DatabaseMonitorPassword = *dbmpFlag

	log.Println("Port = ", settings.Port)
	log.Println("PortService = ", settings.PortService)
	log.Println("DB host dump = ", settings.DatabaseMonitorHost)

	PrintBinPath()
}

func readServerConfig() error {
	type Config struct {
		Server struct {
			Name             string `yaml:"name" envconfig:"Name"`
			IP               string `yaml:"ip" envconfig:"IP"`
			IsBackupServer   bool   `yaml:"is_backup_server" envconfig:"IsBackupServer"`
			IsDatabaseServer bool   `yaml:"is_database_server" envconfig:"IsDatabaseServer"`
		} `yaml:"server"`
	}

	f, err := os.Open("LinuxMonitorGo_config.yml")
	if err != nil {
		return err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return err
	}

	settings.Server = cfg.Server.Name
	settings.ServerIP = cfg.Server.IP
	settings.IsBackupServer = cfg.Server.IsBackupServer
	settings.IsDatabaseServer = cfg.Server.IsDatabaseServer

	return nil
}
