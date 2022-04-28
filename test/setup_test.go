package test

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"timeline/config"
	"timeline/models"
	"timeline/server"

	"github.com/jinzhu/gorm"
)

var (
	DB     *gorm.DB
	err    error
	Server *httptest.Server
)

type Record struct {
	Value string
}

func setup() {
	c := config.MustGetAppConfig()

	var err error
	DB, err = gorm.Open("mysql", c.DB)
	if err != nil {
		panic(err)
	}

	DB.DropTableIfExists(&models.User{})
	DB.AutoMigrate(&models.User{})

	createData()

	router := server.New(DB)
	Server = httptest.NewServer(router)
}

func createData() {
	RunSQL(DB, `
		INSERT INTO users (id, uid, email) VALUES (1, 369, '1048196021@qq.com');
	`)
}

func RunSQL(db *gorm.DB, sqls string) {
	for _, sql := range strings.Split(sqls, "\n") {
		if strings.TrimSpace(sql) != "" {
			db.Exec(strings.TrimSpace(sql))
		}
	}
}

func GetRecords(db *gorm.DB, tableName string, columns string, extra ...string) string {
	var (
		extraSQL = ""
		results  = []string{}
		records  = []Record{}
	)
	if len(extra) > 0 {
		extraSQL = extra[0]
	}
	db.Raw(fmt.Sprintf(`SELECT CONCAT_WS(',', %v) AS value FROM %v %v`, columns, tableName, extraSQL)).Scan(&records)
	for _, record := range records {
		results = append(results, record.Value)
	}
	return strings.Join(results, "; ")
}
