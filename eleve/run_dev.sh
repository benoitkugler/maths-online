# run the main project in API dev mode, to test the app with a local server
~/Android/Sdk/platform-tools/adb reverse tcp:1323 tcp:1323
flutter run -t lib/main_eleve_mobile.dart --release --dart-define=mode=dev