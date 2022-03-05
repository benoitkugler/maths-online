echo "Pulling git..." &&
git pull &&
echo "Building (downloading deps if needed)..." && 
go build main.go &&
echo "Done." && 
echo "Removing deps to free disk space..." &&
go clean -cache -modcache && 
echo "Done." &&
echo "Running dry..." &&
PORT=8000 IP=aaaa::2:a2cc ./main -dry &&
echo "All good !"