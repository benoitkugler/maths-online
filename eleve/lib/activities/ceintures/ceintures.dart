import 'package:eleve/activities/ceintures/animations.dart';
import 'package:eleve/activities/ceintures/api.dart';
import 'package:eleve/activities/ceintures/seance.dart';
import 'package:eleve/exercice/congratulations.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/shared/settings_shared.dart';
import 'package:eleve/types/src_prof_ceintures.dart';
import 'package:eleve/types/src_sql_ceintures.dart';
import 'package:flutter/material.dart';

class CeinturesActivityIcon extends StatelessWidget {
  final void Function() onTap;

  const CeinturesActivityIcon(this.onTap, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        RawMaterialButton(
          onPressed: onTap,
          child: Image.asset("assets/images/yellow-belt.png",
              width: 68, height: 60),
        ),
        const Padding(
          padding: EdgeInsets.only(top: 8, bottom: 6),
          child: Text("Ceintures de calcul"),
        ),
      ],
    );
  }
}

class CeinturesStart extends StatefulWidget {
  final CeinturesAPI api;
  final UserSettings settings;
  final void Function(String id) saveAnonymousID;
  const CeinturesStart(this.api, this.settings, this.saveAnonymousID,
      {super.key});

  StudentTokens get tokens =>
      StudentTokens(settings.ceinturesAnonymousID, settings.studentID);

  @override
  State<CeinturesStart> createState() => _CeinturesStartState();
}

class _CeinturesStartState extends State<CeinturesStart> {
  late Future<GetEvolutionOut> loader;

  @override
  void initState() {
    loader = widget.api.getEvolution(widget.tokens);
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(title: const Text("Ceintures de calcul")),
        body: FutureBuilder<GetEvolutionOut>(
            future: loader,
            builder: (context, snapshot) {
              final data = snapshot.data;
              return snapshot.error != null
                  ? ErrorCard(
                      "Impossible de charger la progression.", snapshot.error)
                  : data == null
                      ? const Center(child: CircularProgressIndicator())
                      : data.has
                          ? _CeinturesEvolution(
                              widget.api, widget.tokens, data.evolution)
                          : _CreateEvolution(createEvolution);
            }));
  }

  void createEvolution(Level level) async {
    final CreateEvolutionOut res;
    try {
      res = await widget.api
          .createEvolution(CreateEvolutionIn(widget.settings.studentID, level));
    } catch (e) {
      if (!mounted) return;
      showError("Impossible de créer le parcours.", e, context);
      return;
    }
    widget.saveAnonymousID(res.anonymousID);
    setState(() {
      loader = Future.sync(() => GetEvolutionOut(true, res.evolution));
    });
  }
}

class _CreateEvolution extends StatefulWidget {
  final void Function(Level level) onCreate;
  const _CreateEvolution(this.onCreate, {super.key});

  @override
  State<_CreateEvolution> createState() => _CreateEvolutionState();
}

class _CreateEvolutionState extends State<_CreateEvolution> {
  Level? level;
  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Padding(
          padding: const EdgeInsets.all(8.0),
          child: Text(
            "Choisis ton niveau",
            style: Theme.of(context).textTheme.titleLarge,
          ),
        ),
        Column(
            mainAxisSize: MainAxisSize.min,
            children: Level.values
                .map((e) => RadioListTile<Level>(
                      title: Text(levelLabel(e)),
                      value: e,
                      groupValue: level,
                      onChanged: (value) => setState(() {
                        level = value ?? Level.seconde;
                      }),
                    ))
                .toList()),
        ElevatedButton(
            onPressed: level == null ? null : () => widget.onCreate(level!),
            child: const Text("Démarrer !"))
      ],
    );
  }
}

class _CeinturesEvolution extends StatefulWidget {
  final CeinturesAPI api;
  final StudentTokens tokens;
  final StudentEvolution initialEvolution;

  const _CeinturesEvolution(this.api, this.tokens, this.initialEvolution,
      {super.key});

  @override
  State<_CeinturesEvolution> createState() => _CeinturesEvolutionState();
}

class _CeinturesEvolutionState extends State<_CeinturesEvolution> {
  late StudentEvolution evolution;

  @override
  void initState() {
    _init();
    super.initState();
  }

  @override
  void didUpdateWidget(covariant _CeinturesEvolution oldWidget) {
    _init();
    super.didUpdateWidget(oldWidget);
  }

  void _init() {
    evolution = widget.initialEvolution;
  }

  @override
  Widget build(BuildContext context) {
    final scheme = evolution.scheme;
    final level = evolution.level;
    final pending = Map.fromEntries(
        evolution.pending.map((p) => MapEntry(p.domain, p.rank)));

    return GridView.count(
        crossAxisCount: 2,
        padding: const EdgeInsets.all(4),
        children: Domain.values
            .where((d) => level.index >= scheme.levels[d.index].index)
            .map((d) {
          final rank = evolution.advance[d.index];
          final nextStat = rank.next == null
              ? null
              : evolution.stats[d.index][rank.next!.index];
          final stage = Stage(d, rank);
          return _StageTile(scheme, stage, !pending.containsKey(d), nextStat,
              () => _launch(d, pending[d]!));
        }).toList());
  }

