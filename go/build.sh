echo "Pulling git..." &&
git pull &&
echo "Building (downloading deps if needed)..." && 
go build main.go 
echo "Removing deps to free disk space..." &&
go clean -cache -modcache && 
echo "Running dry..." &&
PORT=8000 IP=aaaa::2:a2cc ./main -dry &&
echo "Done."