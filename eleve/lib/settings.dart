import 'dart:convert';
import 'dart:io';

import 'package:flutter/material.dart';
import 'package:logging/logging.dart';
import 'package:path_provider/path_provider.dart';

typedef UserSettings = Map<String, String>;

class Settings extends StatefulWidget {
  const Settings({Key? key}) : super(key: key);

  @override
  State<Settings> createState() => _SettingsState();
}

const studentPseudoKey = "client-pseudo";
const studentIDKey = "client-id";

class _SettingsState extends State<Settings> {
  UserSettings settings = {};
  var pseudoController = TextEditingController();

  @override
  void initState() {
    _loadUserSettings();
    super.initState();
  }

  void _loadUserSettings() async {
    settings = await loadUserSettings();
    setState(() {
      pseudoController.text = settings[studentPseudoKey] ?? "";
    });
  }

  void _saveUserSettings() async {
    settings[studentPseudoKey] = pseudoController.text;
    await saveUserSettings(settings);
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      backgroundColor: Theme.of(context).colorScheme.secondary,
      content: const Text("Paramètres enregistrés"),
    ));
    Navigator.of(context).pop();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text("Paramètres"),
      ),
      body: Padding(
        padding: const EdgeInsets.all(10),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
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
                    ))
              ],
            ),
            const SizedBox(height: 30),
            Center(
                child: ElevatedButton(
              child: const Text("Enregistrer"),
              onPressed: _saveUserSettings,
            ))
          ],
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
