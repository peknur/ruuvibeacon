package ruuvibeacon

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/peknur/ruuvitag"
)

// Publisher type
type Publisher = func(ctx context.Context, e *Envelope)

// Publishers hold functions to publish envelope
var Publishers = make(map[string]Publisher)

// Reading represents Readinged data
type Reading struct {
	DeviceID        string
	Version         string
	Humidity        float32
	Temperature     float32
	Pressure        uint32
	AccelerationX   float32
	AccelerationY   float32
	AccelerationZ   float32
	BatteryVoltage  float32
	TXPower         int8
	MovementCounter uint8
	Timestamp       time.Time
}

// Envelope represent Beacon data in some point in time
type Envelope struct {
	Host    string
	Tick    int
	Time    time.Time
	Started time.Time
	Data    []Reading
}

// Beacon object
type Beacon struct {
	hostname   string
	scanner    ruuvitag.Scanner
	data       map[string]Reading
	mu         sync.RWMutex
	started    time.Time
	tick       int
	publishers []Publisher
}

// WebViewHandler is http handler
func (b *Beacon) WebViewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	js, err := json.MarshalIndent(b, "", " ")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(js)
}

// Readings returns latest device readings as array
func (b *Beacon) Readings() []Reading {
	b.mu.RLock()
	defer b.mu.RUnlock()

	d := make([]Reading, len(b.data))
	i := 0
	for _, v := range b.data {
		d[i] = v
		i++
	}
	return d
}

// EncodeEnvelope creates envelope for publisher
func (b *Beacon) EncodeEnvelope() *Envelope {
	return &Envelope{
		Host:    b.hostname,
		Tick:    b.tick,
		Time:    time.Now(),
		Started: b.started,
		Data:    b.Readings(),
	}
}

// MarshalJSON is encodes an envelope
func (b *Beacon) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.EncodeEnvelope())
}

// Publisher publishes envelopes in intervals set by tick
func (b *Beacon) Publisher() {
	ticker := time.NewTicker(time.Duration(b.tick) * time.Second)
	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(b.tick)*time.Second)
			defer cancel()
			for _, p := range b.publishers {
				go p(ctx, b.EncodeEnvelope())
			}
		}
	}
}

// Scan starts scanner and updates beacon readings
func (b *Beacon) Scan() {
	output := b.scanner.Start()
	for {
		data, ok := <-output
		if ok == false {
			log.Println("scanner closed channel")
			return
		}
		b.setData(data.DeviceID(), newReading(data))
	}
}

func (b *Beacon) setData(key string, value Reading) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.data[key] = value
}

// newReading from Measurement
func newReading(m ruuvitag.Measurement) Reading {
	return Reading{
		m.DeviceID(),
		fmt.Sprintf("%d", m.Format()),
		m.Humidity(),
		m.Temperature(),
		m.Pressure(),
		m.AccelerationX(),
		m.AccelerationY(),
		m.AccelerationZ(),
		m.BatteryVoltage(),
		m.TXPower(),
		m.MovementCounter(),
		m.Timestamp()}
}

func loadOuputs(outputs string) []Publisher {
	if outputs == "" {
		log.Panicln("invalid -output flag value")
	}
	selected := strings.Split(outputs, ",")
	p := make([]Publisher, 0)
	for _, s := range selected {
		f, ok := Publishers[s]
		if ok {
			log.Printf("load '%s' publisher", s)
			p = append(p, f)
		}
	}
	return p
}

// Run ruuvibeacon
func Run() {
	var tick int
	var scannerBufferSize int
	var httpPort string
	var outputs string
	flag.IntVar(&tick, "tick", 60, "Data sending interval (seconds)")
	flag.IntVar(&scannerBufferSize, "buffer", 10, "Scanner buffer size. How many measurements are buffered.")
	flag.StringVar(&httpPort, "port", "8080", "httpd server port")
	flag.StringVar(&outputs, "output", "log", "comma separated list of outputs")
	flag.Parse()

	scanner, err := ruuvitag.OpenScanner(scannerBufferSize)
	if err != nil {
		log.Fatal(err)
	}
	defer scanner.Stop()
	hostname, _ := os.Hostname()
	beacon := Beacon{
		scanner:    scanner,
		data:       make(map[string]Reading),
		started:    time.Now(),
		tick:       tick,
		hostname:   hostname,
		publishers: loadOuputs(outputs),
	}
	go beacon.Scan()
	go beacon.Publisher()

	http.HandleFunc("/", beacon.WebViewHandler)
	<-httpListenAndServer(":" + httpPort)
	log.Println("shutting down ..")
}
