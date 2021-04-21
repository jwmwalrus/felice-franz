package kafkian

import (
	"github.com/jwmwalrus/felice-n-franz/internal/base"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func ProduceMessage(env base.Environment, msg *kafka.Message) (err error) {
	var p *kafka.Producer
	ec := kafka.ConfigMap{}
	for k, v := range env.Configuration {
		if k == "group.id" {
			continue
		}
		ec[k] = v
	}
	p, err = kafka.NewProducer(&ec)
	if err != nil {
		return
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				toast := toastMsg{
					ToastType: "info",
					Title:     "From Producer",
					Message:   "Message delivered",
					Topic:     *ev.TopicPartition.Topic,
					Partition: ev.TopicPartition.Partition,
					Offset:    ev.TopicPartition.Offset,
				}
				if ev.TopicPartition.Error != nil {
					toast.ToastType = "error"
					toast.Message = "Delivery failed"
				}
				sendToast(toast)
			}
		}
	}()

	p.Produce(msg, nil)

	// Wait for message deliveries before shutting down
	p.Flush(60 * 1000)
	return
}