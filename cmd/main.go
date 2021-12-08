package main

import (
	"database/sql"
	"encoding/json"
	"github.com/camila-costa/imersao-gateway/adapter/broker/kafka"
	"github.com/camila-costa/imersao-gateway/adapter/factory"
	"github.com/camila-costa/imersao-gateway/adapter/presenter/transaction"
	"github.com/camila-costa/imersao-gateway/usecase/process_transaction"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	repositoryFactory := factory.NewRepositoryDatabaseFactory(db)
	repository := repositoryFactory.CreateTransactionRepository()
	configMapProducer := &ckafka.ConfigMap{
		"bootstrap.servers": "host.docker.internal:9094",
	}
	kafkaPresenter := transaction.NewTransactionKafkaPresenter()
	producer := kafka.NewKafkaProducer(configMapProducer, kafkaPresenter)

	var msgChan = make(chan *ckafka.Message)
	configMapConsumer := &ckafka.ConfigMap{
		"bootstrap.servers": "host.docker.internal:9094",
		"client.id":         "goapp",
		"group.id":          "goapp",
	}
	topics := []string{"transactions"}
	consumer := kafka.NewConsumer(configMapConsumer, topics)
	go consumer.Consume(msgChan)

	usecase := process_transaction.NewProcessTransaction(repository, producer, "transactions_result")

	for msg := range msgChan {
		var input process_transaction.TransactionDtoInput
		json.Unmarshal(msg.Value, &input)
		usecase.Execute(input)
	}
}