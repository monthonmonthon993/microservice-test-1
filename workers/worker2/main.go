package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	db := connectToMongoDB()
	consume(ctx, db)
}

const (
	topicB         = "topicB"
	broker1Address = "localhost:9091"
)

func consume(ctx context.Context, db *mongo.Database) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker1Address},
		Topic:   topicB,
		GroupID: "my-group-B",
		GroupBalancers: []kafka.GroupBalancer{
			kafka.RoundRobinGroupBalancer{},
		}})
	// the `ReadMessage` method blocks until we receive the next event
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("worker1: message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = db.Collection("greeter").InsertOne(ctx, bson.D{{"key", string(m.Key)}, {"value", string(m.Value)}, {"time", m.Time}})
		if err != nil {
			fmt.Println(err)
			panic("cannot insert")
		}
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}

func NewDB(client *mongo.Client) *mongo.Database {
	return client.Database("db-test")
}

func connectToMongoDB() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://username:secret@localhost:27017/?authSource=admin&readPreference=primary&appname=MongoDB%20Compass&ssl=false"))
	if err != nil {
		panic(err)
	}
	fmt.Println("connect to mongodb successfully")
	return NewDB(client)
}
