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
                          ? _Evolution(
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

class _Evolution extends StatefulWidget {
  final CeinturesAPI api;
  final StudentTokens tokens;
  final StudentEvolution initialEvolution;

  const _Evolution(this.api, this.tokens, this.initialEvolution, {super.key});

  @override
  State<_Evolution> createState() => _EvolutionState();
}

class _EvolutionState extends State<_Evolution> {
  late StudentEvolution evolution;
  Domain? selected;

  Map<Domain, Rank> get pending =>
      Map.fromEntries(evolution.pending.map((p) => MapEntry(p.domain, p.rank)));

  Domain? get suggested => (evolution.suggestionIndex != -1)
      ? evolution.pending[evolution.suggestionIndex].domain
      : null;

  @override
  void initState() {
    _init();
    super.initState();
  }

  @override
  void didUpdateWidget(covariant _Evolution oldWidget) {
    _init();
    super.didUpdateWidget(oldWidget);
  }

  void _init() {
    evolution = widget.initialEvolution;
    selected = suggested;
  }

  @override
  Widget build(BuildContext context) {
    final scheme = evolution.scheme;
    final level = evolution.level;
    return Column(
      children: [
        Expanded(
            child: ListView(
                children: Domain.values
                    .where((d) => level.index >= scheme.levels[d.index].index)
                    .map((d) {
          final rank = evolution.advance[d.index];
          final nextStat = rank.next == null
              ? null
              : evolution.stats[d.index][rank.next!.index];
          final locked = !pending.containsKey(d);
          return _DomainTile(
              d, rank, nextStat, d == selected, d == suggested, locked, () {
            locked ? _showLockInfo(d) : setState(() => selected = d);
          });
        }).toList())),
        Padding(
          padding: const EdgeInsets.all(8.0),
          child: ElevatedButton(
              onPressed: selected == null ? null : _launch,
              child: const Text("Lancer la séance !")),
        )
      ],
    );
  }

  void _showLockInfo(Domain d) {
    final current = evolution.advance[d.index];
    final next = current.next;
    if (next == null) return;
    showDialog<void>(
        context: context,
        builder: (context) =>
            _LockInfoDialog(evolution.scheme, Stage(d, next)));
  }

  void _launch() {
    final stage = Stage(
      selected!,
      pending[selected!]!,
    );
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
      selected = suggested;
    });
    if (isSucces) {
      Navigator.of(context).pop(); // remove the question view
      _showSucces(stage);
    }
  }
}

extension on Stat {
  int get essais => failure + success;
  int get reussite => essais == 0 ? 0 : (100 * success / essais).round();
}

class _DomainTile extends StatelessWidget {
  final Domain domain;
  final Rank rank;
  final Stat? stat; // null at the end
  final bool selected;
  final bool recommanded;
  final bool locked;
  final void Function() onTap;
  const _DomainTile(this.domain, this.rank, this.stat, this.selected,
      this.recommanded, this.locked, this.onTap,
      {super.key});

  @override
  Widget build(BuildContext context) {
    return ListTile(
        enabled: stat != null,
        onTap: onTap,
        selected: selected,
        leading: rank == Rank.startRank
            ? const Icon(Icons.arrow_right)
            : CeintureIcon(rank),
        title: Text(domainLabel(domain)),
        trailing: stat == null
            ? const Icon(Icons.check)
            : recommanded
                ? const Icon(Icons.star)
                : locked
                    ? const Icon(Icons.lock)
                    : null,
        subtitle: stat == null
            ? const Text("Parcours terminé.")
            : Text("${stat!.reussite}% de réussite - ${stat!.essais} essais"));
  }
}

extension on Stage {
  bool equals(Stage other) => domain == other.domain && rank == other.rank;
}

class _LockInfoDialog extends StatelessWidget {
  final Scheme scheme;
  final Stage stage;
  const _LockInfoDialog(this.scheme, this.stage, {super.key});

  List<Stage> get needed => scheme.ps
      .where((element) => element.pending.equals(stage))
      .map((p) => p.need)
      .toList();

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      icon: const Icon(Icons.lock),
      title: const Text("Séance verrouillée !"),
      content: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          const Text(
              "Tu as besoin de réussir d'abord les ceintures suivantes :"),
          const SizedBox(height: 16),
          Wrap(
              alignment: WrapAlignment.center,
              children: needed
                  .map((e) => Card(
                      shape: const RoundedRectangleBorder(
                          borderRadius: BorderRadius.all(Radius.circular(8))),
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
        ],
      ),
    );
  }
}
