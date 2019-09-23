# nautilusServer
take home interview for nautilus labs

## build and run

1. clone  or download this repository
2. cd into the root of the repository
3. make sure tests pass with `go test`
4. `go build`
5. `./nautilusServer <CSVFileName.csv>`

## sample requests
start and end query parameters are expected to be in unix epoch time as shown below
```
curl -X GET 'http://localhost:8080/total_distance?start=1561528800&end=1561604400'

curl -X GET 'http://localhost:8080/total_fuel?start=1561528800&end=1561604400'

curl -X GET 'http://localhost:8080/efficiency?start=1561528800&end=1561604400'
```
