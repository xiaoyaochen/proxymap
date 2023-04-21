package db

import (
	"log"
	"net/http"
	"net/url"
)

type StoreInfo struct {
	Url        string      `json:"url"`
	Host       string      `json:"host"`
	Port       string      `json:"port"`
	Method     string      `json:"method"`
	Scheme     string      `json:"scheme"`
	ReqBody    string      `json:"req_body"`
	ReqHeaders http.Header `json:"-"`

	StatusCode  int         `json:"status_code"`
	Title       string      `json:"title"`
	Length      int         `json:"lenght"`
	RespHeaders http.Header `json:"resp_headers"`
	RespBody    string      `json:"resp_body"`
	ContentType string      `json:"content_type"`
	Extension   string      `json:"extention"`
}

type DB interface {
	Push(string, []byte) error
	Close() error
}

func NewMqProducer(name string, addr string) DB {
	var db DB
	schema := GetSchema(addr)
	switch schema {
	case "redis":
		db = NewRedisProducer(name, addr)
		return db
	case "amqp":
		db = NewRbProducer(name, addr)
		return db
	case "http":
		db = NewESProducer(name, addr)
		return db
	case "https":
		db = NewESProducer(name, addr)
		return db
	case "mongo":
		db = NewMongoProducer(name, addr)
		return db
	default:
		return nil
	}
}

func GetSchema(Url string) string {
	u, err := url.Parse(Url)
	if err != nil {
		log.Println(err)
		return ""
	}
	return u.Scheme
}
