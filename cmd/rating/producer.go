package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"main/rating/model"
	"main/util"
	"os"

	"github.com/apache/pulsar-client-go/pulsar"
)

func main() {
	var config string
	var data string
	flag.StringVar(&config, "config", ".env", "Configuration path")
	flag.StringVar(&data, "data", "data.json", "Load Rating data")
	flag.Parse()
	cfg := util.LoadConfig(config)

	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               cfg.PulsarURL,
		ConnectionTimeout: cfg.ConnectionTimeout,
		OperationTimeout:  cfg.OperationTimeout,
	})
	if err != nil {
		log.Fatal("failed to create pulsar client:", err)
	}
	defer client.Close()

	fmt.Println("Creating a Pulsar producer")

	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: cfg.TopicName,
	})
	if err != nil {
		log.Fatal("failed to create producer:", err)
	}
	defer producer.Close()

	fmt.Println("Reading rating events from file:", data)
	ratingEvents, err := readRatingEvents(data)
	if err != nil {
		log.Fatal("failed to read events from file:", err)
	}

	if err := produceRatingEvents(producer, cfg.TopicName, ratingEvents); err != nil {
		log.Fatal("failed produce rating events:", err)
	}

	if err := producer.Flush(); err != nil {
		panic(err)
	}
}

func readRatingEvents(fileName string) ([]model.RatingEvent, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var ratings []model.RatingEvent
	if err := json.NewDecoder(f).Decode(&ratings); err != nil {
		return nil, err
	}
	return ratings, nil
}

func produceRatingEvents(producer pulsar.Producer, topicName string, ratingEvents []model.RatingEvent) error {
	for _, ratingEvent := range ratingEvents {
		encodedEvent, err := json.Marshal(ratingEvent)
		if err != nil {
			return err
		}

		msgId, err := producer.Send(context.Background(), &pulsar.ProducerMessage{
			Payload: encodedEvent,
		})
		if err != nil {
			return err
		}
		fmt.Printf("Produced event to topic %s: msgId = %s, value = %s\n", topicName, msgId.String(), string(encodedEvent))
	}
	return nil
}
