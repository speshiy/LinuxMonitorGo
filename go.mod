module github.com/speshiy/LinuxMonitorGo

go 1.15

require (
	github.com/gin-contrib/cors v0.0.0-20170318125340-cf4846e6a636
	github.com/gin-gonic/gin v1.5.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/guillermo/go.procmeminfo v0.0.0-20131127224636-be4355a9fb0e
	github.com/jinzhu/gorm v1.9.16
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/robfig/cron v1.2.0
	golang.org/x/sys v0.0.0-20210423082822-04245dca01da // indirect
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/go-resty/resty => gopkg.in/resty.v1 v1.11.0
