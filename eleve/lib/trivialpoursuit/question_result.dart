import 'package:eleve/trivialpoursuit/events.gen.dart';
import 'package:flutter/material.dart';

import 'categories.dart';

class WantNextTurnNotification extends Notification {
  final WantNextTurn event;
  WantNextTurnNotification(this.event);
}

/// [QuestionResult] shows the result of the last question
class QuestionResult extends StatefulWidget {
  final PlayerAnswerResult event;

  const QuestionResult(this.event, {Key? key}) : super(key: key);

  @override
  State<QuestionResult> createState() => _QuestionResultState();
}

class _QuestionResultState extends State<QuestionResult> {
  bool markQuestion = false;

  @override
  Widget build(BuildContext context) {
    final backgroundColor =
        widget.event.success ? Colors.lightGreen.shade400 : Colors.red;
    final content = Row(mainAxisSize: MainAxisSize.min, children: [
      Padding(
        padding: const EdgeInsets.only(right: 5),
        child: Icon(widget.event.success
            ? const IconData(0xe156, fontFamily: 'MaterialIcons')
            : const IconData(0xe868, fontFamily: 'MaterialIcons')),
      ),
      Text(widget.event.success
          ? "Bonne réponse, bravo !"
          : "Réponse incorrecte, dommage...")
    ]);

    return Scaffold(
      body: Padding(
        padding: const EdgeInsets.all(8.0),
        child: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
            children: [
              const Text(
                "Résultat",
                style: TextStyle(fontSize: 22),
              ),
              Card(
                  color: backgroundColor,
                  child: Padding(
                    padding: const EdgeInsets.all(12.0),
                    child: content,
                  )),
              if (!widget.event.success)
                _ExpectedAnswer(
                    widget.event.correctAnwser, widget.event.categorie.color),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  InkWell(
                    onTap: () => setState(() {
                      markQuestion = !markQuestion;
                    }),
                    borderRadius: BorderRadius.circular(5),
                    child: Padding(
                      padding: const EdgeInsets.only(right: 8.0),
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
                      onPressed: () =>
                          WantNextTurnNotification(WantNextTurn(markQuestion))
                              .dispatch(context),
                      child:
                          Row(mainAxisSize: MainAxisSize.min, children: const [
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

class _ExpectedAnswer extends StatelessWidget {
  final String content;
  final Color color;

  const _ExpectedAnswer(this.content, this.color, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        border: Border.all(color: color),
        borderRadius: BorderRadius.circular(6),
      ),
      padding: const EdgeInsets.all(12),
      child: Text("La bonne réponse est : $content"),
    );
  }
}
