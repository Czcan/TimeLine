package test

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/Czcan/TimeLine/config"
	"github.com/Czcan/TimeLine/internal/models"
	"github.com/Czcan/TimeLine/internal/server"
	"github.com/ddliu/go-httpclient"
	"github.com/go-chi/chi"
	"github.com/tidwall/gjson"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	err    error
	Server *httptest.Server
	Router *chi.Mux
)

type Record struct {
	Value string
}

func setup() {
	c := config.MustGetAppConfig()
	DB, err = gorm.Open(mysql.Open(c.DB), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB.Migrator().DropTable(&models.User{}, &models.Folder{}, &models.Note{},
		&models.Collection{}, &models.Account{}, &models.Comment{}, &models.Liker{}, &models.Task{})
	DB.AutoMigrate(&models.User{}, &models.Folder{}, &models.Note{},
		&models.Collection{}, &models.Account{}, &models.Comment{}, &models.Liker{}, &models.Task{})

	createdData()

	jwtClient := &JWTClientMock{}
	Router = server.New(DB, jwtClient)
	Server = httptest.NewServer(Router)
}

func Get(url string) string {
	client := httpclient.NewHttpClient()
	resp, _ := client.Get(Server.URL + url)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func SingeGet(token string, url string, values url.Values) string {
	client := httpclient.NewHttpClient()
	client.Headers = make(map[string]string)
	client.Headers["Authorization"] = token
	params := make(map[string]string)
	for index, value := range values {
		params[index] = value[0]
	}
	resp, _ := client.Get(Server.URL+url, params)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func Post(url string, values url.Values) string {
	client := httpclient.NewHttpClient()
	params := make(map[string]string)
	for index, value := range values {
		params[index] = value[0]
	}
	resp, _ := client.Post(Server.URL+url, params)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func SingePost(token string, url string, values url.Values) string {
	client := httpclient.NewHttpClient()
	client.Headers = make(map[string]string)
	client.Headers["Authorization"] = token
	params := make(map[string]string)
	for index, value := range values {
		params[index] = value[0]
	}
	resp, _ := client.Post(Server.URL+url, params)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func SingePutJson(token string, url string, values url.Values) string {
	client := httpclient.NewHttpClient()
	client.Headers = make(map[string]string)
	client.Headers["Authorization"] = token
	params := make(map[string]string)
	for index, value := range values {
		params[index] = value[0]
	}
	resp, _ := client.PutJson(Server.URL+url, params)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func SingePut(token string, url string, values url.Values) string {
	client := httpclient.NewHttpClient()
	client.Headers = make(map[string]string)
	client.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	client.Headers["Authorization"] = token
	params := strings.NewReader(values.Encode())
	resp, _ := client.Put(Server.URL+url, params)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func SingeDelete(token string, url string, values url.Values) string {
	client := httpclient.NewHttpClient()
	client.Headers = make(map[string]string)
	client.Headers["Authorization"] = token
	params := make(map[string]string)
	for index, value := range values {
		params[index] = value[0]
	}
	resp, _ := client.Delete(Server.URL+url, params)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func createdData() {
	RunSQL(DB, `
		INSERT INTO users (id, email, uid, password, nick_name) VALUES (1, 'test1@qq.com', '123123','123456','name');
		INSERT INTO users (id, email, password) VALUES (2, 'test2@qq.com', '123456');
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

func ExtractDate(body string, filedName string, new string) string {
	values := gjson.Get(body, filedName).Array()
	for _, value := range values {
		if value.String() == "0" {
			continue
		}
		body = strings.Replace(body, value.String(), new, -1)
	}
	return body
}
