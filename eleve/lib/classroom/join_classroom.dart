import 'dart:convert';

import 'package:eleve/build_mode.dart';
import 'package:eleve/shared/pin.dart';
import 'package:eleve/shared/students.gen.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

class JoinClassroomRoute extends StatefulWidget {
  final BuildMode buildMode;
  const JoinClassroomRoute(this.buildMode, {Key? key}) : super(key: key);

  @override
  State<JoinClassroomRoute> createState() => _JoinClassroomRouteState();
}

class _JoinClassroomRouteState extends State<JoinClassroomRoute> {
  List<StudentHeader> studentProposals = [];
  final codeController = TextEditingController();

  void _onValidCode(String code) async {
    try {
      final uri = Uri.parse(widget.buildMode
          .serverURL("/api/classroom/attach", query: {"code": code}));
      final resp = await http.get(uri);
      setState(() {
        studentProposals = listStudentHeaderFromJson(jsonDecode(resp.body));
      });
    } catch (e) {
      _showError(e);
    }
  }

  void _showError(dynamic error) {
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      duration: const Duration(seconds: 5),
      backgroundColor: Theme.of(context).colorScheme.error,
      content: Text("Une erreur est survenue : $error"),
    ));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(title: const Text("Rejoindre une classe")),
        body: studentProposals.isEmpty
            ? Pin("Code de la classe", codeController, _onValidCode)
            : ListView(
                children: studentProposals
                    .map((student) => ListTile(
                          title: Text(student.label),
                        ))
                    .toList()));
  }
}
