flutter build web --base-href=/test-eleve/ && 
echo "Moving build to server/static/eleve..." && 
rm -r ../server/static/eleve/ && 
mkdir ../server/static/eleve && 
cp -r build/web/* ../server/static/eleve &&
echo "Done."