echo "Pulling git..." &&
git pull &&
echo "Downloading deps..." &&
go mod download && 
echo "Building..." && 
go build main.go 
echo "Removing deps to free disk space..." &&
go clean -cache -modcache && 
echo "Running dry..." &&
PORT=8000 IP=aaaa::2:a2cc ./main -dry &&
echo "Done."