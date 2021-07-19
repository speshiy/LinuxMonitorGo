package cpurge

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/speshiy/LinuxMonitorGo/_main/controllers/cbackup"
	"github.com/speshiy/LinuxMonitorGo/_main/models/muser"
	"github.com/speshiy/LinuxMonitorGo/database"
	"github.com/speshiy/LinuxMonitorGo/settings"
)

// Очищает неискользуемые базы данных (больше года)
func PurgeUnusedDatabases(c *gin.Context) {
	var err error
	var DBMain database.DBUserConnect
	var DBClient database.DBUserConnect
	var userConnections []muser.UserConnectionData

	err = DBMain.Open(settings.DBHostDump, "main", "rtu", settings.DBRTUP, "UTC")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}
	defer DBMain.Close()

	err = DBClient.Open(settings.DBHostDump, "client", "rtu", settings.DBRTUP, "UTC")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}
	defer DBClient.Close()

	log.Println("Получение коннектов ко всем аккаунтам...")

	err = muser.GetConnectionsDataByDBName(c, "", &userConnections, DBMain.Connect)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	// Бежим по коннектам и проверяем активность аккаунта
	// 1) Транзакции
	// 2) Оплаты
	// 3) Партнёр или нет
	// 4) Мобилка
	for _, uc := range userConnections {
		log.Println(uc.DBName, "начат процесс проверки на удаление...")

		if uc.DBName == "db_pssolo_bk" {
			continue
		}

		func() {
			var DBUser database.DBUserConnect

			err = DBUser.Open(settings.DetermineHost(uc.DBHost), uc.DBName, "rtu", settings.DBRTUP, "UTC")
			if err != nil {
				return
			}
			defer DBUser.Close()

			var transactionCreatedAt time.Time
			var transactionsQuantity int

			sql := `SELECT
						IFNULL(MAX(ct.created_at), DATE("2019-03-01 00:00:01")) AS created_at,
						COUNT(*) AS transactionsQuantity
					FROM
						d_cards_transactions AS ct`

			row := DBUser.Connect.Raw(sql).Select("transactionCreatedAt, transactionsQuantity").Row()
			err = row.Scan(&transactionCreatedAt, &transactionsQuantity)
			if err != nil {
				if strings.Contains(err.Error(), "no rows in result set") {
					transactionCreatedAt = time.Date(2019, 03, 01, 00, 00, 01, 00, time.UTC)
				} else {
					c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
					return
				}
			}

			daysTransaction := time.Since(transactionCreatedAt).Hours() / 24
			daysCreated := time.Since(uc.CreatedAt).Hours() / 24

			// Если аккаунт создан более года наза и транзакций не было больше 180 дней и их в аккаунте меньше 500
			// и это не партнёр, то удаляем базу
			if daysCreated >= 240 && daysTransaction >= 240 && transactionsQuantity <= 500 {
				log.Println(uc.DBName, "удаление инициализировано...")
				err = DropDatabase(c, DBMain.Connect, DBClient.Connect, DBUser.Connect, uc)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
					return
				}
			}
		}()
	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "message": "Неиспользуемые аккаунты были удалены успешно"})
}

// DropDatabase удаляем базу
func DropDatabase(c *gin.Context, DBMain *gorm.DB, DBClient *gorm.DB, DBUser *gorm.DB, userConnection muser.UserConnectionData) error {
	var userPrimeID uint
	var isPartner uint

	sql := `SELECT
				u.id as userPrimeID,
				IFNULL(u.is_partner, 0) as isPartner
			FROM
				s_users AS u
			WHERE
				u.connection_data_id = ? AND
				u.is_prime = 1`

	row := DBMain.Raw(sql, userConnection.ID).Select("userPrimeID, isPartner").Row()
	err := row.Scan(&userPrimeID, &isPartner)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil
		} else {
			return err
		}
	}

	if isPartner == 1 {
		return nil
	}

	if userPrimeID <= 0 {
		return nil
	}

	log.Println(userConnection.DBName, "бэкап начат...")

	// Бэкапируем базу
	err = cbackup.DoBackup(c, userConnection.DBName, DBMain, nil, false, true)
	if err != nil {
		return err
	}

	log.Println(userConnection.DBName, "очиcтка таблиц в DBClient...")

	// Удаляем все связи с компанией в DBClients
	r := DBClient.Exec("DELETE FROM s_clients_users_link WHERE user_prime_id = ?", userPrimeID)
	if r.Error != nil {
		return r.Error
	}

	// Удаляем все компании в DBClients
	r = DBClient.Exec("DELETE FROM s_company WHERE user_prime_id = ?", userPrimeID)
	if r.Error != nil {
		return r.Error
	}

	// Удаляем все промокода в DBClients
	r = DBClient.Exec("DELETE FROM s_promocodes WHERE user_prime_id = ?", userPrimeID)
	if r.Error != nil {
		return r.Error
	}

	// Удаляем все столы в DBClients
	r = DBClient.Exec("DELETE FROM s_rooms_tables WHERE user_prime_id = ?", userPrimeID)
	if r.Error != nil {
		return r.Error
	}

	log.Println(userConnection.DBName, "очиcтка таблиц в DBMain...")

	// Удаляем все url's в DBMain
	r = DBMain.Exec("DELETE FROM sys_website_urls WHERE user_id = ?", userPrimeID)
	if r.Error != nil {
		return r.Error
	}

	// Удаляем все компании в DBMain
	r = DBMain.Exec("DELETE FROM s_users_companies WHERE user_id = ?", userPrimeID)
	if r.Error != nil {
		return r.Error
	}

	// Удаляем все логины в DBMain
	r = DBMain.Exec("DELETE FROM s_users WHERE connection_data_id = ?", userConnection.ID)
	if r.Error != nil {
		return r.Error
	}

	// Удаляем базу пользователя на старом хосте
	err = database.DropDatabaseOnHost(userConnection.DBHost, userConnection.DBName, userConnection.DBUser)
	if err != nil {
		return err
	}

	// Удаляем коннекты
	r = DBMain.Exec("UPDATE s_users_connection_data_t SET deleted_at = NOW() WHERE id = ?", userConnection.ID)
	if r.Error != nil {
		return r.Error
	}

	return nil
}
