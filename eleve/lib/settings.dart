import 'package:eleve/build_mode.dart';
import 'package:eleve/classroom/join_classroom.dart';
import 'package:eleve/classroom/recap.dart';
import 'package:eleve/shared/settings_shared.dart';
import 'package:flutter/material.dart';
import 'package:package_info_plus/package_info_plus.dart';

class Settings extends StatefulWidget {
  final BuildMode buildMode;
  final SettingsStorage handler;

  const Settings(this.buildMode, this.handler, {Key? key}) : super(key: key);

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
    final newSettings = await widget.handler.load();
    setState(() {
      settings = newSettings;
    });
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
    await widget.handler.save(settings);
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
    await widget.handler.save(settings);
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
    await widget.handler.save(settings);
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
