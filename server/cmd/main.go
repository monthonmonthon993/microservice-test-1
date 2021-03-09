package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

var count = 0

func main() {
	// produce to kafka
	ctx := context.Background()
	// go consume(ctx)
	r := gin.Default()
	r.POST("/greeter", func(c *gin.Context) {
		greeter := Greeter{}

		if err := c.ShouldBindJSON(&greeter); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		now := time.Now()
		fmt.Println(count)
		if count == 0 {
			produce(ctx, greeter, now, topicA)
		} else if count == 1 {
			produce(ctx, greeter, now, topicB)
		} else {
			produce(ctx, greeter, now, topicC)
			count = -1
		}
		count++
		c.JSON(200, gin.H{
			"code":          "OK",
			"received_time": now,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

type Greeter struct {
	ID     int    `json:"msg_id"`
	Sender string `json:"sender"`
	MSG    string `json:"msg"`
}

// the topic and broker address are initialized as constants
const (
	topicA         = "topicA"
	topicB         = "topicB"
	topicC         = "topicC"
	broker1Address = "localhost:9091"
	// broker2Address = "localhost:9092"
	// broker3Address = "localhost:9093"
)

func produce(ctx context.Context, greeter Greeter, now time.Time, topic string) {

	// intialize the writer with the broker addresses, and the topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker1Address},
		Topic:   topic,
	})

	// each kafka message has a key and value. The key is used
	// to decide which partition (and consequently, which broker)
	// the message gets published on
	err := w.WriteMessages(ctx, kafka.Message{
		Key: []byte(strconv.Itoa(greeter.ID) + now.Format(time.RFC3339)),
		// create an arbitrary message payload for the value
		Value: []byte(fmt.Sprintf(`%d %s: %s %s`, greeter.ID, greeter.Sender, greeter.MSG, now.Format(time.RFC3339))),
		Time:  now,
	})
	if err != nil {
		panic("could not write message " + err.Error())
	}

	// log a confirmation once the message is written
	fmt.Println("writes:", count)
	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
