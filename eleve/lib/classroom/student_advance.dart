import 'dart:math';

import 'package:eleve/shared/progression_bar.dart';
import 'package:eleve/types/src_sql_events.dart';
import 'package:flutter/material.dart';

class AdvanceSummary extends StatelessWidget {
  final StudentAdvance advance;

  const AdvanceSummary(this.advance, {super.key});

  @override
  Widget build(BuildContext context) {
    final color = _rankColors[advance.rank];
    return InkWell(
      borderRadius: BorderRadius.circular(4),
      onTap: () => Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (context) => _EventsDetailsView(advance),
      )),
      child: Container(
        padding: const EdgeInsets.all(8),
        decoration: BoxDecoration(
            border: Border.all(color: Colors.lightBlue.withValues(alpha: 0.5)),
            borderRadius: BorderRadius.circular(4)),
        child: Column(
          children: [
            rankIcon(advance.rank),
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 4.0),
              child: Container(
                  decoration: BoxDecoration(
                      boxShadow: [
                        BoxShadow(color: color, spreadRadius: 2, blurRadius: 2)
                      ],
                      color: Colors.white.withValues(alpha: 0.2),
                      borderRadius: BorderRadius.circular(4),
                      border: Border.all(color: color)),
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
                      ])),
            ),
            if (advance.pointsNextRank > advance.pointsCurrentRank)
              Padding(
                padding: const EdgeInsets.symmetric(vertical: 4.0),
                child: _RankBar(advance.pointsCurrentRank, advance.totalPoints,
                    advance.pointsNextRank),
              ),
            if (advance.flames > 0)
              Padding(
                padding: const EdgeInsets.symmetric(vertical: 4.0),
                child: FlamesBar(advance.flames),
              )
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
  Color.fromARGB(255, 176, 190, 190),
  Color.fromARGB(255, 216, 197, 23),
  Color.fromARGB(255, 100, 216, 210)
];

Image rankIcon(int rank, {double width = 80}) {
  rank = rank.clamp(0, _rankColors.length - 1);
  return Image.asset("assets/images/ranks/rank-$rank.png", width: width);
}

class _RankBar extends StatelessWidget {
  final int start;
  final int current;
  final int end;

  // assume start <= current <= end
  const _RankBar(this.start, this.current, this.end);

  @override
  Widget build(BuildContext context) {
    final advance = (current - start).toDouble() / (end - start);
    return Row(
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
    );
  }
}

class _EventsDetailsView extends StatelessWidget {
  final StudentAdvance advance;

  const _EventsDetailsView(this.advance);

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
  const _EventOccurences(this.event, this.occurences);

  @override
  Widget build(BuildContext context) {
    return ListTile(
      leading: Chip(
          elevation: occurences == 0 ? 0 : 4,
          backgroundColor: occurences == 0 ? Colors.grey : Colors.teal,
          label: Text("$occurences")),
      title: Text(eventKLabel(event)),
      titleAlignment: ListTileTitleAlignment.center,
    );
  }
}

class FlamesBar extends StatelessWidget {
  final int flames;
  const FlamesBar(this.flames, {super.key});

  @override
  Widget build(BuildContext context) {
    final color = Colors.orangeAccent.shade700;
    final iconsNumber = min(flames, 10);
    return Container(
      decoration: BoxDecoration(
          boxShadow: [BoxShadow(color: color, blurRadius: 2, spreadRadius: 2)],
          borderRadius: BorderRadius.circular(4),
          color: Colors.white),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: [
          Row(
              mainAxisSize: MainAxisSize.min,
              children: List.filled(
                  iconsNumber,
                  Image.asset("assets/images/fire.png",
                      width: 20, color: color))),
          Chip(
            backgroundColor: color,
            elevation: 4,
            visualDensity: VisualDensity.compact,
            label: Text("$flames"),
          )
        ],
      ),
    );
  }
}