  void _launch(Domain domain, Rank pending) {
    final stage = Stage(domain, pending);
    Navigator.of(context).push(MaterialPageRoute<void>(
      builder: (context) => Seance(
        widget.api,
        widget.tokens,
        stage,
        (b, ev) => _onValid(stage, b, ev),
      ),
    ));
  }

  void _showSucces(Stage stage) {
    showDialog<void>(
        context: context,
        builder: (context) => AlertDialog(
            title: const Text("Ceinture débloquée !"),
            content: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                UnlockAnimation(stage.rank.color, domainLabel(stage.domain)),
                const SizedBox(height: 50),
                Text(pickCongratulationMessage(),
                    style: const TextStyle(fontStyle: FontStyle.italic)),
              ],
            )));
  }

  void _onValid(Stage stage, bool isSucces, StudentEvolution newEvolution) {
    setState(() {
      evolution = newEvolution;
    });
    if (isSucces) {
      Navigator.of(context).pop(); // remove the question view
      _showSucces(stage);
    }
  }
}

class _StageTile extends StatelessWidget {
  final Scheme scheme;
  final Stage stage;
  final bool locked;
  final Stat? pendingStat;
  final void Function() onLaunch;
  const _StageTile(
      this.scheme, this.stage, this.locked, this.pendingStat, this.onLaunch,
      {super.key});

  @override
  Widget build(BuildContext context) {
    final pending = stage.rank.next;
    return InkWell(
      borderRadius: BorderRadius.circular(4),
      onTap: () => _showStage(context),
      child: Card(
        elevation: locked ? 0 : 4,
        shadowColor: locked ? Colors.grey : Colors.white,
        color: Colors.teal.withOpacity(0.2),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            Text(
              domainLabel(stage.domain),
              style: Theme.of(context).textTheme.headlineSmall,
            ),
            pending == null
                ? const Icon(Icons.check, size: 48)
                : locked
                    ? const Icon(Icons.lock, size: 32)
                    : CeintureIcon(
                        pending,
                        scale: 1.3,
                        withBackground: true,
                      )
          ],
        ),
      ),
    );
  }

  void _showStage(BuildContext context) async {
    final launch = await showDialog<bool>(
        context: context,
        builder: (context) =>
            _StageDetailsDialog(scheme, stage, locked, pendingStat));
    if (launch == true) onLaunch();
  }
}

extension on Stat {
  int get essais => failure + success;
  int get reussite => essais == 0 ? 0 : (100 * success / essais).round();
}

class _StageDetailsDialog extends StatelessWidget {
  final Scheme scheme;
  final Stage stage;
  final bool locked;
  final Stat? pendingStat;
  const _StageDetailsDialog(
      this.scheme, this.stage, this.locked, this.pendingStat,
      {super.key});

  @override
  Widget build(BuildContext context) {
    final pending = stage.rank.next;

    final needed = scheme.ps
        .where((s) =>
            s.pending.domain == stage.domain && s.pending.rank == pending)
        .map((p) => p.need)
        .toList();

    return AlertDialog(
      icon: pending == null
          ? const Icon(Icons.check)
          : locked
              ? const Icon(Icons.lock)
              : null,
      title: Text(domainLabel(stage.domain)),
      content: Column(
          mainAxisSize: MainAxisSize.min,
          children: pending == null
              ? const [
                  Text("Domaine terminé. Bravo !"),
                ]
              : locked
                  ? [
                      const Text(
                          "Tu as besoin de réussir d'abord les ceintures suivantes :"),
                      const SizedBox(height: 16),
                      Wrap(
                          alignment: WrapAlignment.center,
                          children: needed
                              .map((e) => Card(
                                  shape: const RoundedRectangleBorder(
                                      borderRadius:
                                          BorderRadius.all(Radius.circular(8))),
                                  child: Padding(
                                    padding: const EdgeInsets.all(8.0),
                                    child: Row(
                                      mainAxisSize: MainAxisSize.min,
                                      children: [
                                        CeintureIcon(e.rank),
                                        const SizedBox(width: 10),
                                        Text(domainLabel(e.domain)),
                                      ],
                                    ),
                                  )))
                              .toList()),
                    ]
                  : [
                      if (stage.rank != Rank.startRank) ...[
                        const Text("Niveau actuel"),
                        const SizedBox(height: 12),
                        CeintureIcon(stage.rank, withBackground: true),
                        const SizedBox(height: 12),
                        const Icon(Icons.arrow_downward),
                        const SizedBox(height: 12),
                      ],
                      const Text("Prochain niveau"),
                      const SizedBox(height: 12),
                      CeintureIcon(
                        pending,
                        withBackground: true,
                      ),
                      const SizedBox(height: 16),
                      if (pendingStat != null)
                        Text(
                          "${pendingStat!.reussite}% de réussite (sur ${pendingStat!.essais} réponses)",
                          style: Theme.of(context).textTheme.labelMedium,
                        ),
                      const SizedBox(height: 32),
                      ElevatedButton(
                          style: ElevatedButton.styleFrom(
                              elevation: 2, shadowColor: Colors.teal.shade200),
                          child: const Text("Lancer la séance !"),
                          onPressed: () => Navigator.of(context).pop(true)),
                    ]),
    );
  }
}
