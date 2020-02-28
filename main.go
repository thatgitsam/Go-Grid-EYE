package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/knieriem/serport"
)

type gridEyeSensor struct {
	thermistor float32
	pixel      [8][8]float32
}

var sensor gridEyeSensor

func main() {

	// Open Serial Port
	port, name, err := serport.Choose("", "b115200")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("# connected to", name)

	// Create a new Scanner
	stream := bufio.NewScanner(port)

	go func() {
		// Read data and process
		for stream.Scan() {
			// Get this packet
			s := stream.Text()
			// deal with this packet being corrupt
			if i := strings.Index(s, "***"); i == -1 {
				continue
			}
			// Trim data if garbage ahead of *** start
			if i := strings.Index(s, "***"); i > 0 {
				s = s[i:]
			}
			// Trim prefix
			s = strings.TrimPrefix(s, "***")
			// Exit if length is wrong
			if len(s) != 130 {
				continue
			}

			// Get Thermister bytes
			hb := []byte(s[:2])
			if !(hb[1]&0b00001000 == 0) {
				hb[1] &= 0b00000111
				// Convert to Float and multiple by constant
				sensor.thermistor = float32(binary.LittleEndian.Uint16(hb)) * -0.0125
			} else {
				// Convert to Float and multiple by constant
				sensor.thermistor = float32(binary.LittleEndian.Uint16(hb)) * 0.0125
			}
			// Log Out
			// log.Println("Internal Temp:", thermistor)

			// Get Pixel bytes
			s = s[2:]
			// Process Pixels
			for y := 0; y < 8; y++ {
				for x := 0; x < 8; x++ {
					// Get relevant bytes
					t := []byte(s[y*x*2 : y*x*2+2])
					// If second byte is not zero, turn on bits 12-16 to convert to
					// twos compliment
					if !(t[1]&0b00001000 == 0) {
						t[1] |= 0b00000111
					}
					// Convert to Twos compliment
					val := binary.LittleEndian.Uint16(t)
					// Convert to Float and multiple by constant
					sensor.pixel[y][x] = float32(val) * 0.25
				}
			}
		}
	}()

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		px, err := json.Marshal(sensor.pixel)
		if err != nil {
			log.Fatal("Cannot encode to JSON ", err)
		}
		w.Write(px)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

	log.Println("Exiting")
}
