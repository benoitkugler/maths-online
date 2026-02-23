import 'package:eleve/shared/settings_shared.dart';
import 'package:web/web.dart';

class LocalStorageSettings implements SettingsStorage {
  static const _settingsKey = "isyro_settings";

  @override
  Future<UserSettings> load() {
    final json = window.localStorage.getItem(_settingsKey) ?? "{}";
    return Future.sync(() => UserSettings.fromJson(json));
  }

  @override
  Future<void> save(UserSettings settings) {
    window.localStorage.setItem(_settingsKey, settings.toJson());
    return Future.sync(() {});
  }
}
