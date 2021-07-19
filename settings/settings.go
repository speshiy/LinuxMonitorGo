package settings

//Server имя сервера
var Server string

//ServerIP url сервера
var ServerIP string

//IsRelease release flag
var IsRelease = false

//DBHostDump сервер, где лежит список баз данных для дампирования
var DBHostDump = "127.0.0.1"

//DBHostMonitoring сервер, куда отправляются данные по мониторингу
var DBHostMonitoring = "172.104.229.191"

//Port application
var Port string

//PortService application
var PortService string

//DatabaseBackupPassword пользователь БД от бэкапов
var DatabaseBackupUser string

//DatabaseBackupPassword пароль пользователя БД от бэкапов
var DatabaseBackupPassword string

//DBRTUP application
var DBRTUP string

//ResourcesPath path of pictures
var ResourcesPath string

//BackupIntervalHours how often to start backup
var BackupIntervalHours string

//IsBackupServer flag
var IsBackupServer = true

//IsBackupInProgress flag
var IsBackupInProgress bool

//BackupingDBName flag
var BackupingDBName string

//Servers мап серверов глобальных IP и локальных IP, смысл в том что некоторые дата центры на VPS позволяют
//задать локальные IP и взаимодействуют они друг с другом по локальным IP. Ну есессно снимать дампы с
//по локалке быстрее чем через глобальные IP (интернет)
var Servers = map[string]string{
	// Localhost
	"127.0.0.1": "127.0.0.1",
}

// DetermineHost определяем хоста в зависимости от тип выполнения программы (prod, dev). Если это prod режим, то вероятнее всего
// мы на машине в локальной сети VPS и поэтому НЕ меняем IP на глобальный
func DetermineHost(host string) string {
	// Если это не prod сервер, то меняем его IP на глобальный
	if !IsRelease {
		host = Servers[host]
	}

	return host
}
