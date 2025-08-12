// used for registred students

import 'dart:convert';

import 'package:eleve/build_mode.dart';
import 'package:eleve/classroom/student_advance.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/shared/hyperlink.dart';
import 'package:eleve/shared/settings_shared.dart';
import 'package:eleve/types/src_prof_teacher.dart';
import 'package:eleve/types/src_sql_events.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

class ClassroomCard extends StatefulWidget {
  final BuildMode buildMode;
  final String studentID;
  final void Function() onInvalidID;

  const ClassroomCard(this.buildMode, this.studentID, this.onInvalidID,
      {Key? key})
      : super(key: key);

  @override
  State<ClassroomCard> createState() => _ClassroomCardState();
}

class _Profile {
  final StudentClassroomHeader meta;
  final StudentAdvance advance;
  const _Profile(this.meta, this.advance);
}

class _ClassroomCardState extends State<ClassroomCard> {
  late Future<_Profile> meta;

  @override
  void initState() {
    meta = loadMeta();

    super.initState();
  }

  Future<_Profile> loadMeta() async {
    final uri = widget.buildMode.serverURL("/api/classroom/login",
        query: {studentIDKey: widget.studentID});
    final resp = await http.get(uri);
    final res = checkStudentClassroomOutFromJson(jsonDecode(resp.body));
    if (res.isOK) {
      return _Profile(res.meta, res.advance);
    }
    widget.onInvalidID();
    throw ("L'identifiant n'est plus valide");
  }

  @override
  void didUpdateWidget(covariant ClassroomCard oldWidget) {
    try {
      meta = loadMeta();
    } catch (e) {
      showError("Impossible de charger les données.", e, context);
    }

    super.didUpdateWidget(oldWidget);
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<_Profile>(
        future: meta,
        builder: (context, snapshot) {
          if (snapshot.hasError) {
            return ErrorBar(
                "Impossible de charger les données.", snapshot.error!);
          }

          if (snapshot.hasData) {
            return _LoadedProfile(snapshot.data!);
          }
          return const _Loading();
        });
  }
}

class _Loading extends StatelessWidget {
  const _Loading({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return const Column(children: [
      CircularProgressIndicator(value: null),
      SizedBox(height: 20),
      Text("Chargement des données..."),
    ]);
  }
}

class _LoadedProfile extends StatelessWidget {
  final _Profile profile;

  const _LoadedProfile(this.profile, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final meta = profile.meta;
    const detailsStyle = TextStyle(fontSize: 14, fontWeight: FontWeight.bold);

    return Padding(
      padding: const EdgeInsets.all(8.0),
      child: Column(
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.start,
            children: [
              const Icon(Icons.account_circle),
              const SizedBox(width: 12),
              Text(
                "Compte enregistré",
                style: Theme.of(context).textTheme.titleLarge,
              ),
            ],
          ),
          const SizedBox(height: 24),
          AdvanceSummary(profile.advance),
          const SizedBox(height: 12),
          const Divider(thickness: 2),
          const SizedBox(height: 12),
          _FormRow(
              "Nom",
              Text("${meta.student.surname} ${meta.student.name}",
                  style: detailsStyle)),
          _FormRow("Classe", Text(meta.classroomName, style: detailsStyle)),
          ...meta.teachers.map((t) => _FormRow(
              "Contact", hyperlink(t.mail, t.contactURL, style: detailsStyle))),
        ],
      ),
    );
  }
}

class _FormRow extends StatelessWidget {
  final String left;
  final Widget right;
  const _FormRow(this.left, this.right, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4.0),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            left,
            style: const TextStyle(fontSize: 10),
          ),
          right
        ],
      ),
    );
  }
}
