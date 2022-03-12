import 'package:eleve/trivialpoursuit/events.gen.dart';
import 'package:eleve/trivialpoursuit/pie.dart';
import 'package:flutter/material.dart';

List<int> sortedPlayers(Map<int, PlayerStatus> successes) {
  final list = successes.keys.toList();
  list.sort();
  return list;
}

class SuccessRecapScaffold extends StatelessWidget {
  final Map<int, PlayerStatus> successes;

  const SuccessRecapScaffold(this.successes, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text("Succès obtenus")),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 30, vertical: 10),
          child: ListView(
              shrinkWrap: true,
              children: sortedPlayers(successes)
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

class SuccessRecapRow extends StatelessWidget {
  final int ownPlayer;
  final Map<int, PlayerStatus> successes;

  const SuccessRecapRow(this.ownPlayer, this.successes, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: sortedPlayers(successes)
          .where((element) => element != ownPlayer)
          .map((e) => Column(
                children: [
                  Transform.scale(
                      scale: 0.8, child: Pie(2, successes[e]!.success)),
                  Padding(
                    padding: const EdgeInsets.only(bottom: 5),
                    child: Text(successes[e]!.name),
                  )
                ],
              ))
          .toList(),
    );
  }
}
