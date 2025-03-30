package main

import (
	"context"
	"os"
	"time"

	kafka "github.com/segmentio/kafka-go"
	"github.com/smukherj1/k8s-signer/pkg/log"
)

func main() {
	log.Init()
	log.Info("Launching worker.")
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"redpanda-0.redpanda.redpanda.svc.cluster.local:9093", "redpanda-1.redpanda.redpanda.svc.cluster.local:9093"},
		GroupID:  "consumer-group-test",
		Topic:    "test-topic",
		MaxBytes: 10e6,
	})
	ctx := context.Background()
	errors := 0
	errorLimit := 10
	for {
		if errors > errorLimit {
			break
		} else if errors > 0 {
			log.Infof("Sleeping for %vs due to %v consecutive errors.", errors, errors)
			time.Sleep(time.Duration(errors) * time.Second)
		}
		m, err := r.FetchMessage(ctx)
		if err != nil {
			log.Errorf("Error fetching message from Kafka: %v", err)
			errors++
			continue
		}
		log.Infof("Got message [%v] Topic %v, Partition %v, Offset %v, Value %v.", m.Time, m.Topic, m.Partition, m.Offset, string(m.Value))
		if err := r.CommitMessages(ctx, m); err != nil {
			log.Errorf("Error committing message [%v] Topic %v, Partition %v, Offset %v.", m.Time, m.Topic, m.Partition, m.Offset)
			errors++
			continue
		}
		if errors > 0 {
			log.Infof("Resetting error count from %v back to 0.", errors)
			errors = 0
		}
	}
	log.Error("Abnormal exit due to", errorLimit, "consecutive errors.")
	os.Exit(1)
}
