import 'dart:convert';
import 'dart:math';

import 'package:eleve/exercices/question.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:eleve/trivialpoursuit/events.gen.dart';
import 'package:flutter/material.dart';

class QuestionGallery extends StatelessWidget {
  QuestionGallery({Key? key}) : super(key: key);
  final _controller = PageController(initialPage: 1);

  static const _questions = [
    """
{"Title":"Calcul littéral","Content":[{"Data":{"Text":"Développer l’expression : "},"Kind":4},{"Data":{"Content":"\\\\left(x - 6\\\\right)\\\\left(4x - 3\\\\right)","IsInline":true},"Kind":0}]}
""",
    """
{"Title":"Calcul littéral","Content":[{"Data":{"Text":"Écrire sous une seule fraction : "},"Kind":4},{"Data":{"Content":"\\\\frac{1}{3} + \\\\frac{2}{5}","IsInline":true},"Kind":0},{"Data":{"Content":"\\\\frac{1}{3} + \\\\frac{2}{5} = \\\\frac{a}{b}","IsInline":false},"Kind":0},{"Data":{"Text":"Déterminer "},"Kind":4},{"Data":{"Content":"a = ","IsInline":true},"Kind":0},{"Data":{},"Kind":3},{"Data":{"Text":" et "},"Kind":4},{"Data":{"Content":"b = ","IsInline":true},"Kind":0},{"Data":{},"Kind":3}]}
""",
    """
{"Title":"Très longue question horizontale","Content":[{"Data":{"Text":"Écrire sous une seule fraction : "},"Kind":4},{"Data":{"Content":"\\\\frac{1}{3} + \\\\frac{2}{5}","IsInline":true},"Kind":0},{"Data":{"Text":"Écrire sous une seule fraction : "},"Kind":4},{"Data":{"Content":"\\\\frac{1}{3} + \\\\frac{2}{5}","IsInline":true},"Kind":0},{"Data":{"Text":"Écrire sous une seule fraction : "},"Kind":4},{"Data":{"Content":"\\\\frac{1}{3} + \\\\frac{2}{5}","IsInline":true},"Kind":0},{"Data":{"Text":"Écrire sous une seule fraction : "},"Kind":4},{"Data":{"Content":"\\\\frac{1}{3} + \\\\frac{2}{5}","IsInline":true},"Kind":0},{"Data":{"Text":"Écrire sous une seule fraction : "},"Kind":4},{"Data":{"Content":"\\\\frac{1}{3} + \\\\frac{2}{5}","IsInline":true},"Kind":0},{"Data":{"Text":"Écrire sous une seule fraction : "},"Kind":4},{"Data":{"Content":"\\\\frac{1}{3} + \\\\frac{2}{5}","IsInline":true},"Kind":0},{"Data":{"Text":"Écrire sous une seule fraction : "},"Kind":4},{"Data":{"Content":"\\\\frac{1}{3} + \\\\frac{2}{5}","IsInline":true
},"Kind":0}]}
""",
    """
{"Title":"Très longue question verticale","Content":[{"Data":{"Text":"Écrire sous une seule fraction : "},"Kind":4},{"Data":{"Content":"\\\\frac{1}{3} + \\\\frac{2}{5}","IsInline":false},"Kind":0},{"Data":{"Text":"Écrire sous une seule fraction : "},"Kind":4},{"Data":{"Content":"\\\\frac{1}{3} + \\\\frac{2}{5}","IsInline":false},"Kind":0},{"Data":{"Text":"Écrire sous une seule fraction : "},"Kind":4},{"Data":{"Content":"\\\\frac{1}{3} + \\\\frac{2}{5}","IsInline":false},"Kind":0},{"Data":{"Text":"Écrire sous une seule fraction : "},"Kind":4},{"Data":{"Content":"\\\\frac{1}{3} + \\\\frac{2}{5}","IsInline":false},"Kind":0},{"Data":{"Text":"Écrire sous une seule fraction : "},"Kind":4},{"Data":{"Content":"\\\\frac{1}{3} + \\\\frac{2}{5}","IsInline":false},"Kind":0},{"Data":{"Text":"Écrire sous une seule fraction : "},"Kind":4},{"Data":{"Content":"\\\\frac{1}{3} + \\\\frac{2}{5}","IsInline":false},"Kind":0},{"Data":{"Text":"Écrire sous une seule fraction : "},"Kind":4},{"Data":{"Content":"\\\\frac{1}{3} + \\\\frac{2}{5}","IsInline":
false},"Kind":0}]}
""",
  ];

  Question _fromJSON(String json) {
    final question = clientQuestionFromJson(jsonDecode(json));
    return Question(
        question, Categorie.values[Random().nextInt(Categorie.values.length)]);
  }

  @override
  Widget build(BuildContext context) {
    return PageView(
      controller: _controller,
      children: _questions.map(_fromJSON).toList(),
    );
  }
}
