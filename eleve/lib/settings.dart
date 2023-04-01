import 'dart:convert';
import 'dart:io';

import 'package:eleve/audio.dart';
import 'package:eleve/build_mode.dart';
import 'package:eleve/classroom/join_classroom.dart';
import 'package:eleve/classroom/recap.dart';
import 'package:flutter/material.dart';
import 'package:logging/logging.dart';
import 'package:package_info_plus/package_info_plus.dart';
import 'package:path_provider/path_provider.dart';

const studentPseudoKey = "client-pseudo";
const studentIDKey = "client-id";

/// [UserSettings] store the local parameters persisting
/// accross app launches.
class UserSettings {
  String studentPseudo;
  String studentID;
  PlaylistController songs;
  bool hasBeenLaunched;

  UserSettings(
      {this.studentPseudo = "",
      this.studentID = "",
      this.songs = Audio.DefaultPlaylist,
      this.hasBeenLaunched = false});

  String toJson() {
    return jsonEncode({
      studentPseudoKey: studentPseudo,
      studentIDKey: studentID,
      "songs": songs,
      "hasBeenLaunched": hasBeenLaunched,
    });
  }

  factory UserSettings.fromJson(String source) {
    final dict = jsonDecode(source) as Map<String, dynamic>;
    var songs = [0, 1];
    if (dict["songs"] is List) {
      songs = (dict["songs"] as List<dynamic>).map((e) => e as int).toList();
    }
    return UserSettings(
      studentPseudo: dict[studentPseudoKey] as String,
      studentID: dict[studentIDKey] as String,
      songs: songs,
      hasBeenLaunched: dict["hasBeenLaunched"] as bool,
    );
  }
}

class Settings extends StatefulWidget {
  final BuildMode buildMode;
  const Settings(this.buildMode, {Key? key}) : super(key: key);

  @override
  State<Settings> createState() => _SettingsState();
}

class _SettingsState extends State<Settings> {
  UserSettings settings = UserSettings();
  String version = "";

  @override
  void initState() {
    _loadUserSettings();
    _loadVersion();
    super.initState();
  }

  void _loadUserSettings() async {
    settings = await loadUserSettings();
    setState(() {});
  }

  void _loadVersion() async {
    PackageInfo packageInfo = await PackageInfo.fromPlatform();
    setState(() {
      version = packageInfo.version;
    });
  }

  void _savePseudo(String pseudo) async {
    setState(() {
      settings.studentPseudo = pseudo;
    });
    await saveUserSettings(settings);
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      backgroundColor: Theme.of(context).colorScheme.secondary,
      content: const Text("Paramètres enregistrés"),
    ));
    Navigator.of(context).pop(settings);
  }

  void _showJoinClassroom() async {
    final idCrypted = await Navigator.of(context).push(
        MaterialPageRoute<String>(
            builder: (context) => JoinClassroomRoute(widget.buildMode)));
    if (idCrypted == null) {
      return;
    }

    setState(() {
      settings.studentID = idCrypted;
    });
    await saveUserSettings(settings);
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      backgroundColor: Theme.of(context).colorScheme.secondary,
      content: const Text("Classe rejointe avec succès."),
    ));
  }

  void _onInvalidStudentID() async {
    // clear the studentID
    setState(() {
      settings.studentID = "";
    });
    await saveUserSettings(settings);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text("Paramètres"),
      ),
      body: WillPopScope(
        onWillPop: () async {
          Navigator.of(context).pop(settings);
          return false;
        },
        child: Padding(
          padding: const EdgeInsets.all(10),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Expanded(
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  crossAxisAlignment: CrossAxisAlignment.center,
                  children: settings.studentID.isEmpty
                      ? [
                          _NotRegistred(settings.studentPseudo, _savePseudo,
                              _showJoinClassroom)
                        ]
                      : [
                          ClassroomCard(widget.buildMode, settings.studentID,
                              _onInvalidStudentID)
                        ],
                ),
              ),
              if (version.isNotEmpty)
                Align(
                  alignment: Alignment.bottomCenter,
                  child: Text(
                    "Version : $version",
                    style: const TextStyle(fontStyle: FontStyle.italic),
                  ),
                )
            ],
          ),
        ),
      ),
    );
  }
}

Future<File> _settingFile() async {
  final directory = await getApplicationDocumentsDirectory();
  return File('${directory.path}/.isyro_settings.json');
}

Future<UserSettings> loadUserSettings() async {
  try {
    final file = await _settingFile();
    final content = await file.readAsString();
    return UserSettings.fromJson(content);
  } catch (e) {
    Logger.root.info("loading settings: $e");
    return UserSettings();
  }
}

Future<void> saveUserSettings(UserSettings settings) async {
  try {
    final file = await _settingFile();
    await file.writeAsString(settings.toJson());
    Logger.root.info("settings saved in ${file.path}");
  } catch (e) {
    Logger.root.info("saving settings: $e");
  }
  return;
}

class _NotRegistred extends StatefulWidget {
  final String pseudo;
  final void Function(String) onSavePseudo;
  final void Function() onJoinClassroom;

  const _NotRegistred(this.pseudo, this.onSavePseudo, this.onJoinClassroom,
      {Key? key})
      : super(key: key);

  @override
  State<_NotRegistred> createState() => _NotRegistredState();
}

class _NotRegistredState extends State<_NotRegistred> {
  final TextEditingController _controller = TextEditingController();

  @override
  void initState() {
    _controller.text = widget.pseudo;
    super.initState();
  }

  @override
  void didUpdateWidget(covariant _NotRegistred oldWidget) {
    _controller.text = widget.pseudo;
    super.didUpdateWidget(oldWidget);
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(vertical: 6, horizontal: 8),
      decoration: BoxDecoration(
          border: Border.all(width: 2, color: Colors.lightBlue),
          borderRadius: const BorderRadius.all(Radius.circular(6))),
      child: Column(children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            const Text(
              "Nom de joueur :",
              style: TextStyle(fontSize: 16),
            ),
            SizedBox(
                width: 180,
                child: TextField(
                  textAlign: TextAlign.center,
                  controller: _controller,
                  decoration: const InputDecoration(hintText: "Pseudo"),
                ))
          ],
        ),
        const SizedBox(height: 10),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            OutlinedButton(
                onPressed: widget.onJoinClassroom,
                child: const Text("Rejoindre une classe")),
            Center(
                child: ElevatedButton(
              child: const Text("Enregistrer"),
              onPressed: () => widget.onSavePseudo(_controller.text),
            )),
          ],
        )
      ]),
    );
  }
}
