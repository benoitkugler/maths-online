import 'dart:convert';

import 'package:eleve/build_mode.dart';
import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/question.dart';
import 'package:eleve/quotes.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_tasks.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

/// [DecrassageAPI] provides the logic to load a list
/// of questions
abstract class DecrassageAPI extends FieldAPI {
  Future<InstantiateQuestionsOut> loadQuestions(List<int> ids);
  Future<QuestionAnswersOut> evaluateQuestion(EvaluateQuestionIn answer);
}

/// [ServerDecrassageAPI] is the default implementation of
/// [DecrassageAPI], using an http call to the server.
class ServerDecrassageAPI extends ServerFieldAPI implements DecrassageAPI {
  const ServerDecrassageAPI(super.buildMode);

  @override
  Future<InstantiateQuestionsOut> loadQuestions(List<int> ids) async {
    final uri = Uri.parse(buildMode.serverURL("/api/questions/instantiate"));
    final resp = await http.post(uri, body: jsonEncode(ids), headers: {
      'Content-type': 'application/json',
    });
    return listInstantiatedQuestionFromJson(jsonDecode(resp.body));
  }

  @override
  Future<QuestionAnswersOut> evaluateQuestion(EvaluateQuestionIn answer) async {
    final uri = Uri.parse(buildMode.serverURL("/api/questions/evaluate"));
    final resp = await http.post(uri,
        body: jsonEncode(evaluateQuestionInToJson(answer)),
        headers: {
          'Content-type': 'application/json',
        });
    return questionAnswersOutFromJson(jsonDecode(resp.body));
  }
}

class DecrassageQuestionController extends BaseQuestionController {
  final int questionIndex;
  final void Function(QuestionAnswersIn) onValid;

  DecrassageQuestionController(
      this.questionIndex, Question question, FieldAPI api, this.onValid)
      : super(question, api) {
    state.footerQuote = pickQuote();
  }

  @override
  void onPrimaryButtonClick() {
    state.buttonEnabled = false;
    state.buttonLabel = "Correction...";
    onValid(answers());
  }
}

class Decrassage extends StatefulWidget {
  final DecrassageAPI api;
  final List<int> idQuestions;

  const Decrassage(this.api, this.idQuestions, {Key? key}) : super(key: key);

  @override
  _DecrassageState createState() => _DecrassageState();
}

class _DecrassageState extends State<Decrassage> {
  InstantiateQuestionsOut questions = [];

  DecrassageQuestionController? ct; // null when loading

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
        ct = controllerForIndex(0);
      });
    } catch (e) {
      _showError(e);
    }
  }

  void _showError(dynamic error) {
    showError("Une erreur est survenue", error, context);
  }

  DecrassageQuestionController controllerForIndex(int questionIndex) {
    return DecrassageQuestionController(questionIndex,
        questions[questionIndex].question, widget.api, _evaluateQuestion);
  }

  void _selectQuestion(int questionIndex) {
    setState(() {
      ct = controllerForIndex(questionIndex);
    });
  }

  void _evaluateQuestion(QuestionAnswersIn data) async {
    final qc = ct!;
    final QuestionAnswersOut answerResult;
    final questionOrigin = questions[qc.questionIndex];
    try {
      final args = EvaluateQuestionIn(
          AnswerP(questionOrigin.params, data), questionOrigin.id);
      answerResult = await widget.api.evaluateQuestion(args);
    } catch (e) {
      _showError(e);
      return;
    }

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
                  qc.setAnswers(answerResult.expectedAnswers);
                });
              }),
    ));

    if (isValid) {
      if (qc.questionIndex < questions.length - 1) {
        // go to the next question
        setState(() {
          ct = controllerForIndex(ct!.questionIndex + 1);
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
              ? Center(
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: const [
                      Text("Chargement des questions..."),
                      Padding(
                        padding: EdgeInsets.all(12.0),
                        child: CircularProgressIndicator(),
                      ),
                    ],
                  ),
                )
              : QuestionW(
                  ct!,
                  Colors.pink,
                  title: "Question ${ct!.questionIndex + 1}",
                )),
    );
  }
}
