import 'dart:convert';

import 'package:eleve/activities/automatismes/question.dart';
import 'package:eleve/activities/trivialpoursuit/login.dart';
import 'package:eleve/build_mode.dart';
import 'package:eleve/shared/activity_start.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/types/src_automatismes.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_prof_trivial.dart';
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

  Future<LaunchSelfaccessOut> launchTrivial(
    String clientId,
    IdTrivial idTrivial,
  );
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
  Future<LaunchSelfaccessOut> launchTrivial(
    String clientId,
    IdTrivial idTrivial,
  ) async {
    const serverEndpoint = "/api/student/automatismes/trivials/launch";
    final uri = buildMode.serverURL(
      serverEndpoint,
      query: {"client-id": clientId, "trivial-id": idTrivial.toString()},
    );
    final resp = await http.get(
      uri,
      headers: {'Content-type': 'application/json'},
    );
    return launchSelfaccessOutFromJson(checkServerError(resp.body));
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
  final TrivialSettings settings;
  const AutomatismesStart(this.api, this.settings, {super.key});

  @override
  State<AutomatismesStart> createState() => _AutomatismesStartState();
}

class _AutomatismesStartState extends State<AutomatismesStart> {
  var level = LevelTag.premiere;
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
                // properly reset the sublevel
                sublevel = "";
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
            ElevatedButton.icon(
              onPressed: isLoading ? null : _load,
              label: const Text("Afficher"),
              icon: isLoading
                  ? SizedBox(
                      width: 16,
                      height: 16,
                      child: CircularProgressIndicator(),
                    )
                  : null,
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

      final goToTrivials = await Navigator.of(context).push(
        MaterialPageRoute<bool>(
          builder: (context) =>
              _AutomatismesKindSelect(widget.api, widget.settings, res),
        ),
      );
      if ((goToTrivials ?? false) && mounted) Navigator.of(context).pop(true);
    } catch (e) {
      if (mounted) showError("Chargement des données", e, context);
    }
    setState(() {
      isLoading = false;
    });
  }
}

class _AutomatismesKindSelect extends StatelessWidget {
  final AutomatismesAPI api;
  final TrivialSettings settings;
  final GetAutomatismesOut data;
  const _AutomatismesKindSelect(this.api, this.settings, this.data);

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
                : () async {
                    final goToTrivials = await Navigator.of(context).push(
                      MaterialPageRoute<bool>(
                        builder: (_) =>
                            _TrivialList(api, settings, data.trivials),
                      ),
                    );
                    if ((goToTrivials ?? false) && context.mounted) {
                      Navigator.of(context).pop(true);
                    }
                  },
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

class _TrivialList extends StatefulWidget {
  final AutomatismesAPI api;
  final TrivialSettings settings;
  final List<Trivial> list;

  const _TrivialList(this.api, this.settings, this.list);

  @override
  State<_TrivialList> createState() => _TrivialListState();
}

class _TrivialListState extends State<_TrivialList> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text("Choisir un Isy'Triv")),
      body: Column(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Expanded(
            child: ListView(
              children: widget.list
                  .map((trivial) => TrivialRow(trivial, () => _launch(trivial)))
                  .toList(),
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: ElevatedButton.icon(
              onPressed: () => Navigator.of(context).pop(true),
              label: Text("Rejoindre une partie en cours"),
              icon: const Icon(Icons.key),
            ),
          ),
        ],
      ),
    );
  }

  void _launch(Trivial trivial) async {
    showDialog<void>(
      context: context,
      builder: (_) => AlertDialog(
        title: const Text("Lancement de la partie"),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [CircularProgressIndicator()],
        ),
      ),
      barrierDismissible: false,
    );
    try {
      final res = await widget.api.launchTrivial(
        widget.settings.settings.settings.studentID,
        trivial.id,
      );
      if (!mounted) return;
      Navigator.of(context).pop(); // remove dialog
      Navigator.of(context).push(launchGameRoute(res, widget.settings));
    } catch (e) {
      showError("Chargement des données", e, context);
      Navigator.of(context).pop(); // remove dialog
    }
  }
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
      builder: (_) => AlertDialog(
        title: const Text("Chargement"),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [CircularProgressIndicator()],
        ),
      ),
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
          child: Text("AutoMath'ismes"),
        ),
      ],
    );
  }
}
