package main

import (
	"context"
	"crypto/tls"
	"github.com/IBM/sarama"
	"gitlab.vsk.ru/marlin/examples/marlin-kafka-examples/golang/provider"
	"log"
	"strings"
	"time"
)

// struct defining the handler for the consuming Sarama method
type consumerGroupHandler struct {
	msgContent string
}

func (cgh *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (cgh *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (cgh *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(20 * time.Second)
		timeout <- true
		close(timeout)
	}()

	log.Printf("Listening topic `%s`, partition `%v`, offset `%v`", claim.Topic(), claim.Partition(), claim.InitialOffset())

	select {
	case <-timeout:
		log.Println("timout waiting for new message")
	case message := <-claim.Messages():
		if message != nil {
			cgh.msgContent = string(message.Value)
			log.Printf("Message received: value=%s, partition=%d, offset=%d", cgh.msgContent, message.Partition, message.Offset)
			session.MarkMessage(message, "")
		}
	}

	return nil
}

func main() {
	bootstrapServers := "kafka1.example.org:443,kafka2.example.org:443,"
	clientID := "oauth-client-id"
	clientSecret := "oauth-client-secret"
	tokenEndpoint := "https://keycloak.example.org/auth/realms/kafka-authz/protocol/openid-connect/token"
	topic := "example"
	groupID := "example"

	c := sarama.NewConfig()
	c.Net.TLS.Enable = true
	c.Net.SASL.Enable = true
	c.Net.SASL.Mechanism = sarama.SASLTypeOAuth
	c.Net.SASL.TokenProvider = provider.NewTokenProvider(clientID, clientSecret, tokenEndpoint)
	c.Net.TLS.Config = &tls.Config{
		InsecureSkipVerify: true,
	}
	c.Consumer.Offsets.AutoCommit.Enable = true
	c.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(strings.Split(bootstrapServers, ","), groupID, c)
	if err != nil {
		panic(err)
	}

	cgh := &consumerGroupHandler{}
	ctx := context.Background()
	// this method calls the methods handler on each stage: setup, consume and cleanup
	err = consumerGroup.Consume(ctx, []string{topic}, cgh)
	if err != nil {
		panic(err)
	}

	err = consumerGroup.Close()
	if err != nil {
		panic(err)
	}

	log.Println(cgh.msgContent)
}
