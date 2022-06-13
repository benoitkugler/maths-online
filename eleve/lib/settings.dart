import 'dart:convert';
import 'dart:io';

import 'package:eleve/build_mode.dart';
import 'package:eleve/classroom/join_classroom.dart';
import 'package:flutter/material.dart';
import 'package:logging/logging.dart';
import 'package:package_info_plus/package_info_plus.dart';
import 'package:path_provider/path_provider.dart';

const studentPseudoKey = "client-pseudo";
const studentIDKey = "client-id";

class UserSettings {
  String studentPseudo;
  String studentID;
  List<int> songs;

  UserSettings(
      {this.studentPseudo = "",
      this.studentID = "",
      this.songs = const [0, 1]});

  String toJson() {
    return jsonEncode({
      studentPseudoKey: studentPseudo,
      studentIDKey: studentID,
      "songs": songs
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
  var pseudoController = TextEditingController();

  @override
  void initState() {
    _loadUserSettings();
    _loadVersion();
    super.initState();
  }

  void _loadUserSettings() async {
    settings = await loadUserSettings();
    setState(() {
      pseudoController.text = settings.studentPseudo;
    });
  }

  void _loadVersion() async {
    PackageInfo packageInfo = await PackageInfo.fromPlatform();
    setState(() {
      version = packageInfo.version;
    });
  }

  void _saveUserSettings() async {
    settings.studentPseudo = pseudoController.text;
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

  void _showLeaveClassroom() async {
    final successLeaving = await confirmLeaveClassroom(
        widget.buildMode, settings.studentID, context);
    if (!successLeaving) {
      return;
    }
    setState(() {
      settings.studentID = "";
    });
    await saveUserSettings(settings);
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      backgroundColor: Theme.of(context).colorScheme.secondary,
      content: const Text("Classe quittée avec succès."),
    ));
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
                          OutlinedButton(
                              onPressed: _showJoinClassroom,
                              child: const Text("Rejoindre une classe")),
                          const SizedBox(height: 20),
                          Row(
                            children: [
                              const Padding(
                                padding: EdgeInsets.all(8.0),
                                child: Text(
                                  "Nom de joueur :",
                                  style: TextStyle(fontSize: 16),
                                ),
                              ),
                              SizedBox(
                                  width: 200,
                                  child: TextField(
                                    textAlign: TextAlign.center,
                                    controller: pseudoController,
                                    decoration: const InputDecoration(
                                        hintText: "Pseudo"),
                                  ))
                            ],
                          ),
                          const SizedBox(height: 30),
                          Center(
                              child: ElevatedButton(
                            child: const Text("Enregistrer"),
                            onPressed: _saveUserSettings,
                          )),
                        ]
                      : [
                          OutlinedButton(
                              onPressed: _showLeaveClassroom,
                              child: const Text("Quitter ma classe"))
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
