import 'dart:convert';

import 'package:eleve/build_mode.dart';
import 'package:eleve/questions/question.dart';
import 'package:eleve/quotes.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_tasks.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

/// [DecrassageAPI] provides the logic to load a list
/// of questions
abstract class DecrassageAPI {
  Future<InstantiatedQuestionsOut> loadQuestions(List<int> ids);
  Future<QuestionAnswersOut> evaluateQuestion(EvaluateQuestionIn answer);
}

/// [ServerDecrassageAPI] is the default implementation of
/// [DecrassageAPI], using an http call to the server.
class ServerDecrassageAPI implements DecrassageAPI {
  final BuildMode buildMode;
  const ServerDecrassageAPI(this.buildMode);

  @override
  Future<InstantiatedQuestionsOut> loadQuestions(List<int> ids) async {
    final uri = buildMode.serverURL("/api/questions/instantiate");
    final resp = await http.post(uri, body: jsonEncode(ids), headers: {
      'Content-type': 'application/json',
    });
    return listInstantiatedQuestionFromJson(jsonDecode(resp.body));
  }

  @override
  Future<QuestionAnswersOut> evaluateQuestion(EvaluateQuestionIn answer) async {
    final uri = buildMode.serverURL("/api/questions/evaluate");
    final resp = await http.post(uri,
        body: jsonEncode(evaluateQuestionInToJson(answer)),
        headers: {
          'Content-type': 'application/json',
        });
    return questionAnswersOutFromJson(jsonDecode(resp.body));
  }
}

class Decrassage extends StatefulWidget {
  final DecrassageAPI api;
  final List<int> idQuestions;

  const Decrassage(this.api, this.idQuestions, {super.key});

  @override
  _DecrassageState createState() => _DecrassageState();
}

class _DecrassageState extends State<Decrassage> {
  InstantiatedQuestionsOut questions = [];

  int questionIndex = 0;
  QuestionController? ct; // null when loading

  @override
  void initState() {
    _loadQuestions();
    super.initState();
  }

  void _loadQuestions() async {
    try {
      final res = await widget.api.loadQuestions(widget.idQuestions);
      setState(() {
        questions = res;
        questionIndex = 0;
        ct = buildController();
      });
    } catch (e) {
      _showError(e);
    }
  }

  void _showError(dynamic error) {
    showError("Une erreur est survenue", error, context);
  }

  // for the current question
  QuestionController buildController() {
    final out =
        QuestionController.fromQuestion(questions[questionIndex].question);
    out.footerQuote = pickQuote();
    return out;
  }

  void _selectQuestion(int newIndex) {
    setState(() {
      questionIndex = newIndex;
      ct = buildController();
    });
  }

  void _evaluateQuestion() async {
    if (ct == null) return;
    final data = ct!.answers();
    setState(() {
      ct!.setFieldsEnabled(false);
      ct!.buttonEnabled = false;
      ct!.buttonLabel = "Correction...";
    });

    final QuestionAnswersOut answerResult;
    final questionOrigin = questions[questionIndex];
    try {
      final args = EvaluateQuestionIn(
          AnswerP(questionOrigin.params, data), questionOrigin.id);
      answerResult = await widget.api.evaluateQuestion(args);
    } catch (e) {
      _showError(e);
      return;
    }
    if (!mounted) return;

    final isValid = answerResult.results.values.every((element) => element);
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      backgroundColor: isValid ? Colors.lightGreen : Colors.red.shade300,
      duration: Duration(seconds: isValid ? 3 : 10),
      content: Text(isValid ? "Bonne réponse" : "Réponse incorrecte"),
      action: isValid
          ? null
          : SnackBarAction(
              label: "Afficher la réponse",
              textColor: Colors.black87,
              onPressed: () {
                ScaffoldMessenger.of(context).hideCurrentSnackBar();
                setState(() {
                  ct?.setAnswers(answerResult.expectedAnswers);
                });
              }),
    ));

    if (isValid) {
      if (questionIndex < questions.length - 1) {
        // go to the next question
        setState(() {
          questionIndex += 1;
          ct = buildController();
        });
      } else {
        // assume the student has followed the order
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(
          backgroundColor: Colors.lightGreen,
          content: Text("Décrassage terminé. Bon travail !"),
        ));
      }
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
                          value: index,
                          child: Text("Question ${index + 1}"),
                        )),
                onSelected: _selectQuestion,
                child: const Text("Choisir la question"),
              ),
            ),
          )
        ],
      ),
      body: Padding(
          padding: const EdgeInsets.all(10),
          child: ct == null
              ? const Center(
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text("Chargement des questions..."),
                      Padding(
                        padding: EdgeInsets.all(12.0),
                        child: CircularProgressIndicator(),
                      ),
                    ],
                  ),
                )
              : QuestionView(
                  questions[questionIndex].question,
                  ct!,
                  _evaluateQuestion,
                  Colors.yellow,
                  title: "Question ${questionIndex + 1}",
                )),
    );
  }
}
