package kafkian

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jwmwalrus/bnp"
	"github.com/jwmwalrus/felice-franz/base"
	log "github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// CreateConsumers creates consumers for the given environment topics
func CreateConsumers(env base.Environment, topics []string) (err error) {
	var c *kafka.Consumer

	c, err = kafka.NewConsumer(&env.Configuration)

	if err != nil {
		return
	}

	go func() {
		defer c.Close()

		c.SubscribeTopics(topics, nil)

		for {
			msg, err := c.ReadMessage(-1)
			if err != nil {
				log.Error(err)
				continue
			}
			fmt.Printf("%% Message on %s:\n%s\n",
				msg.TopicPartition, string(msg.Value))
			if msg.Headers != nil {
				fmt.Printf("%% Headers: %v\n", msg.Headers)
			}

			var payload []byte
			payload, err = getPayLoadFromMessage(msg)
			if err != nil {
				log.Error(err)
				continue
			}
			err = WS.WriteMessage(websocket.TextMessage, payload)
			bnp.LogOnError(err)
		}
	}()
	return
}

// CreateConsumerPoll creates consumers for the given environment topics
func CreateConsumerPoll(env base.Environment, topics []string) (err error) {
	var c *kafka.Consumer

	c, err = kafka.NewConsumer(&env.Configuration)

	if err != nil {
		return
	}

	go func() {
		defer c.Close()

		c.SubscribeTopics(topics, nil)

		for {
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}
			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				// But in this example we choose to terminate
				// the application if all brokers are down.
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					break
				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}()
	return
}

func getPayLoadFromMessage(m *kafka.Message) (payload []byte, err error) {
	type ht struct {
		Key, Value string
	}

	headers := []ht{}
	for _, h := range m.Headers {
		headers = append(headers, ht{h.Key, string(h.Value)})
	}

	flat := struct {
		Topic         string
		Partition     int32
		Offset        kafka.Offset
		Value         string
		Key           string
		Timestamp     time.Time
		TimestampType kafka.TimestampType
		Headers       []ht
	}{
		*m.TopicPartition.Topic,
		m.TopicPartition.Partition,
		m.TopicPartition.Offset,
		string(m.Value),
		string(m.Key),
		m.Timestamp,
		m.TimestampType,
		headers,
	}
	payload, err = json.Marshal(flat)
	return
}
