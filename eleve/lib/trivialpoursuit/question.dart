import 'package:eleve/trivialpoursuit/events.gen.dart';
import 'package:flutter/widgets.dart';

// TODO: implements this widget
class QuestionRoute extends StatelessWidget {
  final ShowQuestion question;

  const QuestionRoute(this.question, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(20),
      child: Text("Question : ${question.question}"),
    );
  }
}
