echo "Pulling git..." &&
git pull &&
echo "Downloading deps..." &&
go mod download && 
echo "Building..." && 
go build main.go 
echo "Removing deps to free disk space..." &&
go clean -cache -modcache && 
echo "Running dry..." &&
./main -dry  &&
echo "Done."