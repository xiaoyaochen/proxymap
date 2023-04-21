package db

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ESProducer struct {
	name        string
	logger      *log.Logger
	connection  *elasticsearch.Client
	isConnected bool
	rdbCtx      context.Context
}

func NewESProducer(name string, addr string) *ESProducer {
	producer := ESProducer{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		name:   name,
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cfg := elasticsearch.Config{
		Addresses: []string{
			addr,
		},
		Transport: tr,
	}
	eslient, _ := elasticsearch.NewClient(cfg)
	_, err := eslient.Ping()
	if err != nil {
		// 连接失败
		log.Fatalf("Error connecting to the server: %s", err.Error())
	}
	log.Println("Attempting to connect")
	producer.connection = eslient
	producer.isConnected = true
	return &producer
}

func (producer *ESProducer) Push(docid string, doc []byte) error {
	if !producer.isConnected {
		return errors.New("failed to push push: not connected")
	}
	req := esapi.IndexRequest{
		Index:      producer.name,
		DocumentID: docid,
		Body:       bytes.NewReader(doc),
		Refresh:    "true",
	}
	res, _ := req.Do(context.Background(), producer.connection)
	defer res.Body.Close()
	return nil
}

func (producer *ESProducer) Close() error {
	producer.isConnected = false
	return nil
}
