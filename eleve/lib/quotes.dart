import 'dart:math';

import 'package:flutter/material.dart';

class QuoteData {
  final String content;
  final String comment; // may be empty
  final String author; // may be empty
  const QuoteData(this.content, this.comment, this.author);
}

const _quotes = [
  QuoteData(
    "Je ne perds jamais : soit je gagne, soit j'apprends.",
    "",
    "Nelson Mandela",
  ),
  QuoteData(
    "Le bonheur n'est pas d'avoir tout ce que l'on désire mais d'apprécier ce que l'on a.",
    "",
    "Paulo Coelho",
  ),
  QuoteData(
    "La vie mettra des pierres sur ta route. A toi de décider d'en faire des murs ou des ponts.",
    "",
    "Coluche",
  ),
  QuoteData(
    "Un cœur prompt à la colère est en grand danger",
    "Ne t'énerve pas si vite !",
    "Ménandre",
  ),
  QuoteData(
      "On ne peut faire saigner un cœur sans en être éclaboussé.",
      "Ne crois pas que le mal fait à l'autre ne t'atteint pas aussi !",
      "Paul Masson"),
  QuoteData(
      "Apprendre par cœur est bien, apprendre par le cœur est mieux.",
      "Apprendre les maths c'est bien, comprendre les maths c'est mieux ",
      "Paul Masson"),
  QuoteData(
      "Que de sottises du cœur viennent de sa fatigue !",
      "Il vaut donc mieux se taire que de penser ou dire dans ces conditions !",
      "Anne Barratin"),
  QuoteData(
    "Le coeur fait tout, le reste est inutile.",
    "",
    "Jean de la Fontaine",
  ),
  QuoteData(
    "Qu'est-ce qui rempli tout le coeur ? L'amour.",
    "",
    "Victor Hugo",
  ),
  QuoteData(
    "La chute n'est pas un échec. L'échec, c'est de rester là où l'on est tombé.",
    "",
    "Socrate",
  ),
  QuoteData(
    "Les mots peuvent être mortels... fais-en bon usage !",
    "",
    "Jean Chalon",
  ),
  QuoteData(
      "Les mots peuvent blesser un cœur, tuer un rêve, briser une relation : veille sur ta bouche !",
      "",
      ""),
  QuoteData(
      "On dit souvent de prendre son mal en patience. Et si on prenait son bien en urgence ?",
      "",
      ""),
  QuoteData(
      "Ne laisse pas le monde changer ton sourire mais laisse ton sourire changer le monde !",
      "",
      ""),
  QuoteData(
      "Si tu vis par les compliments des gens tu mourras par leurs critiques.",
      "",
      ""),
  QuoteData(
      "Sois la raison pour laquelle quelqu'un sourit aujourd'hui :)", "", ""),
  QuoteData(
      "Le processus de ton coeur compte plus que la finalité de tes projets.",
      "",
      ""),
  QuoteData(
      "Je ne cherche pas à être qui j’étais avant, je cherche à exprimer qui je suis maintenant !",
      "",
      ""),
  QuoteData(
      "L’erreur c’est de croire que l’on a besoin de rien et de se satisfaire d’une vie sans lendemain.",
      "",
      ""),
  QuoteData("Qu’importe, l’aube poindra.", "", ""),
  QuoteData("Choisir c'est renoncer.", "", ""),
  QuoteData(
      "Autrui fait de moi ce que je suis... parfois ce que je ne suis pas…",
      "",
      ""),
  QuoteData(
      "N’oublie jamais : ta vie a de la valeur, elle est précieuse…", "", ""),
  QuoteData("Ta vie est trop précieuse pour être jetée en l’air.", "", ""),
  QuoteData("Quelle que soit la situation, ne perds pas espoir !", "", ""),
  QuoteData(
      "Ne dis jamais que tu n'as pas de valeur... tu as simplement besoin de la découvrir",
      "",
      ""),
  QuoteData(
      "Tout est possible, laisse toi surprendre... et si on essayait", "", ""),
  QuoteData("Etre en apprentissage, ce n'est pas être incompétent.", "", ""),
];

QuoteData pickQuote() {
  final index = Random().nextInt(_quotes.length);
  return _quotes[index];
}

class Quote extends StatelessWidget {
  final QuoteData data;
  const Quote(this.data, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    const shadows = [
      BoxShadow(color: Colors.white, blurRadius: 10, spreadRadius: 5)
    ];
    const fontSize = 16.0;
    return Padding(
      padding: const EdgeInsets.all(15.0),
      child: Column(
        children: data.author.isEmpty
            ? [
                Text(
                  data.content,
                  style: const TextStyle(shadows: shadows, fontSize: fontSize),
                )
              ]
            : [
                Text(data.content + " " + data.author,
                    style: const TextStyle(
                        fontStyle: FontStyle.italic,
                        shadows: shadows,
                        fontSize: fontSize)),
                if (data.comment.isNotEmpty) ...[
                  const SizedBox(height: 10),
                  Text(data.comment),
                ]
              ],
      ),
    );
  }
}
