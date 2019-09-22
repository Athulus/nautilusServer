# nautilusServer
take home interview for nautilus labs

## build and run

1. `go build`
2. `./nautilusServer <CSVFileName.csv>`

## sample requests
```
curl -X GET 'http://localhost:8080/total_distance?start=1561528800&end=1561604400'

curl -X GET 'http://localhost:8080/total_fuel?start=1561528800&end=1561604400'

curl -X GET 'http://localhost:8080/efficiency?start=1561528800&end=1561604400'
```
