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
    await loadUserSettings();
  });

  test('settings ...', () async {
    final settings = {"key": "value", "id": "789"};
    await saveUserSettings(settings);

    final settings2 = await loadUserSettings();
    expect(mapEquals(settings, settings2), equals(true));
  });
}
