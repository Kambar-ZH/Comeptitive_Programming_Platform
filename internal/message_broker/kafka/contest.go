package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"site/internal/datastruct"
)

const (
	topic = "contest"
)

type (
	ContestBroker struct {
		producer      sarama.SyncProducer
		consumerGroup sarama.ConsumerGroup

		consumeHandler contestConsumeHandler
		clientId       string
	}

	contestConsumeHandler struct {
		ready chan bool
	}
)

func NewContestBroker(clientId string) *ContestBroker {
	return &ContestBroker{
		clientId: clientId,
		consumeHandler: contestConsumeHandler{
			ready: make(chan bool),
		},
	}
}

func (c *ContestBroker) Connect(ctx context.Context, brokers []string) error {
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.Return.Successes = true
	producerConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producer, err := sarama.NewSyncProducer(brokers, producerConfig)
	if err != nil {
		panic(err)
	}
	c.producer = producer

	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Return.Errors = true
	consumerGroup, err := sarama.NewConsumerGroup(brokers, c.clientId, consumerConfig)
	if err != nil {
		panic(err)
	}
	c.consumerGroup = consumerGroup

	go func() {
		for {
			err := c.consumerGroup.Consume(ctx, []string{topic}, c.consumeHandler)
			if err != nil {
				log.Println(err)
			}
			if ctx.Err() != nil {
				return
			}
			c.consumeHandler.ready = make(chan bool)
		}
	}()
	<-c.consumeHandler.ready

	return nil
}

func (c *ContestBroker) CreateContest(contest *datastruct.Contest) error {
	msg := &datastruct.ContestMessage{
		Command: datastruct.ContestCommandCreate,
		Contest: contest,
	}

	msgRaw, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	contestMsg := new(datastruct.ContestMessage)
	_ = json.Unmarshal(msgRaw, contestMsg)
	fmt.Println(contestMsg)

	_, _, err = c.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msgRaw),
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *ContestBroker) Close() error {
	if err := c.producer.Close(); err != nil {
		return err
	}

	if err := c.consumerGroup.Close(); err != nil {
		return err
	}

	return nil
}

func (c contestConsumeHandler) Setup(session sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

func (c contestConsumeHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c contestConsumeHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		contestMsg := new(datastruct.ContestMessage)
		if err := json.Unmarshal(msg.Value, contestMsg); err != nil {
			return err
		}
		switch contestMsg.Command {
		case datastruct.ContestCommandCreate:
			fmt.Println(contestMsg.Contest)
		default:
			fmt.Println("Undefined command")
		}
		session.MarkMessage(msg, "")
	}
	return nil
}
