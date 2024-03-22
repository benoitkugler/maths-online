import 'package:eleve/types/src_sql_events.dart';
import 'package:flutter/material.dart';

class AdvanceSummary extends StatelessWidget {
  final Stats advance;

  const AdvanceSummary(this.advance, {super.key});

  @override
  Widget build(BuildContext context) {
    final color = _rankColors[advance.rank];
    return InkWell(
      onTap: () => Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (context) => _EventsDetailsView(advance),
      )),
      child: Column(
        children: [
          Card(
              shadowColor: color,
              elevation: 6,
              color: color,
              child: Padding(
                padding: const EdgeInsets.all(8.0),
                child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                    children: [
                      Text(_rankLabels[advance.rank],
                          style: Theme.of(context).textTheme.titleMedium),
                      Chip(
                          backgroundColor: Colors.teal,
                          label: Row(
                            children: [
                              Text("${advance.totalPoints}"),
                              const SizedBox(width: 8),
                              const Icon(Icons.check, size: 18),
                            ],
                          )),
                    ]),
              )),
          _rankIcon(advance.rank),
        ],
      ),
    );
  }
}

const _rankLabels = [
  "Novice",
  "Guilde de Pythagore",
  "Guilde de Thalès",
  "Guilde d'Al-Kashi",
  "Guilde de Newton",
  "Guilde de Gauss",
  "Guilde d'Einstein",
];

const _rankColors = [
  Colors.transparent,
  Color.fromARGB(255, 158, 122, 55),
  Colors.grey,
  Color.fromARGB(255, 148, 112, 100),
  Color.fromARGB(255, 202, 219, 218),
  Color.fromARGB(255, 216, 197, 23),
  Color.fromARGB(255, 112, 241, 235)
];

Image _rankIcon(int rank) {
  return Image.asset("assets/images/ranks/rank-$rank.png", width: 80);
}

class _EventsDetailsView extends StatelessWidget {
  final Stats advance;

  const _EventsDetailsView(this.advance, {super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text("Détails de tes succès"),
      ),
      body: ListView(
          children: EventK.values
              .map((e) => _EventOccurences(e, advance.occurences[e.index]))
              .toList()),
    );
  }
}

class _EventOccurences extends StatelessWidget {
  final EventK event;
  final int occurences;
  const _EventOccurences(this.event, this.occurences, {super.key});

  @override
  Widget build(BuildContext context) {
    return ListTile(
      leading: occurences == 0
          ? const Text("-")
          : Chip(backgroundColor: Colors.teal, label: Text("$occurences")),
      title: Text(eventKLabel(event)),
      titleAlignment: ListTileTitleAlignment.center,
    );
  }
}
