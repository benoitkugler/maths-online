import 'dart:convert';
import 'dart:math';

import 'package:eleve/build_mode.dart';
import 'package:eleve/questions/question.dart';
import 'package:eleve/questions/types.gen.dart';
import 'package:eleve/quotes.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

class QuestionGallery extends StatefulWidget {
  final BuildMode buildMode;

  const QuestionGallery(this.buildMode, {Key? key}) : super(key: key);

  @override
  State<QuestionGallery> createState() => _QuestionGalleryState();
}

class _QuestionGalleryState extends State<QuestionGallery> {
  List<Question> questions = [];

  int? currentQuestionIndex;
  Color? questionColor;

  @override
  void initState() {
    _loadQuestions();
    super.initState();

    Future.delayed(const Duration(milliseconds: 500), _showSummary);
  }

  void _loadQuestions() async {
    final server = Uri.parse(widget.buildMode.serverURL("/questions"));
    try {
      final resp = await http.get(server);
      setState(() {
        questions = (jsonDecode(resp.body) as List<dynamic>)
            .map(questionFromJson)
            .toList();
      });
    } catch (e) {
      print("ERROR: $e");
    }
  }

  Future<QuestionSyntaxCheckOut> _checkSyntaxCall(
      CheckQuestionSyntaxeNotification v) async {
    final pageIndex = currentQuestionIndex!;
    final uri =
        Uri.parse(widget.buildMode.serverURL("/questions/syntaxe/$pageIndex"));
    final resp = await http.post(uri,
        body: jsonEncode(questionSyntaxCheckInToJson(v.data)),
        headers: {
          'Content-type': 'application/json',
        });
    return questionSyntaxCheckOutFromJson(jsonDecode(resp.body));
  }

  Future<QuestionAnswersOut> _validateCall(ValidQuestionNotification v) async {
    final pageIndex = currentQuestionIndex!;
    final uri =
        Uri.parse(widget.buildMode.serverURL("/questions/answer/$pageIndex"));
    final resp = await http
        .post(uri, body: jsonEncode(questionAnswersInToJson(v.data)), headers: {
      'Content-type': 'application/json',
    });
    return questionAnswersOutFromJson(jsonDecode(resp.body));
  }

  void _validate(ValidQuestionNotification v, BuildContext context) async {
    final rep = await _validateCall(v);
    final crible = rep.results;
    final isValid = crible.values.every((element) => element);
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      backgroundColor: isValid ? Colors.lightGreen : Colors.red,
      content: Text(isValid
          ? "Bonne réponse"
          : "Certains champs sont incorrects ($crible)."),
    ));
  }

  void _showQuestion(int index) {
    setState(() {
      currentQuestionIndex = index;
      questionColor = Color.fromARGB(255, Random().nextInt(256),
          Random().nextInt(256), Random().nextInt(256));
    });
  }

  void _showSummary() {
    Navigator.of(context).push(
      MaterialPageRoute<void>(
          builder: (ct) => Scaffold(
                appBar: AppBar(),
                body: ListView(
                  children: List<Widget>.generate(
                      questions.length,
                      (index) => ListTile(
                            title: Text(
                                "(${index + 1}) " + questions[index].title),
                            onTap: () {
                              _showQuestion(index);
                              Navigator.of(ct).pop();
                            },
                          )).toList(),
                ),
              )),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(actions: [
        TextButton(onPressed: _showSummary, child: const Text("Sommaire"))
      ]),
      body: Padding(
        padding: const EdgeInsets.all(8.0),
        child: currentQuestionIndex != null
            ? QuestionW(
                widget.buildMode,
                questions[currentQuestionIndex!],
                questionColor!,
                (v) => _validate(v, context),
                footerQuote: pickQuote(),
                blockOnSubmit: false,
              )
            : const Text("Chargement"),
      ),
    );
  }
}
