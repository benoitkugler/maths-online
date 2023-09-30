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
      PlaylistController? songs,
      Map<String, String>? trivialGameMetas,
      this.hasBeenLaunched = false})
      : trivialGameMetas = trivialGameMetas ?? {},
        songs = songs ?? PlaylistController([], false);

  String toJson() {
    return jsonEncode({
      studentPseudoKey: studentPseudo,
      studentIDKey: studentID,
      "songs": songs.toJson(),
      "trivialGameMetas": trivialGameMetas,
      "hasBeenLaunched": hasBeenLaunched,
    });
  }

  factory UserSettings.fromJson(String source) {
    final dict = jsonDecode(source) as Map<String, dynamic>;

    PlaylistController? songs;
    if (dict.containsKey("songs")) {
      songs = PlaylistController.fromJson(dict["songs"]);
    }

    final gameMetas = (dict["trivialGameMetas"] ?? <String, dynamic>{})
        as Map<String, dynamic>;

    return UserSettings(
      studentPseudo: (dict[studentPseudoKey] ?? "") as String,
      studentID: (dict[studentIDKey] ?? "") as String,
      songs: songs ?? Audio.defaultPlaylist,
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
