package common

import (
	"errors"
	"flag"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/speshiy/LinuxMonitorGo/settings"
	"golang.org/x/crypto/bcrypt"
	validator "gopkg.in/go-playground/validator.v9"
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
	dbhdumpFlag := flag.String("dbhdump", "127.0.0.1", "DB host dump")
	portFlag := flag.String("port", "9090", "Port")
	portServiceFlag := flag.String("port_service", "9091", "PortService")
	dbrpFlag := flag.String("dbrp", "1", "DBRP")
	dbrtupFlag := flag.String("dbrtup", "1", "DBRP")
	resourcePathFlag := flag.String("resources_path", "resources", "Resources path")
	backupIntervalHoursFlag := flag.String("backup_interval", "48", "BackupIntervalHours")
	//this command MUST be
	flag.Parse()

	// Получаем конфиг сервера
	readServerConfig()

	settings.IsRelease = *releaseFlag

	settings.DBHostDump = *dbhdumpFlag
	settings.Port = *portFlag
	settings.PortService = *portServiceFlag
	settings.DBRP = *dbrpFlag
	settings.DBRTUP = *dbrtupFlag
	if len(settings.DBRP) == 0 {
		settings.DBRP = "<GENERATE>"
	}

	// *nix /var/www/tuvis/resources/
	settings.ResourcesPath = *resourcePathFlag
	settings.BackupIntervalHours = *backupIntervalHoursFlag

	log.Println("Port = ", settings.Port)
	log.Println("PortService = ", settings.PortService)
	log.Println("DB host dump = ", settings.DBHostDump)

	PrintBinPath()

	initValidator()
}

func readServerConfig() error {
	type Config struct {
		Server struct {
			Name string `yaml:"name" envconfig:"Name"`
			IP   string `yaml:"ip" envconfig:"IP"`
		} `yaml:"server"`
	}

	f, err := os.Open("/root/tuvis_config.yml")
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

//GetNewUUID returns new UUID
func GetNewUUID() (string, error) {
	uid := uuid.NewV4()

	return uid.String(), nil
}

func random(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

//GetNewRandomValue return new random value from start to end
func GetNewRandomValue(start int, finish int) string {
	return strconv.FormatUint(uint64(random(start, finish)), 10)
}

//GetNewClientPassword return new pasword
func GetNewClientPassword() string {
	return strconv.FormatUint(uint64(random(10000, 99999)), 10)
}

func rightPad(s string, padStr string, pLen int) string {
	return s + strings.Repeat(padStr, pLen)
}

//GetCardBarcodeFormatted return barcode for card in format 0000-0000-000000
func GetCardBarcodeFormatted(barcode string) string {
	splittedBarcode := strings.Split(barcode, "-")
	barcode = rightPad(splittedBarcode[0], "0", 4-len(splittedBarcode[0]))
	splittedBarcode[0] = rightPad(splittedBarcode[0], "0", 4-len(splittedBarcode[0]))
	splittedBarcode[1] = rightPad(splittedBarcode[1], "0", 4-len(splittedBarcode[1]))
	return strings.Join(splittedBarcode, "-")
}

/******************************************/
/*PASSWORD*/
/******************************************/

var fishKey = "8043uri2-dlldmg;032i5;gs34"

//HashAndSaltPassword is hashin incoming password
func HashAndSaltPassword(pwd []byte) (string, error) {

	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

//ComparePasswords compare passwords
func ComparePasswords(hashedPwd string, plainPwd []byte) (bool, error) {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		if err != bcrypt.ErrMismatchedHashAndPassword {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

/******************************************/
/*VALIDATOR*/
/******************************************/

//Validate validator
var Validate *validator.Validate

//InitValidator create new global validator
func initValidator() {
	Validate = validator.New()
}

//GetValidationNewSimpleError return a simple error from validation error
func GetValidationNewSimpleError(err error) error {
	var errorString string
	for _, err := range err.(validator.ValidationErrors) {
		errorString += "Field " + err.Namespace() + " has validation Error on '" + err.ActualTag() + "' = '" + err.Param() + "'"
	}
	return errors.New(errorString)
}
