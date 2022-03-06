# equivalent for go run src/main.go
cd src && 
go build main.go && 
mv main .. && 
cd .. && 
./main -dev
