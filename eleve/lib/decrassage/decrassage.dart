import 'dart:convert';

import 'package:eleve/build_mode.dart';
import 'package:eleve/exercices/question.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:eleve/quotes.dart';
import 'package:eleve/shared_gen.dart' as shared;
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

class Decrassage extends StatefulWidget {
  final List<int> idQuestions;
  final BuildMode buildMode;

  const Decrassage(this.idQuestions, this.buildMode, {Key? key})
      : super(key: key);

  @override
  _DecrassageState createState() => _DecrassageState();
}

class _DecrassageState extends State<Decrassage> {
  shared.InstantiateQuestionsOut questions = [];
  int? currentQuestionIndex;
  Map<int, Answer>? currentAnswer;

  @override
  void initState() {
    _loadQuestions();
    super.initState();
  }

  shared.InstantiatedQuestion? get currentQuestion =>
      currentQuestionIndex == null ? null : questions[currentQuestionIndex!];

  void _loadQuestions() async {
    try {
      final uri =
          Uri.parse(widget.buildMode.serverURL("/api/questions/instantiate"));
      final resp =
          await http.post(uri, body: jsonEncode(widget.idQuestions), headers: {
        'Content-type': 'application/json',
      });
      setState(() {
        questions =
            shared.listInstantiatedQuestionFromJson(jsonDecode(resp.body));
        currentQuestionIndex = 0;
        currentAnswer = null;
      });
    } catch (e) {
      _showError(e);
    }
  }

  void _showError(dynamic error) {
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      duration: const Duration(seconds: 5),
      backgroundColor: Theme.of(context).colorScheme.error,
      content: Text("Une erreur est survenue : $error"),
    ));
  }

  void _selectQuestion(int questionIndex) {
    setState(() {
      currentQuestionIndex = questionIndex;
      currentAnswer = null;
    });
  }

  void _evaluateQuestion(ValidQuestionNotification data) async {
    try {
      final uri =
          Uri.parse(widget.buildMode.serverURL("/api/questions/evaluate"));
      final args = shared.EvaluateQuestionIn(
          data.data, currentQuestion!.params, currentQuestion!.id);
      final resp = await http.post(uri,
          body: jsonEncode(shared.evaluateQuestionInToJson(args)),
          headers: {
            'Content-type': 'application/json',
          });
      final answerResult = questionAnswersOutFromJson(jsonDecode(resp.body));
      final isValid = answerResult.results.values.every((element) => element);
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        backgroundColor: isValid ? Colors.lightGreen : Colors.red.shade200,
        duration: Duration(seconds: isValid ? 2 : 4),
        content: Text(isValid ? "Bonne réponse" : "Réponse incorrecte"),
        action: isValid
            ? null
            : SnackBarAction(
                label: "Afficher la réponse",
                onPressed: () => setState(() {
                  currentAnswer = answerResult.expectedAnswers;
                }),
              ),
      ));
      if (isValid) {
        if (currentQuestionIndex! < questions.length - 1) {
          // go to the next question
          setState(() {
            currentQuestionIndex = currentQuestionIndex! + 1;
            currentAnswer = null;
          });
        } else {
          // assume the student has followed the order
          ScaffoldMessenger.of(context).showSnackBar(const SnackBar(
            backgroundColor: Colors.lightGreen,
            content: Text("Décrassage terminé. Bon travail !"),
          ));
        }
      }
    } catch (e) {
      _showError(e);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        automaticallyImplyLeading: true,
        title: const Text("Décrassage"),
        actions: [
          Center(
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 8.0),
              child: PopupMenuButton(
                itemBuilder: (context) => List<PopupMenuEntry<int>>.generate(
                    questions.length,
                    (index) => PopupMenuItem(
                          child: Text("Question ${index + 1}"),
                          value: index,
                        )),
                child: const Text("Choisir la question"),
                onSelected: _selectQuestion,
              ),
            ),
          )
        ],
      ),
      body: Padding(
          padding: const EdgeInsets.all(10),
          child: questions.isEmpty
              ? Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: const [
                      Text("Chargement"),
                      Padding(
                        padding: EdgeInsets.all(12.0),
                        child: CircularProgressIndicator(),
                      ),
                    ],
                  ),
                )
              : QuestionW(
                  widget.buildMode,
                  currentQuestion!.question,
                  Colors.pink,
                  _evaluateQuestion,
                  answer: currentAnswer,
                  blockOnSubmit: false,
                  footerQuote: pickQuote(),
                  timeout: null,
                  title: "Question ${currentQuestionIndex! + 1}",
                )),
    );
  }
}
