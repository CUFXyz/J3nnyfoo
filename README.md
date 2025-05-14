# J3nnyfoo

J3nnyfoo is a small project on GoLang what getting json data and storing them into PostgreSQL
#
### Work in progress
* Auth system
* Metrics (hope so)

## Usage
### *Docker*
```bash
dockerbuild -t j3nnyfoo .
dockerrun -p 9090:9090 j3nnyfoo
```
#
### *W/o Docker*
```bash
make run
```
### or
```bash
go run ./cmd/main.go
```
