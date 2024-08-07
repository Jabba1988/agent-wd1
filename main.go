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

	 "encoding/json"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type sysinfo struct {
	Cpu          []string `json:"cpu"`
	Fqdn         []string `json:"fqdn"`
	Memory       []string `json:"memory"`
	Swap         []string `json:"swap"`
	Mount_points []string `json:"mp"`
	Ip           []string `json:"ip"`
}

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
	s := sysinfo{Cpu: get_cpu_info()}
	s.Memory, s.Swap = get_ram()
	s.Mount_points = get_mount_points()
	s.Fqdn = get_fqdn()
	s.Ip = get_ips()
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
	// var sslCert string
	// var sslCertKey string 
	flag.StringVar(&svar, "bootstrap", "localhost:9092", "kafka host with port")
	// flag.StringVar(&sslCert, "cert", "", "path to cert.pem")
	// flag.StringVar(&sslCertKey, "cert-key", "", "path to key.pem")
	flag.Parse()
	config := &kafka.ConfigMap{
		"bootstrap.servers": svar,
		// "ssl.certificat.pem": sslCert,
		// "ssl.key.pem": sslCertKey,
	}
	producer, err := kafka.NewProducer(config)
	if err != nil {
		log.Fatal(err)
	}
	topic := "awd"
	for true {
		//  value := fmt.Sprintf("message-%d", i)
		// value := []byte(fmt.Sprintf("%v", payloadStruct()))
		value := payloadStruct()
		msg, _ := json.Marshal(value)
		err := producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny} ,
			Value: []byte(msg),
		}, nil)
		if err != nil {
			log.Fatal("Failed to produce message %d: %v\n", err)	
		} else {
			fmt.Printf("Produced message %s  %s\n", time.Now(), value)
		}
		time.Sleep(10 * time.Second)
	}
}