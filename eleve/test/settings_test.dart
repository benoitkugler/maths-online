import 'package:eleve/shared/audio.dart';
import 'package:eleve/shared/settings_mobile.dart';
import 'package:eleve/shared/settings_shared.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:logging/logging.dart';

void main() {
  Logger.root.level = Level.ALL; // defaults to Level.INFO
  Logger.root.onRecord.listen((record) {
    print('${record.level.name}: ${record.time}: ${record.message}');
  });

  test('settings ...', () async {
    TestWidgetsFlutterBinding.ensureInitialized();
    const MethodChannel channel =
        MethodChannel('plugins.flutter.io/path_provider');
    TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger
        .setMockMethodCallHandler(channel, (MethodCall methodCall) async {
      return "/tmp";
    });

    final handler = FileSettings();

    final settings = UserSettings(
        studentID: "msùlsùd",
        songs: PlaylistController([4, 56], true),
        trivialGameMetas: {"1": "2"});
    await handler.save(settings);

    final settings2 = await handler.load();
    expect(settings.studentID, equals(settings2.studentID));
    expect(
        listEquals(settings.songs.songs, settings2.songs.songs), equals(true));
    expect(settings.songs.random, settings2.songs.random);
    expect(mapEquals(settings.trivialGameMetas, settings2.trivialGameMetas),
        equals(true));

    final oldVersion = UserSettings.fromJson("""{"songs" :[0, 1]}""");
    expect(listEquals(oldVersion.songs.songs, [0, 1]), true);
  });
}
