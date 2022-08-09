import 'dart:convert';

import 'package:eleve/build_mode.dart';
import 'package:eleve/homework/sheet.dart';
import 'package:eleve/homework/types.gen.dart';
import 'package:eleve/settings.dart';
import 'package:eleve/shared/errors.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

const _serverEndpoint = "/api/student/homework/sheets";

typedef Sheets = List<SheetProgression>;

/// Homework is the entry point for the homework
/// activity
class Homework extends StatefulWidget {
  final BuildMode buildMode;
  final String studentID;

  const Homework(this.buildMode, this.studentID, {Key? key}) : super(key: key);

  @override
  State<Homework> createState() => _HomeworkState();
}

class _HomeworkState extends State<Homework> {
  late Future<Sheets> sheets;

  @override
  void initState() {
    sheets = loadSheets();

    super.initState();
  }

  Future<Sheets> loadSheets() async {
    final uri = Uri.parse(widget.buildMode
        .serverURL(_serverEndpoint, query: {studentIDKey: widget.studentID}));
    final resp = await http.get(uri);
    return listSheetProgressionFromJson(jsonDecode(resp.body));
  }

  @override
  void didUpdateWidget(covariant Homework oldWidget) {
    sheets = loadSheets();
    super.didUpdateWidget(oldWidget);
  }

  void showSummary(Sheets data) async {
    final sheet = await Navigator.of(context)
        .push(MaterialPageRoute<SheetProgression>(builder: (context) {
      return _SheetList(data, (selected) => Navigator.pop(context, selected));
    }));
    if (sheet != null) {
      // TODO: go to this sheet
    }
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<Sheets>(
      future: sheets,
      builder: (context, snapshot) {
        final Widget body;
        if (snapshot.hasError) {
          body = Center(
              child: ErrorBar(
                  "Impossible de charger les données.", snapshot.error!));
        } else if (snapshot.hasData) {
          body = SheetHome(widget.buildMode, snapshot.data!);
        } else {
          body = const _Loading();
        }
        return Scaffold(
          appBar: AppBar(
            actions: [
              TextButton(
                  onPressed: snapshot.hasData
                      ? () => showSummary(snapshot.data!)
                      : null,
                  child: const Text("Fiches"))
            ],
          ),
          body: body,
        );
      },
    );
  }
}

//   @override
//   Widget build(BuildContext context) {
//     return Container(
//       padding: const EdgeInsets.all(10),
//       decoration: BoxDecoration(
//         border: Border.all(color: Colors.lightBlue, width: 2),
//         borderRadius: const BorderRadius.all(Radius.circular(6)),
//       ),

//   }
// }

class _Loading extends StatelessWidget {
  const _Loading({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Column(children: const [
      CircularProgressIndicator(value: null),
      SizedBox(height: 20),
      Text("Chargement des fiches..."),
    ]);
  }
}

// shows the most interesting sheet, or nothing
class SheetHome extends StatelessWidget {
  final BuildMode buildMode;
  final Sheets sheets;
  const SheetHome(this.buildMode, this.sheets, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    if (sheets.isEmpty) {
      return const Center(child: Text("Aucun devoir n'est planifié."));
    }

    // select the most recent one
    sheets
        .sort(((a, b) => a.sheet.deadline.isAfter(b.sheet.deadline) ? 1 : -1));
    return SheetW(buildMode, sheets[0]);
  }
}

class _SheetList extends StatelessWidget {
  final Sheets sheets;
  final void Function(SheetProgression) onSelect;
  const _SheetList(this.sheets, this.onSelect, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text("Liste des fiches"),
      ),
      body: ListView(
        children: sheets
            .map((e) => GestureDetector(
                onTap: () => onSelect(e), child: _SheetSummary(e)))
            .toList(),
      ),
    );
  }
}

class _SheetSummary extends StatelessWidget {
  final SheetProgression sheet;
  const _SheetSummary(this.sheet, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Text(sheet.sheet.title);
  }
}

class HomeworkActivityIcon extends StatelessWidget {
  final void Function() onTap;

  const HomeworkActivityIcon(this.onTap, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        GestureDetector(
          onTap: onTap,
          child: Image.asset("lib/images/homework.png", width: 68),
        ),
        const Padding(
          padding: EdgeInsets.only(top: 8, bottom: 6),
          child: Text("Travail à la maison"),
        ),
      ],
    );
  }
}
