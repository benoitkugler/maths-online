# run the main project in API dev mode, to test the app with a local server
echo "Forwarding localhost to device..." &&
~/Android/Sdk/platform-tools/adb reverse tcp:1323 tcp:1323 &&
echo "Running..." &&
flutter run -t lib/main_eleve_mobile.dart -v --dart-define=mode=dev #--release 