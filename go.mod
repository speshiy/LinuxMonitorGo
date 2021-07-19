module github.com/speshiy/LinuxMonitorGo

go 1.15

require (
	github.com/dchest/captcha v0.0.0-20200903113550-03f5f0333e1f // indirect
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/gin-contrib/cors v0.0.0-20170318125340-cf4846e6a636
	github.com/gin-gonic/gin v1.5.0
	github.com/go-playground/validator v9.31.0+incompatible // indirect
	github.com/go-resty/resty v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/snappy v0.0.3 // indirect
	github.com/guillermo/go.procmeminfo v0.0.0-20131127224636-be4355a9fb0e
	github.com/jinzhu/gorm v1.9.16
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/robfig/cron v1.2.0
	github.com/satori/go.uuid v1.2.0
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	github.com/xlab/closer v0.0.0-20190328110542-03326addb7c2 // indirect
	golang.org/x/crypto v0.0.0-20191205180655-e7c4368fe9dd
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5 // indirect
	gopkg.in/go-playground/validator.v9 v9.29.1
	gopkg.in/resty.v1 v1.12.0 // indirect
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/go-resty/resty => gopkg.in/resty.v1 v1.11.0
