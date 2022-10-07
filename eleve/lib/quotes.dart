import 'dart:math';

import 'package:flutter/material.dart';

class QuoteData {
  final String content;
  final String comment; // may be empty
  final String author; // may be empty
  const QuoteData(this.content, this.comment, this.author);

  bool get isEmpty => content.isEmpty && comment.isEmpty && author.isEmpty;
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
      "Le problème n'est pas d'être têtu mais de l'être pour les bonnes choses.",
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
  QuoteData("Sois ressourcé et fortifié !", "", ""),
  QuoteData("Garde le cap, persévère !", "", ""),
  QuoteData(
      "Ta valeur n'est pas définie par ce que tu réussis ou non !", "", ""),
  QuoteData("Ne te fies pas à tout ce que tu penses ou ressens.", "", ""),
  QuoteData(
      "Ne compte pas les choses que tu fais mais fais les choses qui comptent.",
      "",
      ""),
  QuoteData(
      "On ne voit bien qu'avec le coeur. L'essentiel est invisible pour les yeux.",
      "",
      "Le petit prince"),
  QuoteData("Il faut toute la vie pour apprendre à vivre.", "", " Sénèque"),
  QuoteData(
      "La vie n'est pas d'attendre que les orages passent, c'est d'apprendre à danser sous la pluie.",
      "",
      "Sénèque"),
  QuoteData(
      "Je refuse que la peur de l'échec m'empêche de faire ce qui importe vraiment.",
      "",
      ""),
  QuoteData("Ne comptez pas les jours, faites que chaque jour compte.", "",
      "Muhammed Ali"),
  QuoteData(
      "Les décisions les plus difficiles à prendre sont celles qui vous présentent des chemins au bout desquels vous ne serez plus la même personne.",
      "",
      " Nelson Mandela"),
  QuoteData("Vis comme si tu n'avais besoin de rien de plus", "", ""),
  QuoteData("Le manque de temps n'est rien d'autre que le manque de priorité.",
      "", "Timothy Ferriss"),
  QuoteData(
      "La seule limite à notre épanouissement de demain sera nos doutes d'aujourd'hui.",
      "",
      "Franklin D.Roosevelt"),
  QuoteData("Le changement est une porte qui ne s'ouvre que de l'intérieur.",
      "", "Tom Peters"),
  QuoteData("Le bonheur le plus doux est celui que l'on partage.", "",
      "Jacques Delille"),
  QuoteData("Un enseignant, un livre, un stylo, peuvent changer le monde.", "",
      "Malala Yousafzai"),
  QuoteData(
      "Le monde déteste le changement. C'est pourtant la première chose qui lui a permis de progresser.",
      "",
      "Charles F.Kettering"),
  QuoteData(
      "L'art de vivre consiste en un subtil mélange entre lâcher prise et tenir bon.",
      "",
      "Henri Lewis"),
  QuoteData("Il est grand temps de rallumer les étoiles.", "",
      "Guillaume Apollinaire"),
  QuoteData("Changer, c'est d'abord changer de point de vue.", "",
      "Jean Bertrand Pontalis"),
  QuoteData(
      "Celui qui déplace une montagne commence par déplacer de petites pierres.",
      "",
      "Confucius"),
  QuoteData(
      "Attendre d'en savoir assez pour agir en toute lumière, c'est se condamner à l'inaction.",
      "",
      "Jean Rostand"),
  QuoteData(
      "L'expérience ce n'est pas ce qui nous arrive, c'est ce que nous faisons avec ce qui nous arrive.",
      "",
      "Aldous Huxley"),
  QuoteData(
      "Le bonheur n'est pas une destination à atteindre, c'est une manière de voyager.",
      "",
      "Margaret Lee Runbeck"),
  QuoteData("Les mots sont plus forts que les armes.", "", "Laurent Jacqua"),
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
