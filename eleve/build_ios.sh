# use the following command to collect shader warmup
# flutter run -t lib/main_eleve_mobile.dart --profile --dart-define=mode=debug --cache-sksl --purge-persistent-cache
# then press M
# mv flutter_01.sksl.json flutter_ios.sksl.json
# flutter build ipa -t lib/main_eleve_mobile.dart --bundle-sksl-path flutter_ios.sksl.json
flutter build ipa -t lib/main_eleve_mobile.dart
flutter build ios -t lib/main_eleve_mobile.dart # to locally preview