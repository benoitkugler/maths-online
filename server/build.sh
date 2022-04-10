echo "Pulling git..." &&
git pull &&
echo "Entering src/" &&
cd src && 
echo "Building (downloading deps if needed)..." && 
go build *.go &&
echo "Done." && 
# echo "Removing deps to free disk space..." &&
# go clean -cache -modcache && 
# echo "Done." &&
echo "Moving executable and leaving source..."
cd .. && 
mv src/main . &&
echo "Running dry..." &&
PORT=8000 IP=aaaa::2:a2cc ./main -dry &&
echo "All good !"