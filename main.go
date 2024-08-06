package main

import (
	"flag"
	"fmt"
	"time"
	// "github.com/IBM/sarama"
	"log"
	//"os"
	"os/exec"
	"regexp"
	"strings"

	// "encoding/json"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type sysinfo struct {
	fqdn         []string
	cpu          []string
	memory       []string
	swap         []string
	mount_points []string
	ip           []string
}

// type ProducerMessage struct {
// 	Topic string // The Kafka topic for this message.
// 	// The partitioning key for this message. Pre-existing Encoders include
// 	// StringEncoder and ByteEncoder.
// 	Key Encoder
// 	// The actual message to store in Kafka. Pre-existing Encoders include
// 	// StringEncoder and ByteEncoder.
// 	Value Encoder

// 	// The headers are key-value pairs that are transparently passed
// 	// by Kafka between producers and consumers.
// 	Headers []RecordHeader

// 	// This field is used to hold arbitrary data you wish to include so it
// 	// will be available when receiving on the Successes and Errors channels.
// 	// Sarama completely ignores this field and is only to be used for
// 	// pass-through data.
// 	Metadata interface{}

// 	// Offset is the offset of the message stored on the broker. This is only
// 	// guaranteed to be defined if the message was successfully delivered and
// 	// RequiredAcks is not NoResponse.
// 	Offset int64
// 	// Partition is the partition that the message was sent to. This is only
// 	// guaranteed to be defined if the message was successfully delivered.
// 	Partition int32
// 	// Timestamp can vary in behavior depending on broker configuration, being
// 	// in either one of the CreateTime or LogAppendTime modes (default CreateTime),
// 	// and requiring version at least 0.10.0.
// 	//
// 	// When configured to CreateTime, the timestamp is specified by the producer
// 	// either by explicitly setting this field, or when the message is added
// 	// to a produce set.
// 	//
// 	// When configured to LogAppendTime, the timestamp assigned to the message
// 	// by the broker. This is only guaranteed to be defined if the message was
// 	// successfully delivered and RequiredAcks is not NoResponse.
// 	Timestamp time.Time
// 	// contains filtered or unexported fields
// }



func init() {

}



func get_cpu_info() []string {
	cpu, err := exec.Command("lscpu").Output()
	if err != nil {
		log.Fatal(err)
	}
	var a string = string(cpu)
	b := strings.Split(a, "\n")

	return b
}

func get_ram() ([]string, []string) {
	ram, err := exec.Command("free").Output()
	if err != nil {
		log.Fatal(err)
	}
	r, _ := regexp.Compile(`(?:Mem:         +\d*)`)
	sw, _ := regexp.Compile(`(?:Swap:        +\d*)`)
	var a = r.FindAllString(string(ram), 100)
	var b = sw.FindAllString(string(ram), 100)
	return a, b
}
func get_fqdn() []string {
	fqdn, err := exec.Command("hostname",).Output()
	if err != nil {
		log.Fatal(err)
	}
	var a string = string(fqdn)
	b := strings.Split(a, "\n")
	return b
}
func get_mount_points() []string {
	mp, err := exec.Command("df", "-h").Output()
	if err != nil {
		log.Fatal(err)
	}
	r, _ := regexp.Compile(`(?:/dev/mapper/[a-z,_,-]* *\d*\S *\d*\S\d*\S+\s*\d*\S\d*\S\s*\d*\S*\s*\S*)`)
	mps := r.FindAllString(string(mp), 100)
	var a []string = []string(mps)
	return a
}
func get_ips() []string {
	ips, err := exec.Command("ip", "a").Output()
	if err != nil {
		log.Fatal(err)
	}
	r, _ := regexp.Compile(`(?:inet\s*\d*\S\d*\S\d*\S\d*)`)
	var a []string = r.FindAllString(string(ips), 100)
	return a
}
func payloadStruct() *sysinfo {
	s := sysinfo{cpu: get_cpu_info()}
	s.memory, s.swap = get_ram()
	s.mount_points = get_mount_points()
	s.fqdn = get_fqdn()
	s.ip = get_ips()
	return &s
}

type TopicPartition struct {
	Topic       *string
	Partition   int32
	Offset      int
	Metadata    *string
	Error       error
	LeaderEpoch *int32 // LeaderEpoch or nil if not available
}

func main() {
	var svar string 
	flag.StringVar(&svar, "bootstrap", "localhost:9092", "kafka host with port")
	flag.Parse()
	config := &kafka.ConfigMap{
		"bootstrap.servers": svar,
	}
	producer, err := kafka.NewProducer(config)
	if err != nil {
		log.Fatal(err)
	}
	topic := "awd"
	for true {
		//  value := fmt.Sprintf("message-%d", i)
		value := []byte(fmt.Sprintf("%v", payloadStruct()))
		err := producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny} ,
			Value: []byte(value),
		}, nil)
		if err != nil {
			log.Fatal("Failed to produce message %d: %v\n", err)	
		} else {
			fmt.Printf("Produced message %s  %s\n", time.Now(), value)
		}
		time.Sleep(10 * time.Second)
	}
}