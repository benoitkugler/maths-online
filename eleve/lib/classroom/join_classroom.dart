import 'dart:convert';

import 'package:eleve/build_mode.dart';
import 'package:eleve/shared/date_field.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/shared/pin.dart';
import 'package:eleve/shared/settings_shared.dart';
import 'package:eleve/types/src_prof_teacher.dart';
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
  String code = "";

  void _onValidCode(String code) async {
    this.code = code;
    try {
      final uri = widget.buildMode
          .serverURL("/api/classroom/attach", query: {"code": code});
      final resp = await http.get(uri);
      setState(() {
        studentProposals =
            listStudentHeaderFromJson(checkServerError(resp.body));
      });
    } catch (e) {
      _showError(e);
    }
  }

  void _showError(dynamic error) {
    showError("Impossible de rejoindre la classe.", error, context);
  }

  void _onSelected(StudentHeader student) async {
    final date = await showDialog<String>(
        context: context,
        builder: (context) {
          return Dialog(
            child: Card(
              child: Padding(
                padding:
                    const EdgeInsets.symmetric(horizontal: 8.0, vertical: 12),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    const Text("Confirme en entrant ta date de naissance"),
                    DateField((date) {
                      Navigator.of(context).pop(date);
                    })
                  ],
                ),
              ),
            ),
          );
        });
    if (date != null) {
      _validAttach(student, date);
    }
  }

  void _validAttach(StudentHeader student, String date) async {
    final device = await loadUserDeviceName();
    try {
      final uri = widget.buildMode.serverURL("/api/classroom/attach");
      final args = AttachStudentToClassroom2In(
        code,
        student.id,
        date,
        device,
      );
      final resp = await http.post(uri,
          body: jsonEncode(attachStudentToClassroom2InToJson(args)),
          headers: {
            'Content-type': 'application/json',
          });
      final result =
          attachStudentToClassroom2OutFromJson(checkServerError(resp.body));
      if (result.errInvalidBirthday) {
        _showError("Date de naissance invalide.");
        return;
      }

      // pop the route with the result
      if (!mounted) return;
      Navigator.of(context).pop(result.idCrypted);
    } catch (e) {
      _showError(e);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(title: const Text("Rejoindre une classe")),
        body: studentProposals.isEmpty
            ? Pin("Code de la classe", _onValidCode)
            : Column(
                children: [
                  const Padding(
                    padding: EdgeInsets.symmetric(vertical: 20.0),
                    child: Text(
                      "Qui es-tu ?",
                      style: TextStyle(fontSize: 18),
                    ),
                  ),
                  Expanded(
                    child: ListView(
                        children: studentProposals
                            .map((student) => ListTile(
                                  title: Text(student.label),
                                  onTap: () => _onSelected(student),
                                ))
                            .toList()),
                  ),
                ],
              ));
  }
}

Future<bool> _leaveClassroom(BuildMode bm, String idCrypted) async {
  final uri =
      bm.serverURL("/api/classroom/attach", query: {"id-crypted": idCrypted});
  final resp = await http.delete(uri);
  return resp.statusCode == 200;
}

Future<bool> confirmLeaveClassroom(
    BuildMode bm, String idCrypted, BuildContext context) async {
  final ok = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
            title: const Text("Confirmer"),
            content: const Text("Es-tu sûr(e) de quitter la classe ?"),
            actions: [
              TextButton(
                  onPressed: () async {
                    try {
                      final ok = await _leaveClassroom(bm, idCrypted);
                      Navigator.of(context).pop(ok);
                    } catch (e) {
                      showError("Erreur", e, context);
                    }
                  },
                  child: const Text("Quitter"))
            ],
          ));

  return ok ?? false;
}
