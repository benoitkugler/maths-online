import 'package:eleve/trivialpoursuit/timeout_bar.dart';
import 'package:flutter/material.dart';

import 'events.gen.dart';
import 'pie.dart';

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

  @override
  Widget build(BuildContext context) {
    return DecoratedBox(
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
                "Question du thème ${widget.question.categorie}",
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
                onPressed: _enabledValid
                    ? () {
                        SubmitAnswerNotification(_controller.text)
                            .dispatch(context);
                      }
                    : null,
                child: const Text(
                  "Valider",
                  style: TextStyle(fontSize: 18),
                ),
              ),
              TimeoutBar(widget.timeout),
            ],
          ),
        ),
      ),
    );
  }
}
