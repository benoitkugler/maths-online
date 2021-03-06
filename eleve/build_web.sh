flutter build web -t lib/main_eleve_web.dart --base-href=/test-eleve/  && 
echo "Moving build to server/static/eleve..." && 
rm -r ../server/static/eleve/ && 
mkdir ../server/static/eleve && 
cp -r build/web/* ../server/static/eleve &&
echo "Fixing bug https://github.com/flutter/flutter/issues/53745..." && 
sed -i -e 's/return cache.addAll/cache.addAll/g' ../server/static/eleve/flutter_service_worker.js &&
echo "Done."