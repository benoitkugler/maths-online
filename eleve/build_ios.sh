# use the following command to collect shader warmup
# flutter run --profile --dart-define=mode=debug --cache-sksl --purge-persistent-cache
# then press M
# mv flutter_01.sksl.json flutter_ios.sksl.json
flutter build ipa --bundle-sksl-path flutter_ios.sksl.json