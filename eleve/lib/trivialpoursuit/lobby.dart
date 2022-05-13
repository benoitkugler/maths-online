import 'package:eleve/quotes.dart';
import 'package:flutter/material.dart';

class GameLobby extends StatelessWidget {
  final Map<int, String> players;
  final int player;

  const GameLobby(this.players, this.player, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final sorted = players.keys.toList();
    sorted.sort();
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: [
          const Text(
            "En attente d'autres joueurs...",
            style: TextStyle(fontSize: 20),
          ),
          const CircularProgressIndicator(),
          Quote(pickQuote()),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 8.0),
            child: Wrap(
              spacing: 20,
              runSpacing: 15,
              alignment: WrapAlignment.spaceEvenly,
              children: sorted
                  .map((e) => _PlayerAvatar(players[e]!, e == player))
                  .toList(),
            ),
          ),
        ],
      ),
    );
  }
}

class _PlayerAvatar extends StatelessWidget {
  final String name;
  final bool isCurrentPlayer;
  const _PlayerAvatar(this.name, this.isCurrentPlayer, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return DecoratedBox(
      decoration: BoxDecoration(
          boxShadow: [
            BoxShadow(
                color: isCurrentPlayer ? Colors.yellow : Colors.white,
                blurRadius: 5)
          ],
          color:
              isCurrentPlayer ? Colors.yellow : Colors.white.withOpacity(0.8),
          borderRadius: const BorderRadius.all(Radius.circular(5))
          // boxShadow:
          ),
      child: Card(
        child: Padding(
          padding: const EdgeInsets.all(8.0),
          child: Text(name),
        ),
      ),
    );
  }
}
