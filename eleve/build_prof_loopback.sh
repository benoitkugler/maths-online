flutter build web -t lib/main_prof_loopback.dart --base-href=/prof-loopback-app/ && 
# uncomment to build in "debug" mode
# flutter build web --profile --dart-define=Dart2jsOptimization=O0 -t lib/main_prof_loopback.dart --base-href=/prof-loopback-app/ && 
echo "Moving build to server/static/prof_loopback..." && 
rm -r ../server/static/prof_loopback/ && 
mkdir ../server/static/prof_loopback && 
cp -r build/web/* ../server/static/prof_loopback &&
echo "Removing unused music..." && 
rm ../server/static/prof_loopback/assets/lib/music/* &&
echo "Fixing bug https://github.com/flutter/flutter/issues/53745..." && 
sed -i -e 's/return cache.addAll/cache.addAll/g' ../server/static/prof_loopback/flutter_service_worker.js &&
echo "Done."