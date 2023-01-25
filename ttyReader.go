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

var myMap = map[string]string{
	"0.0.0":     "Serialnummer ",
	"0.0.1":     "SAT Seriennummer ",
	"1.128.0":   "Inkasso-Zählwerk ",
	"0.1.0":     "Kumulierungszähler ",
	"0.1.2.99":  "Uhrzeit (ÜF) ",
	"1.2.0":     "Kum. Max A+ ",
	"2.2.0":     "Kum. Max A- ",
	"1.4.0":     "lauf. Max A+ (Minuten Wert) =>Zähler ",
	"2.4.0":     "lauf. Max A- (Minuten Wert) =>Zähler ",
	"1.6.0":     "aktuelles Max A+ Datum Uhrzeit ",
	"1.6.0.99":  "Max Vorperiode A+ Datum Uhrzeit ",
	"2.6.0":     "aktuelles Max A- Datum Uhrzeit ",
	"2.6.0.99":  "Max Vorperiode A- Datum Uhrzeit ",
	"1.8.0":     "Energie A+ Tariflos ",
	"1.8.0.99":  "Energie A+ Tariflos Vorwerte ",
	"1.8.1":     "Energie A+ Tarif 1 ",
	"1.8.1.99":  "Energie A+ Tarif 1 Vorwerte ",
	"1.8.2":     "Energie A+ Tarif 2 ",
	"1.8.2.99":  "Energie A+ Tarif 2 Vorwerte ",
	"1.8.3":     "Energie A+ Tarif 3 ",
	"1.8.3.99":  "Energie A+ Tarif 3 Vorwerte ",
	"1.8.4":     "Energie A+ Tarif 4 ",
	"1.8.4.99":  "Energie A+ Tarif 4 Vorwerte ",
	"1.8.5":     "Energie A+ Tarif 5 ",
	"1.8.5.99":  "Energie A+ Tarif 5 Vorwerte ",
	"1.8.6":     "Energie A+ Tarif 6 ",
	"1.8.6.99":  "Energie A+ Tarif 6 Vorwerte ",
	"2.8.0":     "Energie A- Tariflos ",
	"2.8.0.99":  "Energie A- Tariflos Vorwerte ",
	"2.8.1":     "Energie A- Tarif 1 ",
	"2.8.1.99":  "Energie A- Tarif 1 Vorwerte ",
	"2.8.2":     "Energie A- Tarif 2 ",
	"2.8.2.99":  "Energie A- Tarif 2 Vorwerte ",
	"2.8.3":     "Energie A- Tarif 3 ",
	"2.8.3.99":  "Energie A- Tarif 3 Vorwerte ",
	"2.8.4":     "Energie A- Tarif 4 ",
	"2.8.4.99":  "Energie A- Tarif 4 Vorwerte ",
	"2.8.5":     "Energie A- Tarif 5 ",
	"2.8.5.99":  "Energie A- Tarif 5 Vorwerte ",
	"2.8.6":     "Energie A- Tarif 6 ",
	"2.8.6.99":  "Energie A- Tarif 6 Vorwerte ",
	"3.8.1":     "Energie R+ ",
	"3.8.1.99":  "Energie R+ Vorwerte ",
	"4.8.1":     "Energie R- ",
	"4.8.1.99":  "Energie R- Vorwerte ",
	"0.9.1":     "aktuelle Uhrzeit ",
	"0.9.2":     "aktuelles Datum ",
	"0.2.0":     "Firmware ID (Eichbereich) ",
	"C.60.5.1":  "Hardware ID (Revision.Produktzustand) ",
	"C.60.5.2":  "Firmware ID (Nicht-Eichbereich) ",
	"1.7.0":     "momentane Wirkleistung P+ ",
	"2.7.0":     "momentane Wirkleistung P- ",
	"3.7.0":     "momentane Blindleistung Q+ ",
	"4.7.0":     "momentane Blindleistung Q- ",
	"14.7":      "momentane Frequenz ",
	"32.7":      "momentaner Spannungswert L1 ",
	"52.7":      "momentaner Spannungswert L2 ",
	"72.7":      "momentaner Spannungswert L3 ",
	"31.7":      "momentaner Stromwert L1 ",
	"51.7":      "momentaner Stromwert L2 ",
	"71.7":      "momentaner Stromwert L3 ",
	"91.7":      "momentaner Stromwert N ",
	"81.7.4":    "momentaner Phasenwinkel U1-I1 ",
	"81.7.15":   "momentaner Phasenwinkel U2-I2 ",
	"81.7.26":   "momentaner Phasenwinkel U3-I3 ",
	"81.7.1":    "momentaner Phasenwinkel U1-U2 ",
	"81.7.12":   "momentaner Phasenwinkel U2-U3 ",
	"81.7.20":   "momentaner Phasenwinkel U3-U1 ",
	"32.36.0":   "Überspg. Zähler L1 + Datum +Zeit ",
	"52.36.0":   "Überspg. Zähler L2 + Datum +Zeit ",
	"72.36.0":   "Überspg. Zähler L3 + Datum +Zeit ",
	"32.32.0":   "Unterspg. Zähler L1 + Datum +Zeit ",
	"52.32.0":   "Unterspg. Zähler L2 + Datum +Zeit ",
	"72.32.0":   "Unterspg. Zähler L3 + Datum +Zeit ",
	"C.2.1":     "letzte Parametrierung. ",
	"C.7.0":     "Spg.ausfall L1-L3+ Datum +Zeit ",
	"C.7.1":     "Spg.ausfall L1+ Datum +Zeit ",
	"C.7.2":     "Spg.ausfall L2+ Datum +Zeit ",
	"C.7.3":     "Spg.ausfall L3+ Datum +Zeit ",
	"P.01":      "Lastprofilspeicher P+ ",
	"C.C.1":     "Strom-ADC-Overflow ",
	"C.1.8.1":   "Registriergrenze 1.8.1 ",
	"C.1.8.2":   "Registriergrenze 1.8.2 ",
	"C.1.8.3":   "Registriergrenze 1.8.3 ",
	"C.1.8.4":   "Registriergrenze 1.8.4 ",
	"C.1.8.5":   "Registriergrenze 1.8.5 ",
	"C.1.8.6":   "Registriergrenze 1.8.6 ",
	"C.2.8.1":   "Registriergrenze 2.8.1 ",
	"C.2.8.2":   "Registriergrenze 2.8.2 ",
	"C.2.8.3":   "Registriergrenze 2.8.3 ",
	"C.2.8.4":   "Registriergrenze 2.8.4 ",
	"C.2.8.5":   "Registriergrenze 2.8.5 ",
	"C.2.8.6":   "Registriergrenze 2.8.6 ",
	"C.60.4.1":  "BCAdr+Rev Zählertyp ",
	"C.60.4.2":  "BCAdr+Rev Netzbereitstellung ",
	"C.60.4.3":  "BCAdr+Rev Ableseeinheit-Zählertyp ",
	"C.60.4.4":  "BCAdr+Rev Lastprofiltyp ",
	"C.60.4.5":  "BCAdr+Rev Ableseeinheit-Lastprofiltyp ",
	"C.60.4.6":  "BCAdr+Rev PQ-Typ ",
	"C.60.4.7":  "BCAdr+Rev Ableseeinheit-PQ-Typ ",
	"C.60.4.8":  "BCAdr+Rev Erweiterungsmodul 0 ",
	"C.60.4.9":  "BCAdr+Rev Erweiterungsmodul 1 ",
	"C.60.4.10": "BCAdr+Rev Erweiterungsmodul 2 ",
	"C.60.4.11": "BCAdr+Rev Inkasso-Typ ",
	"C.70.0":    "Infofeld ",
	"C.70.2":    "Textfeld ",
	"C.71.1":    "Zustand der Abschalteinrichtung ",
	"C.71.2":    "Auslöseschwelle ",
	"C.71.3":    "Manipulationskontakt ",
	"C.71.4":    "Status der Zeitquelle ",
	"C.71.5":    "Betriebszustand ",
	"P.98.xx":   "Logbuch ",
	"L.70.2.1":  "PQ-Kumulierung Datum Uhrzeit ",
	"L.71.1":    "Überspannung 1 L1 ",
	"L.71.1.1":  "Überspannung 1 L1 Vorwert ",
	"L.71.2":    "Überspannung 2 L1 ",
	"L.71.2.1":  "Überspannung 2 L1 Vorwert ",
	"L.71.3":    "Überspannung 3 L1 ",
	"L.71.3.1":  "Überspannung 3 L1 Vorwert ",
	"L.71.4":    "Überspannung 1 L2 ",
	"L.71.4.1":  "Überspannung 1 L2 Vorwert ",
	"L.71.5":    "Überspannung 2 L2 ",
	"L.71.5.1":  "Überspannung 2 L2 Vorwert ",
	"L.71.6":    "Überspannung 3 L2 ",
	"L.71.6.1":  "Überspannung 3 L2 Vorwert ",
	"L.71.7":    "Überspannung 1 L3 ",
	"L.71.7.1":  "Überspannung 1 L3 Vorwert ",
	"L.71.8":    "Überspannung 2 L3 ",
	"L.71.8.1":  "Überspannung 2 L3 Vorwert ",
	"L.71.9":    "Überspannung 3 L3 ",
	"L.71.9.1":  "Überspannung 3 L3 Vorwert ",
	"L.72.1":    "DIPs 1 L1 ",
	"L.72.1.1":  "DIPs 1 L1 Vorwert ",
	"L.72.2":    "DIPs 2 L1 ",
	"L.72.2.1":  "DIPs 2 L1 Vorwert ",
	"L.72.3":    "DIPs 3 L1 ",
	"L.72.3.1":  "DIPs 3 L1 Vorwert ",
	"L.72.4":    "DIPs 1 L2 ",
	"L.72.4.1":  "DIPs 1 L2 Vorwert ",
	"L.72.5":    "DIPs 2 L2 ",
	"L.72.5.1":  "DIPs 2 L2 Vorwert ",
	"L.72.6":    "DIPs 3 L2 ",
	"L.72.6.1":  "DIPs 3 L2 Vorwert ",
	"L.72.7":    "DIPs 1 L3 ",
	"L.72.7.1":  "DIPs 1 L3 Vorwert ",
	"L.72.8":    "DIPs 2 L3 ",
	"L.72.8.1":  "DIPs 2 L3 Vorwert ",
	"L.72.9":    "DIPs 3 L3 ",
	"L.72.9.1":  "DIPs 3 L3 Vorwert ",
	"L.73.1":    "Spannung-Minimum1 ",
	"L.73.1.1":  "Spannung-Minimum1 Vorwert ",
	"L.73.2":    "Spannung-Minimum2 ",
	"L.73.2.1":  "Spannung-Minimum2 Vorwert ",
	"L.73.3":    "Spannung-Minimum3 ",
	"L.73.3.1":  "Spannung-Minimum3 Vorwert ",
	"L.73.4":    "Spannung-Minimum4 ",
	"L.73.4.1":  "Spannung-Minimum4 Vorwert ",
	"L.73.5":    "Spannung-Minimum5 ",
	"L.73.5.1":  "Spannung-Minimum5 Vorwert ",
	"L.73.6":    "Spannung-Minimum6 ",
	"L.73.6.1":  "Spannung-Minimum6 Vorwert ",
	"L.73.7":    "Spannung-Minimum7 ",
	"L.73.7.1":  "Spannung-Minimum7 Vorwert ",
	"L.73.8":    "Spannung-Minimum8 ",
	"L.73.8.1":  "Spannung-Minimum8 Vorwert ",
	"L.73.9":    "Spannung-Minimum9 ",
	"L.73.9.1":  "Spannung-Minimum9 Vorwert ",
	"L.73.10":   "Spannung-Minimum10 ",
	"L.73.10.1": "Spannung-Minimum10 Vorwert ",
	"L.73.11":   "Spannung-Minimum11 ",
	"L.73.11.1": "Spannung-Minimum11 Vorwert ",
	"L.74.1":    "Spannung-Maximum1 ",
	"L.74.1.1":  "Spannung-Maximum1 Vorwert ",
	"L.74.2":    "Spannung-Maximum2 ",
	"L.74.2.1":  "Spannung-Maximum2 Vorwert ",
	"L.74.3":    "Spannung-Maximum3 ",
	"L.74.3.1":  "Spannung-Maximum3 Vorwert ",
	"L.74.4":    "Spannung-Maximum4 ",
	"L.74.4.1":  "Spannung-Maximum4 Vorwert ",
	"L.74.5":    "Spannung-Maximum5 ",
	"L.74.5.1":  "Spannung-Maximum5 Vorwert ",
	"L.74.6":    "Spannung-Maximum6 ",
	"L.74.6.1":  "Spannung-Maximum6 Vorwert ",
	"L.74.7":    "Spannung-Maximum7 ",
	"L.74.7.1":  "Spannung-Maximum7 Vorwert ",
	"L.74.8":    "Spannung-Maximum8 ",
	"L.74.8.1":  "Spannung-Maximum8 Vorwert ",
	"L.74.9":    "Spannung-Maximum9 ",
	"L.74.9.1":  "Spannung-Maximum9 Vorwert ",
	"L.74.10":   "Spannung-Maximum10 ",
	"L.74.10.1": "Spannung-Maximum10 Vorwert ",
	"L.74.11":   "Spannung-Maximum11 ",
	"L.74.11.1": "Spannung-Maximum11 Vorwert ",
	"L.75.1":    "Spannung-Mittelwert1 ",
	"L.75.1.1":  "Spannung-Mittelwert1 Vorwert ",
	"L.75.2":    "Spannung-Mittelwert2 ",
	"L.75.2.1":  "Spannung-Mittelwert2 Vorwert ",
	"L.75.3":    "Spannung-Mittelwert3 ",
	"L.75.3.1":  "Spannung-Mittelwert3 Vorwert ",
	"L.75.4":    "Spannung-Mittelwert4 ",
	"L.75.4.1":  "Spannung-Mittelwert4 Vorwert ",
	"L.75.5":    "Spannung-Mittelwert5 ",
	"L.75.5.1":  "Spannung-Mittelwert5 Vorwert ",
	"L.75.6":    "Spannung-Mittelwert6 ",
	"L.75.6.1":  "Spannung-Mittelwert6 Vorwert ",
	"L.75.7":    "Spannung-Mittelwert7 ",
	"L.75.7.1":  "Spannung-Mittelwert7 Vorwert ",
	"L.75.8":    "Spannung-Mittelwert8 ",
	"L.75.8.1":  "Spannung-Mittelwert8 Vorwert ",
	"L.75.9":    "Spannung-Mittelwert9 ",
	"L.75.9.1":  "Spannung-Mittelwert9 Vorwert ",
	"L.75.10":   "Spannung-Mittelwert10 ",
	"L.75.10.1": "Spannung-Mittelwert10 Vorwert ",
	"L.75.11":   "Spannung-Mittelwert11 ",
	"L.75.11.1": "Spannung-Mittelwert11 Vorwert ",
}

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
	read(s)
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
		if strings.ContainsAny(line, "!") {
			break
		}
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
	val, ok := myMap[id]
	if ok {
		m["id_txt"] = val
	}
	vl := strings.Split(v[1], ")")
	value := vl[0]
	if strings.Index(value, "*") > 0 {
		r := strings.Split(value, "*")
		value = r[0]
		m["unit"] = r[1]
	}
	if value == "0.000" {
		return
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
		m["time_now"] = t.Format(time.RFC3339Nano)
	}

	jsonStr, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		log.Println(jsonStr)
		fmt.Println(string(jsonStr))
	}

}
