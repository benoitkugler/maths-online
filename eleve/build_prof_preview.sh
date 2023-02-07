flutter build web -t lib/main_prof_preview.dart --base-href=/prof-preview-app/ && 
# uncomment to build in "debug" mode
# flutter build web --profile --dart-define=Dart2jsOptimization=O0 -t lib/main_prof_preview.dart --base-href=/prof-preview-app/ && 
echo "Moving build to server/static/prof_preview..." && 
rm -r ../server/static/prof_preview/ && 
mkdir ../server/static/prof_preview && 
cp -r build/web/* ../server/static/prof_preview &&
echo "Removing unused music..." && 
rm ../server/static/prof_preview/assets/assets/music/* &&
echo "Fixing bug https://github.com/flutter/flutter/issues/53745..." && 
sed -i -e 's/return cache.addAll/cache.addAll/g' ../server/static/prof_preview/flutter_service_worker.js &&
echo "Done."