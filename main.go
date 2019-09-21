package main

import (
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type nautilusDataPoint struct {
	timestamp   time.Time
	speed       float64
	consumption float64
}

func main() {

	//read the data in
	csvFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	data, err := readCSV(csvFile)
	if err != nil {
		log.Fatalln(err)
	}

	//all of the http handlers
	http.HandleFunc("/total_distance", totalDistance(data))
	http.HandleFunc("/total_fuel", totalFuel(data))
	http.HandleFunc("/efficiency", efficiency(data))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func totalDistance(data []nautilusDataPoint) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func totalFuel(data []nautilusDataPoint) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		start := params.Get("start")
		end := params.Get("end")
		log.Println(start, end)
	}
}

func efficiency(data []nautilusDataPoint) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func readCSV(fileReader io.Reader) ([]nautilusDataPoint, error) {

	rows, err := csv.NewReader(fileReader).ReadAll()
	if err != nil {
		return nil, err
	}
	dataPoints := make([]nautilusDataPoint, len(rows)-1)
	//start reading at row 2 to skip processing the header
	for i, row := range rows[1:] {
		unix, err := strconv.ParseInt(row[0], 10, 64)
		if err != nil {
			log.Println(err)
		}
		speed, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			log.Println(err)
		}
		cons, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			log.Println(err)
		}
		var point nautilusDataPoint
		point.timestamp = time.Unix(unix, 0)
		point.speed = speed
		point.consumption = cons
		dataPoints[i] = point
	}
	cleanData(dataPoints)
	return dataPoints, nil
}

/*
  If a value is 0, set it to the value from the previous entry in the slice
  an exception to this rule  is if both speed and consumption are 0, that makes sense
*/
func cleanData(dataPoints []nautilusDataPoint) {
	for i, point := range dataPoints {
		if point.consumption == 0 && point.speed == 0 {
			continue
		}
		if point.consumption == 0 {
			point.consumption = dataPoints[i-1].consumption
		}
		if point.speed == 0 {
			point.speed = dataPoints[i-1].speed
		}
		dataPoints[i] = point
	}
}
