package db

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoProducer struct {
	name        string
	logger      *log.Logger
	connection  *mongo.Client
	isConnected bool
	rdbCtx      context.Context
	collection  *mongo.Collection
}

func NewMongoProducer(name string, addr string) *MongoProducer {
	producer := MongoProducer{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		name:   name,
		rdbCtx: context.TODO(),
	}
	client, err := mongo.Connect(producer.rdbCtx, options.Client().ApplyURI(addr))
	if err != nil {
		panic(err)
	}
	if err := client.Ping(producer.rdbCtx, readpref.Primary()); err != nil {
		panic(err)
	}
	producer.connection = client
	producer.isConnected = true
	producer.collection = client.Database("enterprise").Collection(name)
	return &producer
}

func (producer *MongoProducer) Push(docid string, doc []byte) error {
	if !producer.isConnected {
		return errors.New("failed to push push: not connected")
	}
	filter := bson.M{"_id": docid}

	var obj map[string]interface{}
	if err := json.Unmarshal(doc, &obj); err != nil {
		panic(err)
	}
	var update_map map[string]interface{}
	update_map["$set"] = obj
	update := bson.M(update_map)

	opts := options.Update().SetUpsert(true)

	_, err := producer.collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		producer.logger.Println(err)
	}
	return nil
}

func (producer *MongoProducer) Close() error {
	producer.isConnected = false
	err := producer.connection.Disconnect(context.Background())
	if err != nil {
		return err
	}
	return nil
}
