import 'dart:convert';
import 'dart:math';

import 'package:eleve/exercices/question.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:eleve/trivialpoursuit/events.gen.dart' as game;
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

final server = Uri.parse("http://localhost:3030/");

class QuestionGallery extends StatefulWidget {
  const QuestionGallery({Key? key}) : super(key: key);

  @override
  State<QuestionGallery> createState() => _QuestionGalleryState();
}

class _QuestionGalleryState extends State<QuestionGallery> {
  @override
  void initState() {
    _loadQuestions();
    super.initState();
  }

  void _loadQuestions() async {
    try {
      final resp = await http.get(server);
      setState(() {
        questions = jsonDecode(resp.body) as List<dynamic>;
      });
    } catch (e) {
      print("ERROR: $e");
    }
  }

  Future<QuestionSyntaxCheckOut> _checkSyntaxCall(
      CheckQuestionSyntaxeNotification v) async {
    final pageIndex = _controller.page!.toInt();
    final uri = server.replace(path: "syntaxe/$pageIndex");
    final resp = await http.post(uri,
        body: jsonEncode({"ID": v.id, "Answer": answerToJson(v.answer)}),
        headers: {
          'Content-type': 'application/json',
        });
    return questionSyntaxCheckOutFromJson(jsonDecode(resp.body));
  }

  Future<QuestionAnswersOut> _validateCall(ValidQuestionNotification v) async {
    final pageIndex = _controller.page!.toInt();
    final uri = server.replace(path: "answer/$pageIndex");
    final resp = await http
        .post(uri, body: jsonEncode(questionAnswersInToJson(v.data)), headers: {
      'Content-type': 'application/json',
    });
    return questionAnswersOutFromJson(jsonDecode(resp.body));
  }

  List<dynamic> questions = [];
  final _controller = PageController(initialPage: 0);

  void _checkSyntax(
      CheckQuestionSyntaxeNotification v, BuildContext context) async {
    final rep = await _checkSyntaxCall(v);
    if (rep.isValid) {
      print("OK !");
      return;
    }
    final reason = rep.reason;
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      backgroundColor: Colors.red,
      content: Text.rich(TextSpan(children: [
        const TextSpan(text: "Syntaxe invalide: "),
        TextSpan(
            text: reason, style: const TextStyle(fontWeight: FontWeight.bold)),
      ])),
    ));
  }

  void _validate(ValidQuestionNotification v, BuildContext context) async {
    final rep = await _validateCall(v);
    final crible = rep.data;
    final isValid = crible.values.every((element) => element);
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      backgroundColor: isValid ? Colors.lightGreen : Colors.red,
      content: Text(isValid ? "Bonne réponse" : "Réponse incomplète $crible"),
    ));
  }

  NotificationListener _fromJSON(dynamic json, BuildContext context) {
    final question = questionFromJson(json);
    return NotificationListener<CheckQuestionSyntaxeNotification>(
      onNotification: (v) {
        _checkSyntax(v, context);
        return true;
      },
      child: NotificationListener<ValidQuestionNotification>(
        onNotification: (v) {
          _validate(v, context);
          return true;
        },
        child: QuestionPage(
            question,
            game.Categorie
                .values[Random().nextInt(game.Categorie.values.length)]),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return PageView(
      controller: _controller,
      children: questions.map((q) => _fromJSON(q, context)).toList(),
    );
  }
}
