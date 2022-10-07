import 'dart:convert';

import 'package:eleve/build_mode.dart';
import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/question.dart';
import 'package:eleve/questions/types.gen.dart';
import 'package:eleve/quotes.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/shared_gen.dart' as shared;
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

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
  final List<int> idQuestions;
  final BuildMode buildMode;

  const Decrassage(this.idQuestions, this.buildMode, {Key? key})
      : super(key: key);

  @override
  _DecrassageState createState() => _DecrassageState();
}

class _DecrassageState extends State<Decrassage> {
  shared.InstantiateQuestionsOut questions = [];
  // int? currentQuestionIndex;
  // Map<int, Answer>? currentAnswer;

  DecrassageQuestionController? ct; // null when loading

  @override
  void initState() {
    _loadQuestions();
    super.initState();
  }

  // shared.InstantiatedQuestion? get currentQuestion =>
  //     currentQuestionIndex == null ? null : questions[currentQuestionIndex!];

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
    return DecrassageQuestionController(
        questionIndex,
        questions[questionIndex].question,
        ServerFieldAPI(widget.buildMode),
        _evaluateQuestion);
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
      final uri =
          Uri.parse(widget.buildMode.serverURL("/api/questions/evaluate"));
      final args = shared.EvaluateQuestionIn(
          shared.Answer(questionOrigin.params, data), questionOrigin.id);
      final resp = await http.post(uri,
          body: jsonEncode(shared.evaluateQuestionInToJson(args)),
          headers: {
            'Content-type': 'application/json',
          });
      answerResult = questionAnswersOutFromJson(jsonDecode(resp.body));
    } catch (e) {
      _showError(e);
      return;
    }

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
                qc.setAnswers(answerResult.expectedAnswers);
              }),
            ),
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
          child: ct == null
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
                  ct!,
                  Colors.pink,
                  title: "Question ${ct!.questionIndex + 1}",
                )),
    );
  }
}
