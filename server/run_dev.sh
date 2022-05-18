# equivalent for go run src/main.go
cd src && 
go build *.go && 
mv main .. && 
cd .. && 
DEMO_PIN_TRIVIAL=1234 ./main -dev
