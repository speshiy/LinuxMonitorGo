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
	dbmhFlag := flag.String("dbmhFlag", "127.0.0.1", "DBHostMonitoring")
	dbmuFlag := flag.String("dbmuFlag", "root", "DatabaseMonitorUser")
	dbmpFlag := flag.String("dbmpFlag", "1", "DatabaseMonitorPassword")
	//this command MUST be
	flag.Parse()

	// Получаем конфиг сервера
	readServerConfig()

	settings.IsRelease = *releaseFlag

	settings.Port = *portFlag
	settings.PortService = *portServiceFlag
	settings.DBHostMonitoring = *dbmhFlag
	settings.DatabaseMonitorUser = *dbmuFlag
	settings.DatabaseMonitorPassword = *dbmpFlag

	log.Println("Port = ", settings.Port)
	log.Println("PortService = ", settings.PortService)
	log.Println("DB host dump = ", settings.DBHostMonitoring)

	PrintBinPath()
}

func readServerConfig() error {
	type Config struct {
		Server struct {
			Name string `yaml:"name" envconfig:"Name"`
			IP   string `yaml:"ip" envconfig:"IP"`
		} `yaml:"server"`
	}

	f, err := os.Open("/root/LinuxMonitorGo_config.yml")
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

	return nil
}
