import 'package:eleve/shared/progression_bar.dart';
import 'package:eleve/types/src_sql_events.dart';
import 'package:flutter/material.dart';

class AdvanceSummary extends StatelessWidget {
  final StudentAdvance advance;

  const AdvanceSummary(this.advance, {super.key});

  @override
  Widget build(BuildContext context) {
    final advance = StudentAdvance([], 100, 2, 1, 0, 300);
    final color = _rankColors[advance.rank];
    return InkWell(
      borderRadius: BorderRadius.circular(4),
      onTap: () => Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (context) => _EventsDetailsView(advance),
      )),
      child: Container(
        decoration: BoxDecoration(
            border: Border.all(color: Colors.lightBlue.withOpacity(0.5)),
            borderRadius: BorderRadius.circular(4)),
        child: Column(
          children: [
            _rankIcon(advance.rank),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 4.0),
              child: Card(
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
                              elevation: 8,
                              backgroundColor: Colors.teal,
                              label: Row(
                                children: [
                                  Text("${advance.totalPoints}"),
                                  const SizedBox(width: 8),
                                  Image.asset("assets/images/crown.png",
                                      width: 20, color: Colors.white)
                                ],
                              )),
                        ]),
                  )),
            ),
            if (advance.pointsNextRank > advance.pointsCurrentRank)
              _RankBar(advance.pointsCurrentRank, advance.totalPoints,
                  advance.pointsNextRank)
          ],
        ),
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

class _RankBar extends StatelessWidget {
  final int start;
  final int current;
  final int end;

  // assume start <= current <= end
  const _RankBar(this.start, this.current, this.end, {super.key});

  @override
  Widget build(BuildContext context) {
    final advance = (current - start).toDouble() / (end - start);
    return Padding(
      padding: const EdgeInsets.all(8.0),
      child: Row(
        children: [
          Expanded(
            child: ProgressionBar(
                background: Colors.grey,
                layers: [ProgressionLayer(advance, Colors.teal, false)]),
          ),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 4.0),
            child: Text("$current / $end",
                style: Theme.of(context).textTheme.labelSmall),
          )
        ],
      ),
    );
  }
}

class _EventsDetailsView extends StatelessWidget {
  final StudentAdvance advance;

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
      leading: Chip(
          elevation: occurences == 0 ? 0 : 4,
          // shadowColor: occurences == 0 ? null : Colors.white,
          backgroundColor: occurences == 0 ? Colors.grey : Colors.teal,
          label: Text("$occurences")),
      title: Text(eventKLabel(event)),
      titleAlignment: ListTileTitleAlignment.center,
    );
  }
}
