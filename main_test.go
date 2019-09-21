package main

import (
	"io"
	"reflect"
	"strings"
	"testing"
	"time"
)

// func Test_totalDistance(t *testing.T) {
// 	type args struct {
// 		data []nautilusDataPoint
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want func(http.ResponseWriter, *http.Request)
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := totalDistance(tt.args.data); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("totalDistance() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_totalFuel(t *testing.T) {
// 	type args struct {
// 		data []nautilusDataPoint
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want func(http.ResponseWriter, *http.Request)
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := totalFuel(tt.args.data); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("totalFuel() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_efficiency(t *testing.T) {
// 	type args struct {
// 		data []nautilusDataPoint
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want func(http.ResponseWriter, *http.Request)
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := efficiency(tt.args.data); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("efficiency() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

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
			args{strings.NewReader(`timestamp,speed,fuel
1561522200,14.68950611,4.813948324
1561525200,16.63974605,5.931167328
1561528800,16.56744713,5.924517165`)},
			[]nautilusDataPoint{
				{time.Unix(1561522200, 0), 14.68950611, 4.813948324},
				{time.Unix(1561525200, 0), 16.63974605, 5.931167328},
				{time.Unix(1561528800, 0), 16.56744713, 5.924517165},
			},
			false,
		},
		{
			"missing one speed reading",
			args{strings.NewReader(`timestamp,speed,fuel
1561543200,17.75503797,7.846647337
1561546800,,7.914451626
1561550400,17.81794197,7.75468518`)},
			[]nautilusDataPoint{
				{time.Unix(1561543200, 0), 17.75503797, 7.846647337},
				{time.Unix(1561546800, 0), 17.75503797, 7.914451626},
				{time.Unix(1561550400, 0), 17.81794197, 7.75468518},
			},
			false,
		},
		{
			"missing one consumption reading",
			args{strings.NewReader(`timestamp,speed,fuel
1561522200,14.68950611,4.813948324
1561525200,16.63974605,
1561528800,16.56744713,5.924517165`)},
			[]nautilusDataPoint{
				{time.Unix(1561522200, 0), 14.68950611, 4.813948324},
				{time.Unix(1561525200, 0), 16.63974605, 4.813948324},
				{time.Unix(1561528800, 0), 16.56744713, 5.924517165},
			},
			false,
		},
		{
			"small good input with a zero entry",
			args{strings.NewReader(`timestamp,speed,fuel
1561522200,14.68950611,4.813948324
1561525200,0,0
1561528800,16.56744713,5.924517165`)},
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
