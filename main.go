package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
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
	type response struct {
		TotalDistance float64 `json:"total_distance"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		start, end, err := getStartAndEndTime(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		dataSlice := sliceDataByTime(data, start, end)

		//assuming here that speed is constant from one sensor reading to the next
		distance := 0.0
		for i := 1; i < len(dataSlice); i++ {
			timeDifference := dataSlice[i].timestamp.Sub(dataSlice[i-1].timestamp).Hours()
			distance += timeDifference * dataSlice[i].speed
		}

		distanceResposne := response{distance}
		body, err := json.Marshal(distanceResposne)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(body)
	}
}

func totalFuel(data []nautilusDataPoint) func(http.ResponseWriter, *http.Request) {
	type response struct {
		TotalFuel float64 `json:"total_fuel"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		start, end, err := getStartAndEndTime(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		dataSlice := sliceDataByTime(data, start, end)

		//assuming here that consumption is constant from one sensor reading to the next
		fuel := 0.0
		for i := 1; i < len(dataSlice); i++ {
			timeDifference := dataSlice[i].timestamp.Sub(dataSlice[i-1].timestamp).Minutes()
			fuel += timeDifference * dataSlice[i].consumption
		}

		fuelResposne := response{fuel}
		body, err := json.Marshal(fuelResposne)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(body)
	}
}

func efficiency(data []nautilusDataPoint) func(http.ResponseWriter, *http.Request) {
	type response struct {
		Efficiency float64 `json:"efficiency"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		start, end, err := getStartAndEndTime(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		dataSlice := sliceDataByTime(data, start, end)

		//assuming here that consumption adn distance are constant from one sensor reading to the next
		fuel := 0.0
		distance := 0.0
		var mpg []float64
		for i := 1; i < len(dataSlice); i++ {
			timeDifference := dataSlice[i].timestamp.Sub(dataSlice[i-1].timestamp)
			fuel = timeDifference.Minutes() * dataSlice[i].consumption
			distance = timeDifference.Hours() * dataSlice[i].speed
			// mpg = miles/hour * hour/minute * minute/gallon
			mpg = append(mpg, distance/60/fuel)
		}

		efficiencyResposne := response{average(mpg)}
		body, err := json.Marshal(efficiencyResposne)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(body)
	}
}

func getStartAndEndTime(r *http.Request) ( /*start time*/ time.Time /*end time*/, time.Time, error) {
	var startTime, endTime time.Time
	params := r.URL.Query()
	start := params.Get("start")
	end := params.Get("end")
	if start == "" || end == "" {
		//we need both of these parameters so return an error
		//and we should write an http error in the handler
		return startTime, endTime, errors.New("start and end query parameters must be present")
	}
	unix, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		return startTime, endTime, err
	}
	startTime = time.Unix(unix, 0)
	unix, err = strconv.ParseInt(end, 10, 64)
	if err != nil {
		return startTime, endTime, err
	}
	endTime = time.Unix(unix, 0)

	log.Println(start, end)
	return startTime, endTime, nil
}

//this will slice the data set to the appropriate length to aggregate
func sliceDataByTime(data []nautilusDataPoint, start time.Time, end time.Time) []nautilusDataPoint {
	var sub []nautilusDataPoint
	setStart := false
	for i, point := range data {

		// do this top branch until we find the start time
		if !setStart {
			if point.timestamp.Equal(start) {
				setStart = true
				sub = append(sub, point)
			} else if point.timestamp.After(start) {
				setStart = true
				prevPoint := data[i-1]
				prevPoint.timestamp = start
				sub = append(sub, prevPoint)
				sub = append(sub, point)

			}
		} else { //then find end time
			if point.timestamp.Equal(end) {
				sub = append(sub, point)
				break
			} else if point.timestamp.After(end) {
				point.timestamp = end
				sub = append(sub, point)
				break
			}
		}

	}
	return sub
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
			//it is assumed that we will never have missing timestamps
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

// If a value is 0, set it to the value from the previous entry in the slice
// an exception to this rule  is if both speed and consumption are 0, that makes sense
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

//I could be pulling in  a third party lib like gonum
//but if I am only using an average it's probably not worth it
func average(numbers []float64) float64 {
	sum := 0.0
	for _, num := range numbers {
		sum += num
	}
	return sum / float64(len(numbers))
}
