import 'package:eleve/quotes.dart';
import 'package:eleve/shared/title.dart';
import 'package:eleve/types/src_trivial.dart';
import 'package:flutter/material.dart';

import 'categories.dart';

class WantNextTurnNotification extends Notification {
  final WantNextTurn event;
  WantNextTurnNotification(this.event);
}

/// [QuestionResult] shows the result of the last question
class QuestionResult extends StatefulWidget {
  final PlayerID playerID;
  final PlayerAnswerResults event;
  final Map<PlayerID, PlayerStatus> players;
  final void Function() showLastQuestion;

  const QuestionResult(
      this.playerID, this.event, this.players, this.showLastQuestion,
      {Key? key})
      : super(key: key);

  @override
  State<QuestionResult> createState() => _QuestionResultState();
}

class _QuestionResultState extends State<QuestionResult> {
  bool markQuestion = false;

  PlayerAnswerResult get ownResult => widget.event.results[widget.playerID]!;

  List<String> get playersCorrect => widget.players.entries
      .where((entry) => widget.event.results[entry.key]!.success)
      .map((e) => e.value.name)
      .toList();

  List<String> get playersIncorrect => widget.players.entries
      .where((entry) => !widget.event.results[entry.key]!.success)
      .map((e) => e.value.name)
      .toList();

  void _onContinue() {
    WantNextTurnNotification(WantNextTurn(markQuestion)).dispatch(context);
  }

  @override
  Widget build(BuildContext context) {
    final backgroundColor =
        ownResult.success ? Colors.lightGreen.shade400 : Colors.red;
    final content = Row(mainAxisSize: MainAxisSize.min, children: [
      Padding(
        padding: const EdgeInsets.only(right: 5),
        child: Icon(ownResult.success
            ? const IconData(0xe156, fontFamily: 'MaterialIcons')
            : const IconData(0xe868, fontFamily: 'MaterialIcons')),
      ),
      Text(ownResult.success
          ? "Bonne réponse, bravo !"
          : "Réponse incorrecte, dommage...")
    ]);

    final correct =
        playersCorrect.isEmpty ? "" : " à ${playersCorrect.join(", ")} ou";
    final hintSuccess = playersIncorrect.isEmpty
        ? ""
        : "Peut-être peux-tu expliquer à ${playersIncorrect.join(", ")} cette question ?";
    final hint = ownResult.success
        ? hintSuccess
        : "N'hésite pas à demander de l'aide$correct au prof. avant de continuer !";

    return Scaffold(
      appBar: AppBar(automaticallyImplyLeading: false, actions: [
        TextButton(
            onPressed: widget.showLastQuestion,
            child: const Text("Afficher la question"))
      ]),
      body: Padding(
        padding: const EdgeInsets.all(8.0),
        child: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
            children: [
              ColoredTitle("Résultat", widget.event.categorie.color),
              Card(
                color: backgroundColor,
                child: Padding(
                  padding: const EdgeInsets.all(12.0),
                  child: content,
                ),
              ),
              Padding(
                padding: const EdgeInsets.all(8.0),
                child: Text(
                  hint,
                  style: const TextStyle(fontSize: 18),
                ),
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  if (ownResult.askForMask)
                    InkWell(
                      onTap: () => setState(() {
                        markQuestion = !markQuestion;
                      }),
                      borderRadius: BorderRadius.circular(5),
                      child: Padding(
                        padding: const EdgeInsets.only(right: 20.0),
                        child: Row(
                          children: [
                            Checkbox(
                                value: markQuestion,
                                onChanged: (checked) => setState(() {
                                      markQuestion = checked ?? false;
                                    })),
                            const Text("Retenir cette question"),
                          ],
                        ),
                      ),
                    ),
                  ElevatedButton(
                      onPressed: _onContinue,
                      child:
                          const Row(mainAxisSize: MainAxisSize.min, children: [
                        Text("Continuer"),
                        Padding(
                          padding:
                              EdgeInsets.only(left: 4.0, top: 5, bottom: 5),
                          child: Icon(
                              IconData(0xf05bd, fontFamily: 'MaterialIcons')),
                        ),
                      ])),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class WaitForPlayersDialog extends StatelessWidget {
  final List<String> playerNames;

  const WaitForPlayersDialog(this.playerNames, {super.key});

  @override
  Widget build(BuildContext context) {
    return Dialog(
      child: Card(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Padding(
              padding: const EdgeInsets.all(10.0),
              child: Text(
                  "En attente de ${playerNames.join(', ')} pour le prochain tour"),
            ),
            const LinearProgressIndicator(),
            const SizedBox(height: 40),
            Quote(pickQuote()),
          ],
        ),
      ),
    );
  }
}
