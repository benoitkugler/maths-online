import 'dart:convert';

import 'package:device_info_plus/device_info_plus.dart';
import 'package:eleve/shared/audio.dart';
import 'package:flutter/foundation.dart';

const studentPseudoKey = "client-pseudo";
const studentIDKey = "client-id";

/// [UserSettings] store the local parameters persisting
/// accross app launches.
class UserSettings {
  String studentPseudo;
  String studentID;
  PlaylistController songs;
  Map<String, String> trivialGameMetas;
  bool hasBeenLaunched;
  String ceinturesAnonymousID;

  UserSettings({
    this.studentPseudo = "",
    this.studentID = "",
    PlaylistController? songs,
    Map<String, String>? trivialGameMetas,
    this.hasBeenLaunched = false,
    this.ceinturesAnonymousID = "",
  }) : trivialGameMetas = trivialGameMetas ?? {},
       songs = songs ?? Audio.defaultPlaylist;

  String toJson() {
    return jsonEncode({
      studentPseudoKey: studentPseudo,
      studentIDKey: studentID,
      "songs": songs.toJson(),
      "trivialGameMetas": trivialGameMetas,
      "hasBeenLaunched": hasBeenLaunched,
      "ceinturesAnonymousID": ceinturesAnonymousID,
    });
  }

  factory UserSettings.fromJson(String source) {
    final dict = jsonDecode(source) as Map<String, dynamic>;

    PlaylistController? songs;
    if (dict.containsKey("songs")) {
      songs = PlaylistController.fromJson(dict["songs"]);
    }

    final gameMetas =
        (dict["trivialGameMetas"] ?? <String, dynamic>{})
            as Map<String, dynamic>;

    return UserSettings(
      studentPseudo: (dict[studentPseudoKey] ?? "") as String,
      studentID: (dict[studentIDKey] ?? "") as String,
      songs: songs ?? Audio.defaultPlaylist,
      trivialGameMetas: gameMetas.map(
        (key, value) => MapEntry(key, value as String),
      ),
      hasBeenLaunched: (dict["hasBeenLaunched"] ?? false) as bool,
      ceinturesAnonymousID: (dict["ceinturesAnonymousID"] ?? "") as String,
    );
  }
}

abstract class SettingsStorage {
  Future<UserSettings> load();
  Future<void> save(UserSettings settings);
}

/// SettingsHandler holds an active [UserSettings] and a [SettingsStorage]
class SettingsHandler {
  final SettingsStorage _storage;
  UserSettings _settings = UserSettings();

  SettingsHandler(this._storage);

  UserSettings get settings => _settings;

  Future<void> init() async {
    _settings = await _storage.load();
  }

  Future<void> saveStudentID(String studentID) async {
    _settings.studentID = studentID;
    await _storage.save(_settings);
  }

  Future<void> saveStudentPseudo(String studentPseudo) async {
    _settings.studentPseudo = studentPseudo;
    await _storage.save(_settings);
  }

  Future<void> saveSongs(PlaylistController songs) async {
    _settings.songs = songs;
    await _storage.save(_settings);
  }

  Future<void> saveHasBeenLaunched() async {
    _settings.hasBeenLaunched = true;
    await _storage.save(_settings);
  }

  Future<void> saveTrivialMeta(String gameCode, String gameMeta) async {
    _settings.trivialGameMetas[gameCode] = gameMeta;
    await _storage.save(_settings);
  }

  Future<void> removeTrivialMeta(String gameCode) async {
    _settings.trivialGameMetas.remove(gameCode);
    await _storage.save(_settings);
  }

  Future<void> saveCeinturesAnonymousID(String id) async {
    _settings.ceinturesAnonymousID = id;
    await _storage.save(_settings);
  }
}

Future<String> loadUserDeviceName() async {
  final deviceInfoPlugin = DeviceInfoPlugin();

  try {
    if (kIsWeb) {
      final deviceData = await deviceInfoPlugin.webBrowserInfo;
      return deviceData.browserName.name;
    } else {
      switch (defaultTargetPlatform) {
        case TargetPlatform.android:
          final deviceData = await deviceInfoPlugin.androidInfo;
          return "${deviceData.manufacturer} ${deviceData.brand}";
        case TargetPlatform.iOS:
          final deviceData = await deviceInfoPlugin.iosInfo;
          return deviceData.name;
        default:
          return "";
      }
    }
  } catch (e) {
    return "";
  }
}
