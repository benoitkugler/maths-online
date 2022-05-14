import 'package:eleve/build_mode.dart';
import 'package:eleve/decrassage/decrassage.dart';
import 'package:eleve/trivialpoursuit/events.gen.dart';
import 'package:eleve/trivialpoursuit/pie.dart';
import 'package:eleve/trivialpoursuit/success_recap.dart';
import 'package:flutter/material.dart';

class GameEndPannel extends StatelessWidget {
  final BuildMode buildMode;
  final GameEnd data;
  final Map<int, PlayerStatus> players;
  final int ownID;

  const GameEndPannel(this.buildMode, this.data, this.players, this.ownID,
      {Key? key})
      : super(key: key);

  List<int> get winners => data.winners;
  bool get hasWon => winners.contains(ownID);
  Success get ownSuccess => players[ownID]!.success;

  /// may be empty if the teacher disabled decrassage
  List<int> get decrassage => data.questionDecrassageIds[ownID] ?? [];

  void _showRecap(BuildContext context) {
    Navigator.of(context).push(
      MaterialPageRoute<void>(
          builder: (context) => SuccessRecapScaffold(players)),
    );
  }

  void _showDecrassage(BuildContext context) {
    Navigator.of(context).push(
      MaterialPageRoute<void>(
          builder: (context) => Decrassage(decrassage, buildMode)),
    );
  }

  @override
  Widget build(BuildContext context) {
    final List<Widget> congrats = [];
    if (hasWon) {
      congrats.add(const Text("Vous avez gagné, bravo !",
          style: TextStyle(
            color: Colors.yellow,
            fontSize: 25,
          )));
    }
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: [
          const Text(
            "Partie terminée",
            style: TextStyle(fontSize: 25),
          ),
          Pie.asButton(() => _showRecap(context), 2, ownSuccess),
          ...congrats,
          Column(
            children: [
              Padding(
                padding: const EdgeInsets.only(bottom: 20),
                child: Text(
                  winners.length == 1
                      ? "Le gagnant est :"
                      : "Les gagnants sont :",
                  style: const TextStyle(fontSize: 20),
                ),
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                children: winners.map((e) {
                  final name = players[e]!.name;
                  return DecoratedBox(
                      decoration: const BoxDecoration(boxShadow: [
                        BoxShadow(
                          color: Colors.yellow,
                          blurRadius: 5,
                        )
                      ], borderRadius: BorderRadius.all(Radius.circular(10))),
                      child: Card(
                        child: Padding(
                          padding: const EdgeInsets.all(12.0),
                          child: Text(name),
                        ),
                      ));
                }).toList(),
              ),
            ],
          ),
          if (decrassage.isNotEmpty)
            ElevatedButton(
                onPressed: () => _showDecrassage(context),
                child: const Text("Continuer vers le décrassage"))
        ],
      ),
    );
  }
}
