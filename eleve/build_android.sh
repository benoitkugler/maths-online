# use the following command to collect shader warmup
# flutter run --profile --cache-sksl --purge-persistent-cache
# then press M
# mv flutter_01.sksl.json flutter_android.sksl.json
flutter build appbundle --bundle-sksl-path flutter_01.sksl.json