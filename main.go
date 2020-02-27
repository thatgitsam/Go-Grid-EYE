package main

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/tarm/serial"
)

func main() {

	// Configure and connect to Serial Port
	c := &serial.Config{Name: "/dev/ttyACM0", Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	// n, err := s.Write([]byte("test"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Create a buffer for our use
	var b bytes.Buffer

	// Loop reading data as it appears
	for {
		part := make([]byte, 135)
		_, err := s.Read(part)
		b.Write(part)

		if err != nil {
			log.Fatal(err)
		}

		x := bytes.Index(b.Bytes(), []byte("\r\n"))

		if x >= 0 {
			l := bytes.TrimPrefix(b.Bytes(), []byte("\r\n"))
			b.Reset()

			if bytes.Compare(l[:3], []byte("***")) == 0 {
				log.Printf("Head Matched: %x", l[0:5])

				// Process Thermister Value (internal Heat Temp Sensor?)
				// Appears to be 12bit signed
				// 0000 0000 000x XXXX
				// Where 0000 0000 000: Value
				// Where x: 0=Positive, 1=Negative
				// Where XXXX: Ignore

				// Copying from a Python codeset supplied with the Grid-EYE
				// Value Multiplier 0.0125
				var t float64
				tb := l[3:5]
				if tb[1]&8 == 8 {
					tb[1] = tb[1] + 7
					t = float64(binary.LittleEndian.Uint16(tb)) * -0.0125
				} else {
					t = float64(binary.LittleEndian.Uint16(tb)) * 0.0125
				}
				log.Printf("Temp: %f", t)

				// Process Pixels

			}
		}

		// Process if we have full data packet
		// if b. {

		// 	// Output number of bytes read
		// 	log.Printf("Data Size: %d", n)

		// 	// Entire Read Data
		// 	log.Printf("RawData: %q", buffer[:n])

		// 	// Command Head
		// 	head := buffer[:3]

		// 	log.Printf("Data Head: %q", head)
		// 	// Clear out for next round
		// 	buffer = buffer[:0]
		// 	log.Println()
		// }
	}
}
