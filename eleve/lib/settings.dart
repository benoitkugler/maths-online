import 'dart:convert';
import 'dart:io';

import 'package:logging/logging.dart';
import 'package:path_provider/path_provider.dart';

typedef UserSettings = Map<String, String>;

Future<File> _settingFile() async {
  final directory = await getApplicationDocumentsDirectory();
  return File('${directory.path}/.isyro_settings.json');
}

Future<UserSettings> loadUserSettings() async {
  final file = await _settingFile();
  try {
    final content = await file.readAsString();
    final dict = jsonDecode(content) as Map<String, dynamic>;
    return dict.map((key, value) => MapEntry(key, value as String));
  } catch (e) {
    Logger.root.info("loading settings: $e");
    return {};
  }
}

Future<void> saveUserSettings(UserSettings settings) async {
  final file = await _settingFile();
  await file.writeAsString(jsonEncode(settings));
  Logger.root.info("settings saved in ${file.path}");
  return;
}
