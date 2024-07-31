package main

import (
	// "flag"
	"fmt"
	"os"
	"log"
	"encode/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"https://github.com/elastic/go-sysinfo"
	"github.com/zcalusic/sysinfo"
)


type CPU struct {
	Processor       int32    `json:"processor"`
	VendorID        string   `json:"vendor_id"`
	CPUFamily       string   `json:"cpu_family"`
	Model           string   `json:"model"`
	ModelName       string   `json:"model_name"`
	Stepping        string   `json:"stepping"`
	Microcode       string   `json:"microcode"`
	CPUMHz          float32  `json:"cpu_mhz"`
	CacheSize       string   `json:"cache_size"`
	PhysicalID      int32    `json:"physical_id"`
	Siblings        int8     `json:"siblings"`
	CoreID          int32    `json:"core_id"`
	CPUCores        int32    `json:"cpu_cores"`
	APICID          int32    `json:"apicid"`
	InitialAPICID   int32    `json:"initial_apicid"`
	FPU             string   `json:"fpu"`
	FPUException    string   `json:"fpu_exception"`
	CPUIDLevel      string   `json:"cpuid_level"`
	WP              string   `json:"wp"`
	Flags           []string `json:"flags"`
	BogoMIPS        float32  `json:"bogomips"`
	Bugs            []string `json:"bugs"`
	CLFlushSize     uint16   `json:"clflush_size"`
	CacheAlignment  uint16   `json:"cache_alignment"`
	AddressSizes    []string `json:"address_sizes"`
	PowerManagement []string `json:"power_management"`
	TLBSize         string   `json:"tlb_size"`
}

type CPU struct {
	Vendor  string `json:"vendor,omitempty"`
	Model   string `json:"model,omitempty"`
	Speed   uint   `json:"speed,omitempty"`   // CPU clock rate in MHz
	Cache   uint   `json:"cache,omitempty"`   // CPU cache size in KB
	Cpus    uint   `json:"cpus,omitempty"`    // number of physical CPUs
	Cores   uint   `json:"cores,omitempty"`   // number of physical CPU cores
	Threads uint   `json:"threads,omitempty"` // number of logical (HT) CPU cores
}

type Kernel struct {
	Release      string `json:"release,omitempty"`
	Version      string `json:"version,omitempty"`
	Architecture string `json:"architecture,omitempty"`
}

type Memory struct {
	Type  string `json:"type,omitempty"`
	Speed uint   `json:"speed,omitempty"` // RAM data rate in MT/s
	Size  uint   `json:"size,omitempty"`  // RAM size in MB
}

type OS struct {
	Name         string `json:"name,omitempty"`
	Vendor       string `json:"vendor,omitempty"`
	Version      string `json:"version,omitempty"`
	Release      string `json:"release,omitempty"`
	Architecture string `json:"architecture,omitempty"`
}

type StorageDevice struct {
	Name   string `json:"name,omitempty"`
	Driver string `json:"driver,omitempty"`
	Vendor string `json:"vendor,omitempty"`
	Model  string `json:"model,omitempty"`
	Serial string `json:"serial,omitempty"`
	Size   uint   `json:"size,omitempty"` // device size in GB
}

type SysInfo struct {
	Meta    Meta            `json:"sysinfo"`
	Node    Node            `json:"node"`
	OS      OS              `json:"os"`
	Kernel  Kernel          `json:"kernel"`
	Product Product         `json:"product"`
	Board   Board           `json:"board"`
	Chassis Chassis         `json:"chassis"`
	BIOS    BIOS            `json:"bios"`
	CPU     CPU             `json:"cpu"`
	Memory  Memory          `json:"memory"`
	Storage []StorageDevice `json:"storage,omitempty"`
	Network []NetworkDevice `json:"network,omitempty"`
}

func main(){
	process, err := sysinfo.Self()
	if err != nil {
	return err
}


	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id": "test",
		"acks": "all",
	})
	if err != nil {
		fmt.Printf("Failed to create producer %s\n", err)
	}
	deliverch := make(chan kafka.Event, 10000)
	topic := "test-topic"
	err = p.Produce(&kafka.Message{
		TopicPatrition: kafka.TopicPatrition{topic: &topic, Partition: kafka.PartitionAny},
		Value:		[]byte("FOO"),
	},
	deliverch,
)
if err != nil {
	log.Fatal(err)
}
e := <-deliverch
fmt.Printf("%+v\n", e.String())
fmt.Printf("%+v\n", p)
}