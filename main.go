package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"strings"

	"github.com/knieriem/serport"
)

func main() {

	// Open Serial Port
	port, name, err := serport.Choose("", "b115200")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("# connected to", name)

	// Create a new Scanner
	stream := bufio.NewScanner(port)

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
		var thermistor float32
		hb := []byte(s[:2])
		if !(hb[1]&0b00001000 == 0) {
			hb[1] &= 0b00000111
			// Convert to Float and multiple by constant
			thermistor = float32(binary.LittleEndian.Uint16(hb)) * -0.0125
		} else {
			// Convert to Float and multiple by constant
			thermistor = float32(binary.LittleEndian.Uint16(hb)) * 0.0125
		}
		// Log Out
		log.Println("Internal Temp:", thermistor)

		// Get Pixel bytes
		s = s[2:]
		// Process Pixels
		px := [64]float32{}
		for i := 0; i < 64; i++ {
			// Get relevant bytes
			t := []byte(s[i*2 : i*2+2])
			// If second byte is not zero, turn on bits 12-16 to convert to
			// twos compliment
			//fmt.Printf("%08b\n", t)
			if !(t[1]&0b00001000 == 0) {
				t[1] |= 0b00000111
			}
			// if t[1]&8 == 0 {
			// 	t[1] = t[1] + 248
			// }
			// Convert to Twos compliment
			x := binary.LittleEndian.Uint16(t)
			// Convert to Float and multiple by constant
			px[i] = float32(x) * 0.25

		}
		// Output values
		log.Println(px)
	}

	log.Println("Exiting")
}
