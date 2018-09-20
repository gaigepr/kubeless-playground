package kubeless

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kubeless/kubeless/pkg/functions"
	kafka "github.com/segmentio/kafka-go"
)

var (
	kafkaHost  string = os.Getenv("KAFKA_HOST")
	kafkaPort  string = os.Getenv("KAFKA_PORT")
	kafkaTopic string = os.Getenv("KAFKA_TOPIC")
)

func generateMessage(event functions.Event) string {
	return fmt.Sprintf(`{ "id":"%s" , "time":"%s" , "data":"%s" }`, event.EventID, event.EventTime, event.Data)
}

func Handler(event functions.Event, ctx functions.Context) (string, error) {
	topic := kafkaTopic
	partition := 0

	conn, _ := kafka.DialLeader(context.Background(), "tcp", kafkaHost+":"+kafkaPort, topic, partition)
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	message := generateMessage(event)
	log.Println(message)
	conn.WriteMessages(
		kafka.Message{Value: []byte(message)},
	)

	conn.Close()
	return "ok", nil
}
