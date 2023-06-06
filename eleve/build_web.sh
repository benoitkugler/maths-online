flutter build web -t lib/main_eleve_web.dart --base-href=/test-eleve/  && 
echo "Copying build to server/static/eleve..." && 
rm -r ../server/static/eleve/ && 
mkdir ../server/static/eleve && 
cp -r build/web/* ../server/static/eleve &&
echo "Done."