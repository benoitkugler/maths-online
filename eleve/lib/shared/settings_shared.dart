import 'dart:convert';

import 'package:eleve/shared/audio.dart';

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

  UserSettings(
      {this.studentPseudo = "",
      this.studentID = "",
      this.songs = Audio.DefaultPlaylist,
      Map<String, String>? trivialGameMetas,
      this.hasBeenLaunched = false})
      : trivialGameMetas = trivialGameMetas ?? {};

  String toJson() {
    return jsonEncode({
      studentPseudoKey: studentPseudo,
      studentIDKey: studentID,
      "songs": songs,
      "trivialGameMetas": trivialGameMetas,
      "hasBeenLaunched": hasBeenLaunched,
    });
  }

  factory UserSettings.fromJson(String source) {
    final dict = jsonDecode(source) as Map<String, dynamic>;
    var songs = [0, 1];
    if (dict["songs"] is List) {
      songs = (dict["songs"] as List<dynamic>).map((e) => e as int).toList();
    }
    final gameMetas = (dict["trivialGameMetas"] ?? <String, dynamic>{})
        as Map<String, dynamic>;
    return UserSettings(
      studentPseudo: (dict[studentPseudoKey] ?? "") as String,
      studentID: (dict[studentIDKey] ?? "") as String,
      songs: songs,
      trivialGameMetas:
          gameMetas.map((key, value) => MapEntry(key, value as String)),
      hasBeenLaunched: (dict["hasBeenLaunched"] ?? false) as bool,
    );
  }
}

abstract class SettingsStorage {
  Future<UserSettings> load();
  Future<void> save(UserSettings settings);
}
