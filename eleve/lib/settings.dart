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
  String deviceName = "";

  @override
  void initState() {
    _loadUserSettings();
    _loadVersion();
    _loadDeviceName();
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

  void _loadDeviceName() async {
    final name = await loadUserDeviceName();
    setState(() {
      deviceName = name;
    });
  }

  void _savePseudo(String pseudo) async {
    setState(() {
      settings.studentPseudo = pseudo;
    });
    await widget.handler.save(settings);
    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      backgroundColor: Theme.of(context).colorScheme.secondary,
      content: const Text("Paramètres enregistrés"),
    ));
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
        title: const Text("Profil"),
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
                    children: [
                      Container(
                          padding: const EdgeInsets.all(6),
                          decoration: BoxDecoration(
                            border:
                                Border.all(color: Colors.lightBlue, width: 2),
                            borderRadius:
                                const BorderRadius.all(Radius.circular(6)),
                          ),
                          child: settings.studentID.isEmpty
                              ? _NotRegistred(settings.studentPseudo,
                                  _savePseudo, _showJoinClassroom)
                              : ClassroomCard(widget.buildMode,
                                  settings.studentID, _onInvalidStudentID)),
                    ]),
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  if (version.isNotEmpty)
                    Text(
                      "Version : $version",
                      style: const TextStyle(fontStyle: FontStyle.italic),
                    ),
                  if (deviceName.isNotEmpty)
                    Text(
                      deviceName,
                      style: const TextStyle(fontStyle: FontStyle.italic),
                    ),
                ],
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
    return Padding(
      padding: const EdgeInsets.all(8.0),
      child: Column(
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.start,
            children: [
              const Icon(Icons.no_accounts),
              const SizedBox(width: 12),
              Text(
                "Invité",
                style: Theme.of(context).textTheme.titleLarge,
              ),
            ],
          ),
          const SizedBox(height: 24),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Text("Pseudo : "),
              Text(widget.pseudo),
              IconButton(
                  splashRadius: 24,
                  onPressed: showEditPseudo,
                  icon: const Icon(Icons.edit))
            ],
          ),
          const SizedBox(height: 24),
          const Divider(thickness: 4),
          const SizedBox(height: 24),
          ElevatedButton.icon(
              onPressed: widget.onJoinClassroom,
              icon: const Icon(Icons.manage_accounts),
              label: const Text("Rejoindre une classe"))
        ],
      ),
    );
  }

  void showEditPseudo() async {
    final valid = await showDialog<bool>(
        context: context,
        builder: (context) => AlertDialog(
              title: const Text("Modifier son pseudo"),
              content: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  TextField(
                    autofocus: true,
                    textAlign: TextAlign.center,
                    controller: _controller,
                    decoration: const InputDecoration(
                        labelText: "Pseudo",
                        hintText: "Définit ton nom de joueur..."),
                  ),
                  const SizedBox(height: 12),
                  ElevatedButton(
                      onPressed: () {
                        Navigator.of(context).pop(true);
                      },
                      child: const Text("Enregistrer"))
                ],
              ),
            ));
    if (valid ?? false) widget.onSavePseudo(_controller.text);
  }
}
