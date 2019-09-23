package main

import (
	"io"
	"io/ioutil"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_totalDistance(t *testing.T) {
	type args struct {
		data []nautilusDataPoint
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
		{
			"test start & end match entry",
			args{[]nautilusDataPoint{
				{time.Unix(1561582800, 0), 17.8875234, 7.189651947},
				{time.Unix(1561586400, 0), 18.19978543, 7.410584136}, //1 hour
				{time.Unix(1561593600, 0), 18.13383457, 7.373623021}, // 2 hour
				{time.Unix(1561597200, 0), 18.202547, 7.336315167},   // 1 hour
			},
			},
			[]byte("{\"total_distance\":72.42092883}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := totalDistance(tt.args.data)
			r := httptest.NewRequest("GET", "http://nautilus.com/total_distance?start=1561582800&end=1561597200", nil)
			w := httptest.NewRecorder()
			handler(w, r)
			got, _ := ioutil.ReadAll(w.Result().Body)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("totalFuel() = %v, want %v", string(got), string(tt.want))
			}

		})
	}
}

func Test_totalFuel(t *testing.T) {
	type args struct {
		data []nautilusDataPoint
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
		{
			"test start & end match entry",
			args{[]nautilusDataPoint{
				{time.Unix(1561582800, 0), 17.8875234, 7.189651947},
				{time.Unix(1561586400, 0), 18.19978543, 7.410584136}, //60 minutes
				{time.Unix(1561593600, 0), 18.13383457, 7.373623021}, // 120  minutes
				{time.Unix(1561597200, 0), 18.202547, 7.336315167},   // 60 minutes
			},
			},
			[]byte("{\"total_fuel\":1763.0665944}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := totalFuel(tt.args.data)
			r := httptest.NewRequest("GET", "http://nautilus.com/total_distance?start=1561582800&end=1561597200", nil)
			w := httptest.NewRecorder()
			handler(w, r)
			got, _ := ioutil.ReadAll(w.Result().Body)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("totalFuel() = %v, want %v", string(got), string(tt.want))
			}

		})
	}
}

func Test_efficiency(t *testing.T) {
	type args struct {
		data []nautilusDataPoint
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
		{
			"test start & end match entry",
			args{[]nautilusDataPoint{
				{time.Unix(1561582800, 0), 17.8875234, 7.189651947},
				{time.Unix(1561586400, 0), 18.19978543, 7.410584136},
				{time.Unix(1561593600, 0), 18.13383457, 7.373623021},
				{time.Unix(1561597200, 0), 18.202547, 7.336315167},
			},
			},
			[]byte("{\"efficiency\":0.0006854774074310726}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := efficiency(tt.args.data)
			r := httptest.NewRequest("GET", "http://nautilus.com/total_distance?start=1561582800&end=1561597200", nil)
			w := httptest.NewRecorder()
			handler(w, r)
			got, _ := ioutil.ReadAll(w.Result().Body)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("totalFuel() = %v, want %v", string(got), string(tt.want))
			}

		})
	}
}

func Test_readCSV(t *testing.T) {
	type args struct {
		fileReader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []nautilusDataPoint
		wantErr bool
	}{
		{
			"small good input",
			args{strings.NewReader("timestamp,speed,fuel\n" +
				"1561522200,14.68950611,4.813948324\n" +
				"1561525200,16.63974605,5.931167328\n" +
				"1561528800,16.56744713,5.924517165")},
			[]nautilusDataPoint{
				{time.Unix(1561522200, 0), 14.68950611, 4.813948324},
				{time.Unix(1561525200, 0), 16.63974605, 5.931167328},
				{time.Unix(1561528800, 0), 16.56744713, 5.924517165},
			},
			false,
		},
		{
			"missing one speed reading",
			args{strings.NewReader("timestamp,speed,fuel\n" +
				"1561543200,17.75503797,7.846647337\n" +
				"1561546800,,7.914451626\n" +
				"1561550400,17.81794197,7.75468518")},
			[]nautilusDataPoint{
				{time.Unix(1561543200, 0), 17.75503797, 7.846647337},
				{time.Unix(1561546800, 0), 17.75503797, 7.914451626},
				{time.Unix(1561550400, 0), 17.81794197, 7.75468518},
			},
			false,
		},
		{
			"missing one consumption reading",
			args{strings.NewReader("timestamp,speed,fuel\n" +
				"1561522200,14.68950611,4.813948324\n" +
				"1561525200,16.63974605,\n" +
				"1561528800,16.56744713,5.924517165")},
			[]nautilusDataPoint{
				{time.Unix(1561522200, 0), 14.68950611, 4.813948324},
				{time.Unix(1561525200, 0), 16.63974605, 4.813948324},
				{time.Unix(1561528800, 0), 16.56744713, 5.924517165},
			},
			false,
		},
		{
			"small good input with a zero entry",
			args{strings.NewReader("timestamp,speed,fuel\n" +
				"1561522200,14.68950611,4.813948324\n" +
				"1561525200,0,0\n" +
				"1561528800,16.56744713,5.924517165")},
			[]nautilusDataPoint{
				{time.Unix(1561522200, 0), 14.68950611, 4.813948324},
				{time.Unix(1561525200, 0), 0, 0},
				{time.Unix(1561528800, 0), 16.56744713, 5.924517165},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readCSV(tt.args.fileReader)
			if (err != nil) != tt.wantErr {
				t.Errorf("readCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}
