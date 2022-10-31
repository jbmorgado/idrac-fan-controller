package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

var (
	host     = flag.String("H", "", "host")
	user     = flag.String("U", "", "user")
	password = flag.String("P", "", "password")
)

var ipmiArgs = []string{"-I", "lanplus", "-H", "192.168.1.1", "-U", "root", "-P", "password"}
var ipmiReset = []string{"raw", "0x30", "0x30", "0x01", "0x01"}
var ipmiManual = []string{"raw", "0x30", "0x30", "0x01", "0x00"}
var ipmiSpeed = []string{"raw", "0x30", "0x30", "0x02", "0xff"}

type Temperature struct {
	Probe    string `json:"probe"`
	Code     string `json:"code"`
	Status   string `json:"status"`
	Position string `json:"position"`
	Value    int    `json:"value"`
	Unit     string `json:"unit"`
}

func main() {
	flag.Parse()

	out, err := exec.Command("ipmitool", "-I", "lanplus", "-H", *host, "-U", *user, "-P", *password, "sdr", "type", "temperature").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out)
	r := csv.NewReader(strings.NewReader(string(out)))
	// testString := "Inlet Temp       | 04h | ok  |  7.1 | 25 degrees C\nExhaust Temp     | 01h | ok  |  7.1 | 43 degrees C\nTemp             | 0Eh | ok  |  3.1 | 71 degrees C\nTemp             | 0Fh | ok  |  3.2 | 71 degrees C\n"

	// fmt.Printf("%s", testString)
	// r := csv.NewReader(strings.NewReader(testString))

	r.Comma = '|'

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var temperatures []Temperature

	for record := range records {
		var tReading = strings.TrimSpace(records[record][4])
		var tNumber, _ = strconv.Atoi(strings.Split(tReading, " ")[0])
		var tUnit = strings.Split(tReading, " ")[2]

		temperatures = append(temperatures, Temperature{
			Probe:    strings.TrimSpace(records[record][0]),
			Code:     strings.TrimSpace(records[record][1]),
			Status:   strings.TrimSpace(records[record][2]),
			Position: strings.TrimSpace(records[record][3]),
			Value:    tNumber,
			Unit:     tUnit,
		})
	}
	temperatureJson, _ := json.Marshal(temperatures)
	fmt.Println(string(temperatureJson))

}
