import 'package:eleve/shared/settings_shared.dart';
import 'dart:html';

class LocalStorageSettings implements SettingsStorage {
  static const _settingsKey = "isyro_settings";

  @override
  Future<UserSettings> load() {
    final json = window.localStorage[_settingsKey] ?? "{}";
    return Future.sync(() => UserSettings.fromJson(json));
  }

  @override
  Future<void> save(UserSettings settings) {
    window.localStorage[_settingsKey] = settings.toJson();
    return Future.sync(() {});
  }
}
