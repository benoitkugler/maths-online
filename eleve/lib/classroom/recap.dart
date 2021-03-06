// used for registred students

import 'dart:convert';

import 'package:eleve/build_mode.dart';
import 'package:eleve/settings.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/shared/students.gen.dart';
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

class _ClassroomCardState extends State<ClassroomCard> {
  late Future<StudentClassroomHeader> meta;

  @override
  void initState() {
    meta = loadMeta();

    super.initState();
  }

  Future<StudentClassroomHeader> loadMeta() async {
    final uri = Uri.parse(widget.buildMode.serverURL("/api/classroom/login",
        query: {studentIDKey: widget.studentID}));
    final resp = await http.get(uri);
    final res = checkStudentClassroomOutFromJson(jsonDecode(resp.body));
    if (res.isOK) {
      return res.meta;
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
    return Container(
      padding: const EdgeInsets.all(10),
      decoration: BoxDecoration(
        border: Border.all(color: Colors.lightBlue, width: 2),
        borderRadius: const BorderRadius.all(Radius.circular(6)),
      ),
      child: FutureBuilder<StudentClassroomHeader>(
          future: meta,
          builder: (context, snapshot) {
            if (snapshot.hasError) {
              return ErrorBar(
                  "Impossible de charger les données.", snapshot.error!);
            }

            if (snapshot.hasData) {
              return _LoadedClassroom(snapshot.data!);
            }
            return const _Loading();
          }),
    );
  }
}

class _Loading extends StatelessWidget {
  const _Loading({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Column(children: const [
      CircularProgressIndicator(value: null),
      SizedBox(height: 20),
      Text("Chargement des données..."),
    ]);
  }
}

class _LoadedClassroom extends StatelessWidget {
  final StudentClassroomHeader meta;

  const _LoadedClassroom(this.meta, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        _FormRow("Nom", "${meta.student.surname} ${meta.student.name}"),
        _FormRow("Classe", meta.classroomName),
        _FormRow("Mail du prof.", meta.teacherMail),
      ],
    );
  }
}

class _FormRow extends StatelessWidget {
  final String left;
  final String right;
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
            style: const TextStyle(fontSize: 18),
          ),
          Text(
            right,
            style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
          ),
        ],
      ),
    );
  }
}
