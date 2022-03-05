import 'package:eleve/trivialpoursuit/timeout_bar.dart';
import 'package:flutter/material.dart';

import 'events.gen.dart';
import 'pie.dart';

/// SubmitAnswerNotification is emitted when the player
/// validates his answer
class SubmitAnswerNotification extends Notification {
  final String answer;
  SubmitAnswerNotification(this.answer);
}

class QuestionRoute extends StatefulWidget {
  final ShowQuestion question;
  final Duration timeout;

  const QuestionRoute(this.question, this.timeout, {Key? key})
      : super(key: key);

  @override
  State<QuestionRoute> createState() => _QuestionRouteState();
}

class _QuestionRouteState extends State<QuestionRoute> {
  late TextEditingController _controller;

  bool _enabledValid = false;
  bool _waiting = false;

  @override
  void initState() {
    _controller = TextEditingController();
    _controller.addListener(() {
      setState(() {
        _enabledValid = _controller.text.isNotEmpty;
      });
    });
    super.initState();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  void _onValidated() {
    // send the anwser and wait others players
    SubmitAnswerNotification(_controller.text).dispatch(context);

    setState(() {
      _waiting = true;
      _enabledValid = false; // do no permit answering again
    });
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(top: 40, bottom: 10),
      child: DecoratedBox(
        decoration: BoxDecoration(
            color: widget.question.categorie.color,
            borderRadius: const BorderRadius.all(Radius.circular(10))),
        child: Card(
          elevation: 20,
          child: Padding(
            padding: const EdgeInsets.all(40),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
              children: [
                Text(
                  "Thème ${widget.question.categorie}",
                  style: const TextStyle(fontSize: 20),
                ),
                const Text("Quel est le numéro du thème actuel ?",
                    style: TextStyle(fontSize: 18)),
                TextField(
                  controller: _controller,
                  cursorColor: widget.question.categorie.color,
                  decoration: InputDecoration(
                    focusedBorder: OutlineInputBorder(
                        borderSide: BorderSide(
                      color: widget.question.categorie.color,
                    )),
                    border: const OutlineInputBorder(),
                  ),
                ),
                ElevatedButton(
                  onPressed: _enabledValid ? _onValidated : null,
                  child: const Text(
                    "Valider",
                    style: TextStyle(fontSize: 18),
                  ),
                ),
                TimeoutBar(widget.timeout),
                Text(_waiting ? "En attente des autres joueurs..." : "",
                    style: const TextStyle(fontSize: 16)),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
