import 'package:eleve/trivialpoursuit/events.gen.dart';
import 'package:eleve/trivialpoursuit/pie.dart';
import 'package:eleve/trivialpoursuit/success_recap.dart';
import 'package:flutter/material.dart';

class GameEndPannel extends StatelessWidget {
  final List<int> winners;
  final Map<int, PlayerStatus> players;
  final int ownID;

  const GameEndPannel(this.winners, this.players, this.ownID, {Key? key})
      : super(key: key);

  bool get hasWon => winners.contains(ownID);
  Success get ownSuccess => players[ownID]!.success;

  void _showRecap(BuildContext context) {
    Navigator.of(context).push(
      MaterialPageRoute<void>(
          builder: (context) => SuccessRecapScaffold(players)),
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
          )
        ],
      ),
    );
  }
}
