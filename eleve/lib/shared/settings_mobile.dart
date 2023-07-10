import 'dart:io';

import 'package:eleve/shared/settings_shared.dart';
import 'package:logging/logging.dart';
import 'package:path_provider/path_provider.dart';

class FileSettings implements SettingsStorage {
  Future<File> _settingFile() async {
    final directory = await getApplicationDocumentsDirectory();
    return File('${directory.path}/.isyro_settings.json');
  }

  @override
  Future<UserSettings> load() async {
    try {
      final file = await _settingFile();
      final content = await file.readAsString();
      return UserSettings.fromJson(content);
    } catch (e) {
      Logger.root.info("loading settings: $e");
      return UserSettings();
    }
  }

  @override
  Future<void> save(UserSettings settings) async {
    try {
      final file = await _settingFile();
      await file.writeAsString(settings.toJson());
      Logger.root.info("settings saved in ${file.path}");
    } catch (e) {
      Logger.root.info("saving settings: $e");
    }
    return;
  }
}
