import 'dart:convert';

import 'package:eleve/activities/automatismes/question.dart';
import 'package:eleve/activities/trivialpoursuit/login.dart';
import 'package:eleve/build_mode.dart';
import 'package:eleve/shared/activity_start.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/shared/settings_shared.dart';
import 'package:eleve/types/src_automatismes.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_sql_editor.dart';
import 'package:eleve/types/src_sql_trivial.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

abstract class AutomatismesAPI {
  Future<GetAutomatismesOut> load(GetAutomatismesIn args);
  Future<InstantiatedAutomatismeQuestion> instantiateQuestion(
    IdQuestiongroup id,
  );
  Future<QuestionAnswersOut> evaluateQuestion(EvaluateAutomatismeIn _);
}

class AutomatismesServerAPI implements AutomatismesAPI {
  final BuildMode buildMode;
  const AutomatismesServerAPI(this.buildMode);

  @override
  Future<GetAutomatismesOut> load(GetAutomatismesIn args) async {
    const serverEndpoint = "/api/student/automatismes";
    final uri = buildMode.serverURL(serverEndpoint);
    final resp = await http.post(
      uri,
      body: jsonEncode(getAutomatismesInToJson(args)),
      headers: {'Content-type': 'application/json'},
    );
    return getAutomatismesOutFromJson(checkServerError(resp.body));
  }

  @override
  Future<InstantiatedAutomatismeQuestion> instantiateQuestion(
    IdQuestiongroup id,
  ) async {
    const serverEndpoint = "/api/student/automatismes/question";
    final uri = buildMode.serverURL(
      serverEndpoint,
      query: {"idQuestion": id.toString()},
    );
    final resp = await http.get(
      uri,
      headers: {'Content-type': 'application/json'},
    );
    return instantiatedAutomatismeQuestionFromJson(checkServerError(resp.body));
  }

  @override
  Future<QuestionAnswersOut> evaluateQuestion(
    EvaluateAutomatismeIn args,
  ) async {
    const serverEndpoint = "/api/student/automatismes/question";
    final uri = buildMode.serverURL(serverEndpoint);
    final resp = await http.post(
      uri,
      body: jsonEncode(evaluateAutomatismeInToJson(args)),
      headers: {'Content-type': 'application/json'},
    );
    return questionAnswersOutFromJson(checkServerError(resp.body));
  }
}

class AutomatismesStart extends StatefulWidget {
  final AutomatismesAPI api;
  final UserSettings settings;

  const AutomatismesStart(this.api, this.settings, {super.key});

  @override
  State<AutomatismesStart> createState() => _AutomatismesStartState();
}

class _AutomatismesStartState extends State<AutomatismesStart> {
  var level = LevelTag.troisieme;
  var sublevel = "SPE";
  var isLoading = false;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text("Automatismes")),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            const Text("Choisir mon niveau"),
            SegmentedButton<LevelTag>(
              segments: [
                ButtonSegment(
                  value: LevelTag.troisieme,
                  label: Text(levelTagLabel(LevelTag.troisieme)),
                ),
                ButtonSegment(
                  value: LevelTag.premiere,
                  label: Text(levelTagLabel(LevelTag.premiere)),
                ),
              ],
              selected: {level},
              onSelectionChanged: (s) => setState(() {
                level = s.first;
              }),
            ),
            if (level == .premiere)
              SegmentedButton<String>(
                segments: [
                  ButtonSegment(value: "SPE", label: const Text("SPE")),
                  ButtonSegment(
                    value: "SANS SPE",
                    label: const Text("SANS SPE"),
                  ),
                  ButtonSegment(value: "TECHNO", label: const Text("TECHNO")),
                ],
                selected: {sublevel},
                onSelectionChanged: (s) => setState(() {
                  sublevel = s.first;
                }),
              ),
            ElevatedButton(
              onPressed: isLoading ? null : _load,
              child: const Text("Afficher"),
            ),
          ],
        ),
      ),
    );
  }

  void _load() async {
    setState(() {
      isLoading = true;
    });
    try {
      var res = await widget.api.load(GetAutomatismesIn(level, sublevel));
      if (!mounted) return;
      Navigator.of(context).push(
        MaterialPageRoute<void>(
          builder: (context) => _AutomatismesKindSelect(widget.api, res),
        ),
      );
    } catch (e) {
      showError("Chargement des données", e, context);
    }
  }
}

class _AutomatismesKindSelect extends StatelessWidget {
  final AutomatismesAPI api;
  final GetAutomatismesOut data;
  const _AutomatismesKindSelect(this.api, this.data);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text("Mode d'entrainement")),
      body: Column(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: [
          LaunchCard(
            "IsyTriv",
            "Je veux m'entrainer en jouant, seul ou à plusieurs.",
            const Icon(Icons.gamepad),
            data.trivials.isEmpty
                ? null
                : () => Navigator.of(context).push(
                    MaterialPageRoute<void>(
                      builder: (_) => _TrivialList(api, data.trivials),
                    ),
                  ),
          ),
          const Divider(thickness: 4),
          LaunchCard(
            "Questions",
            "Je veux m'entrainer sur une question.",
            const Icon(Icons.list),
            data.questions.isEmpty
                ? null
                : () => Navigator.of(context).push(
                    MaterialPageRoute<void>(
                      builder: (_) => _QuestionList(api, data.questions),
                    ),
                  ),
          ),
        ],
      ),
    );
  }
}

class _TrivialList extends StatelessWidget {
  final AutomatismesAPI api;
  final List<Trivial> list;

  const _TrivialList(this.api, this.list);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text("Choisir un Isy'Triv")),
      body: ListView(
        children: list
            .map((trivial) => TrivialRow(trivial, () => _launch(trivial)))
            .toList(),
      ),
    );
  }

  void _launch(Trivial trivial) {}
}

class _QuestionList extends StatefulWidget {
  final AutomatismesAPI api;
  final List<Questiongroup> list;

  const _QuestionList(this.api, this.list);

  @override
  State<_QuestionList> createState() => _QuestionListState();
}

class _QuestionListState extends State<_QuestionList> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text("Choisir une question")),
      body: ListView(
        children: widget.list
            .map(
              (question) => ListTile(
                title: Text(question.title),
                onTap: () => _load(question),
              ),
            )
            .toList(),
      ),
    );
  }

  void _load(Questiongroup question) async {
    showDialog<void>(
      context: context,
      builder: (_) => AlertDialog(title: const Text("Chargement")),
      barrierDismissible: false,
    );
    try {
      final res = await widget.api.instantiateQuestion(question.id);
      if (!mounted) return;
      Navigator.of(context).pop(); // remove dialog
      Navigator.of(context).push(
        MaterialPageRoute<void>(
          builder: (context) => AutomatismeQuestionW(res, widget.api, () {
            Navigator.of(context).pop();
            _load(question);
          }),
        ),
      );
    } catch (e) {
      showError("Chargement des données", e, context);
    }
  }
}

class AutomatismesActivityIcon extends StatelessWidget {
  final void Function() onTap;

  const AutomatismesActivityIcon(this.onTap, {super.key});

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        RawMaterialButton(
          onPressed: onTap,
          child: Image.asset(
            "assets/images/automatismes.png",
            width: 68,
            height: 60,
          ),
        ),
        const Padding(
          padding: EdgeInsets.only(top: 8, bottom: 6),
          child: Text("Automatismes"),
        ),
      ],
    );
  }
}
