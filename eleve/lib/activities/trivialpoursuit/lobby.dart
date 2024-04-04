import 'package:eleve/classroom/student_advance.dart';
import 'package:eleve/quotes.dart';
import 'package:eleve/types/src_trivial.dart';
import 'package:flutter/material.dart';

class GameLobby extends StatelessWidget {
  final Map<PlayerID, String> players;
  final Map<PlayerID, int> playerRanks;
  final PlayerID player;

  /// if [onStart] is not null, shows "start game" button
  final void Function()? onStart;

  const GameLobby(this.players, this.playerRanks, this.player, this.onStart,
      {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    final sorted = players.keys.toList();
    sorted.sort();
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          const Padding(
            padding: EdgeInsets.symmetric(vertical: 20),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                CircularProgressIndicator(),
                SizedBox(width: 16),
                Text(
                  "En attente d'autres joueurs...",
                  style: TextStyle(fontSize: 20),
                ),
              ],
            ),
          ),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 8.0),
            child: Wrap(
              spacing: 20,
              runSpacing: 15,
              alignment: WrapAlignment.spaceEvenly,
              children: sorted
                  .map((e) =>
                      _PlayerAvatar(players[e]!, playerRanks[e]!, e == player))
                  .toList(),
            ),
          ),
          if (onStart != null)
            ElevatedButton(
                onPressed: onStart,
                child: const Text("Lancer la partie !",
                    style: TextStyle(fontSize: 18))),
          Quote(pickQuote()),
        ],
      ),
    );
  }
}

class _PlayerAvatar extends StatelessWidget {
  final String name;
  final int rank;
  final bool isCurrentPlayer;
  const _PlayerAvatar(this.name, this.rank, this.isCurrentPlayer, {Key? key})
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
          borderRadius: const BorderRadius.all(Radius.circular(5))),
      child: Card(
        child: Padding(
          padding: const EdgeInsets.all(8.0),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              rankIcon(6, width: 42),
              const SizedBox(width: 4),
              Text(name),
            ],
          ),
        ),
      ),
    );
  }
}
