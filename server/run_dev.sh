# equivalent for go run src/main.go
cd src && 
go build *.go && 
mv main .. && 
cd .. && 
./main -dev
