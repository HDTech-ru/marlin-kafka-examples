package main

import (
	"crypto/tls"
	"github.com/IBM/sarama"
	"gitlab.vsk.ru/marlin/examples/marlin-kafka-examples/golang/provider"
	"log"
	"strings"
)

func main() {
	bootstrapServers := "kafka1.example.org:443,kafka2.example.org:443"
	clientID := "oauth-client-id"
	clientSecret := "oauth-client-secret"
	tokenEndpoint := "https://keycloak.example.org/auth/realms/kafka-authz/protocol/openid-connect/token"
	topic := "example"

	c := sarama.NewConfig()
	c.Net.TLS.Enable = true
	c.Net.TLS.Config = &tls.Config{
		InsecureSkipVerify: true,
	}
	c.Net.SASL.Enable = true
	c.Net.SASL.Mechanism = sarama.SASLTypeOAuth
	c.Net.SASL.TokenProvider = provider.NewTokenProvider(clientID, clientSecret, tokenEndpoint)
	c.Producer.RequiredAcks = sarama.WaitForLocal
	c.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(strings.Split(bootstrapServers, ","), c)
	if err != nil {
		panic(err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder("A test message"),
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		panic(err)
	}
	log.Printf("Message sent: partition=%d, offset=%d\n", partition, offset)
}
