import 'package:eleve/trivialpoursuit/events.gen.dart';
import 'package:eleve/trivialpoursuit/pie.dart';
import 'package:flutter/material.dart';

class SuccessRecap extends StatelessWidget {
  final Map<int, PlayerStatus> successes;

  const SuccessRecap(this.successes, {Key? key}) : super(key: key);

  List<int> get sortedPlayers {
    final list = successes.keys.toList();
    list.sort();
    return list;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text("Succès obtenus")),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 30, vertical: 10),
          child: ListView(
              shrinkWrap: true,
              children: sortedPlayers
                  .map((e) => _PlayerTile(e, successes[e]!))
                  .toList()),
        ),
      ),
    );
  }
}

class _PlayerTile extends StatelessWidget {
  final int playerID;
  final PlayerStatus status;

  const _PlayerTile(this.playerID, this.status, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ListTile(
      title: Text(status.name),
      trailing: Hero(
        tag: "recap_$playerID",
        child: Pie(2, status.success),
      ),
      contentPadding: const EdgeInsets.symmetric(horizontal: 8, vertical: 15),
    );
  }
}
