package test

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"net/url"
	"strings"
	"timeline/app/users"
	"timeline/config"
	"timeline/models"
	"timeline/server"

	httpclient "github.com/ddliu/go-httpclient"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
)

var (
	DB     *gorm.DB
	Server *httptest.Server
)

type Record struct {
	Value string
}

func setup() {
	// appconfig := config.Configuration.APPConfig
	configor.Load(&config.Config, "config.yml")
	DB, err := gorm.Open(config.Config.APPConfig.DB)
	if err != nil {
		panic(err)
	}

	DB.DropTableIfExists(&models.User{})
	DB.AutoMigrate(&models.User{})

	// createData()

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

func Post(url string, values url.Values) string {
	client := httpclient.NewHttpClient()
	client.Headers = make(map[string]string)
	params := make(map[string]string)
	for k, value := range values {
		params[k] = value[0]
	}
	resp, _ := client.Post(Server.URL+url, params)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func httpPost(url string, values url.Values) string {
	req := httptest.NewRequest("POST", Server.URL+url, nil)
	w := httptest.NewRecorder()
	userhandler := users.NewHandler(DB)
	userhandler.SendMail(w, req)
	bytes, _ := ioutil.ReadAll(w.Result().Body)

	return string(bytes)
}
