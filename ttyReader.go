package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/ian-kent/go-log/appenders"
	"github.com/ian-kent/go-log/layout"
	"github.com/ian-kent/go-log/log"
	"github.com/tarm/serial"
)

func main() {
	device := flag.String("device", "/dev/ttyIR01", "IR read/write head")
	flag.Parse()
	// fmt.Println("Trying connecting to", *device)
	config := &serial.Config{
		Name:        *device,
		Baud:        300,
		ReadTimeout: 1,
		Size:        7,
		Parity:      serial.ParityEven,
		StopBits:    serial.Stop1,
	}
	s, err := serial.OpenPort(config)
	if err != nil {
		log.Println("Could not open port.")
		log.Fatal(err)
	}

	// sending inital sequence
	_, err = s.Write([]byte("/?!\r\n"))
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(500 * time.Millisecond)

	reader := bufio.NewReader(s)
	for {
		reply, err := reader.ReadBytes('\n')
		if err != nil { // At the end, err will equal io.EOF
			if err != io.EOF {
				log.Println(reply, string(reply))
				log.Println(err)
			}
			break
		}
	}

	// requesting baud rate
	// \x06050\r\n = 9600
	// \x06060\r\n = 19200
	_, err = s.Write([]byte("\x06060\r\n"))
	if err != nil {
		log.Fatal(err)
	}
	s.Close()

	time.Sleep(200 * time.Millisecond) // This sleep is documented and required.

	config = &serial.Config{
		Name:        *device,
		Baud:        19200, // should match the requested baud rate
		ReadTimeout: 1,
		Size:        7,
		Parity:      serial.ParityEven,
		StopBits:    serial.Stop1,
	}
	s, err = serial.OpenPort(config)
	if err != nil {
		log.Println("Could not open port.")
		log.Fatal(err)
	}

	initLogger()
	for {
		read(s)
	}
}

func initLogger() {
	logger := log.Logger()
	logger.SetAppender(appenders.RollingFile("smartmeter.log", true))
	appender := logger.Appender()
	appender.SetLayout(layout.Pattern("%d %p - %m%n"))
}

func read(s *serial.Port) {
	reader := bufio.NewReader(s)
	for {
		reply, err := reader.ReadBytes('\n')
		if err != nil { // At the end, err will equal io.EOF
			if err != io.EOF {
				log.Println(err)
			}
			break
		}
		line := string(reply)
		line = strings.Replace(line, "\n", "", -1)
		line = strings.TrimSpace(line)
		log.Println(line)
		grep(line)
	}
}

func grep(line string) {
	m := make(map[string]string)
	v := strings.Split(line, "(")

	id := v[0]

	if strings.Index(id, "*") > 0 {
		r := strings.Split(id, "*")
		id = r[0]
		m["channel"] = r[1]
	}
	m["id"] = id
	vl := strings.Split(v[1], ")")
	value := vl[0]
	if strings.Index(value, "*") > 0 {
		r := strings.Split(value, "*")
		value = r[0]
		m["unit"] = r[1]
	}
	m["value"] = value
	if len(v) > 2 {
		vl = strings.Split(v[2], ")")
		datestring := vl[0]
		layout := "06-01-02 15:04"
		t, err := time.Parse(layout, datestring)
		if err != nil {
			log.Println("Failed to parse month date", err)
		} else {
			m["time"] = t.Format(time.RFC3339Nano)
		}

	} else {
		t := time.Now()
		m["time"] = t.Format(time.RFC3339Nano)
	}

	jsonStr, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		fmt.Println(string(jsonStr))
	}

}
