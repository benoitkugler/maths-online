import 'package:eleve/settings.dart';
import 'package:flutter/foundation.dart';
import 'package:logging/logging.dart';
import 'package:test/test.dart';

void main() {
  Logger.root.level = Level.ALL; // defaults to Level.INFO
  Logger.root.onRecord.listen((record) {
    print('${record.level.name}: ${record.time}: ${record.message}');
  });

  test('settings ...', () async {
    final settings = UserSettings(studentID: "msùlsùd", songs: [4, 56]);
    await saveUserSettings(settings);

    final settings2 = await loadUserSettings();
    expect(settings.studentID, equals(settings2.studentID));
    expect(listEquals(settings.songs, settings2.songs), equals(true));
  });
}
