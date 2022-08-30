import 'package:eleve/activities/trivialpoursuit/events.gen.dart';
import 'package:eleve/activities/trivialpoursuit/pie.dart';
import 'package:flutter/material.dart';

List<PlayerID> sortedPlayers(Map<PlayerID, PlayerStatus> successes) {
  final list = successes.keys.toList();
  list.sort();
  return list;
}

class SuccessRecapScaffold extends StatelessWidget {
  final Map<PlayerID, PlayerStatus> successes;

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
  final PlayerID playerID;
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
  final PlayerID ownPlayer;
  final Map<PlayerID, PlayerStatus> successes;

  const SuccessRecapRow(this.ownPlayer, this.successes, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ListView(
      shrinkWrap: true,
      scrollDirection: Axis.horizontal,
      reverse: true,
      children: sortedPlayers(successes)
          .where((element) => element != ownPlayer)
          .map((e) {
        final isInactive = successes[e]!.isInactive;
        return Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Transform.scale(
                scale: 0.8,
                child: Pie(
                  2,
                  successes[e]!.success,
                  backgroundColor: isInactive ? Colors.grey : Colors.white,
                )),
            Padding(
              padding: const EdgeInsets.only(bottom: 5),
              child: SizedBox(
                width: 80,
                child: Text(
                  successes[e]!.name,
                  style: TextStyle(
                      fontSize: 12,
                      color: isInactive ? Colors.grey.shade400 : null),
                  textAlign: TextAlign.center,
                ),
              ),
            )
          ],
        );
      }).toList(),
    );
  }
}
