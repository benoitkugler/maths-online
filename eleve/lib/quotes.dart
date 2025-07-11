import 'dart:math';

import 'package:flutter/material.dart';

class QuoteData {
  final String content;
  final String comment; // may be empty
  final String author; // may be empty
  const QuoteData(this.content, this.comment, this.author);

  bool get isEmpty => content.isEmpty && comment.isEmpty && author.isEmpty;
}

final _quotes = [
  const QuoteData(
    "Je ne perds jamais : soit je gagne, soit j'apprends.",
    "",
    "Nelson Mandela",
  ),
  const QuoteData(
    "Le bonheur n'est pas d'avoir tout ce que l'on désire mais d'apprécier ce que l'on a.",
    "",
    "Paulo Coelho",
  ),
  const QuoteData(
    "La vie mettra des pierres sur ta route. A toi de décider d'en faire des murs ou des ponts.",
    "",
    "Coluche",
  ),
  const QuoteData(
    "Un cœur prompt à la colère est en grand danger.",
    "Ne t'énerve pas si vite !",
    "Ménandre",
  ),
  const QuoteData(
      "On ne peut faire saigner un cœur sans en être éclaboussé.",
      "Ne crois pas que les crasses faites à l'autre ne t'atteignent pas aussi !",
      "Paul Masson"),
  const QuoteData(
      "Apprendre par cœur est bien, apprendre par le cœur est mieux.",
      "Apprendre par cœur c'est bien, comprendre c'est mieux !",
      "Paul Masson"),
  const QuoteData(
      "Que de sottises du cœur viennent de sa fatigue !",
      "Il vaut donc mieux se taire que de penser ou dire dans ces conditions !",
      "Anne Barratin"),
  const QuoteData(
    "Le coeur fait tout, le reste est inutile.",
    "",
    "Jean de la Fontaine",
  ),
  const QuoteData(
    "Qu'est-ce qui rempli tout le coeur ? L'amour.",
    "",
    "Victor Hugo",
  ),
  const QuoteData(
    "La chute n'est pas un échec. L'échec, c'est de rester là où l'on est tombé.",
    "",
    "Socrate",
  ),
  const QuoteData(
    "Les mots peuvent être mortels... fais-en bon usage !",
    "",
    "Jean Chalon",
  ),
  const QuoteData(
      "Les mots peuvent blesser un cœur, tuer un rêve, briser une relation : veille sur ta bouche !",
      "",
      ""),
  const QuoteData(
      "On dit souvent de prendre son mal en patience. Et si on prenait son bien en urgence ?",
      "",
      ""),
  const QuoteData(
      "Ne laisse pas le monde changer ton sourire mais laisse ton sourire changer le monde !",
      "",
      ""),
  const QuoteData(
      "Si tu vis par les compliments des gens tu mourras par leurs critiques.",
      "",
      ""),
  const QuoteData(
      "Sois la raison pour laquelle quelqu'un sourit aujourd'hui :)", "", ""),
  const QuoteData(
      "Le processus de ton coeur compte plus que la finalité de tes projets.",
      "",
      ""),
  const QuoteData(
      "Je ne cherche pas à être qui j’étais avant, je cherche à exprimer qui je suis maintenant !",
      "",
      ""),
  const QuoteData(
      "L’erreur c’est de croire que l’on a besoin de rien et de se satisfaire d’une vie sans lendemain.",
      "",
      ""),
  const QuoteData("Qu’importe, l’aube poindra.", "", ""),
  const QuoteData("Choisir c'est renoncer.", "", ""),
  const QuoteData(
      "Autrui fait de moi ce que je suis... parfois ce que je ne suis pas…",
      "",
      ""),
  const QuoteData(
      "Le problème n'est pas d'être têtu mais de l'être pour les bonnes choses.",
      "",
      ""),
  const QuoteData(
      "N’oublie jamais : ta vie a de la valeur, elle est précieuse…", "", ""),
  const QuoteData(
      "Ta vie est trop précieuse pour être jetée en l’air.", "", ""),
  const QuoteData(
      "Quelle que soit la situation, ne perds pas espoir !", "", ""),
  const QuoteData(
      "Ne dis jamais que tu n'as pas de valeur... tu as simplement besoin de la découvrir",
      "",
      ""),
  const QuoteData(
      "Tout est possible, laisse toi surprendre... et si on essayait", "", ""),
  const QuoteData(
      "Etre en apprentissage, ce n'est pas être incompétent.", "", ""),
  const QuoteData("Sois ressourcé et fortifié !", "", ""),
  const QuoteData("Garde le cap, persévère !", "", ""),
  const QuoteData(
      "Ta valeur n'est pas définie par ce que tu réussis ou non !", "", ""),
  const QuoteData("Ne te fies pas à tout ce que tu penses ou ressens.", "", ""),
  const QuoteData(
      "Ne compte pas les choses que tu fais mais fais les choses qui comptent.",
      "",
      ""),
  const QuoteData(
      "On ne voit bien qu'avec le coeur. L'essentiel est invisible pour les yeux.",
      "",
      "Le petit prince"),
  const QuoteData(
      "Il faut toute la vie pour apprendre à vivre.", "", " Sénèque"),
  const QuoteData(
      "La vie n'est pas d'attendre que les orages passent, c'est d'apprendre à danser sous la pluie.",
      "",
      "Sénèque"),
  const QuoteData(
      "Je refuse que la peur de l'échec m'empêche de faire ce qui importe vraiment.",
      "",
      ""),
  const QuoteData("Ne comptez pas les jours, faites que chaque jour compte.",
      "", "Muhammed Ali"),
  const QuoteData(
      "Les décisions les plus difficiles à prendre sont celles qui vous présentent des chemins au bout desquels vous ne serez plus la même personne.",
      "",
      " Nelson Mandela"),
  const QuoteData("Vis comme si tu n'avais besoin de rien de plus", "", ""),
  const QuoteData(
      "Le manque de temps n'est rien d'autre que le manque de priorité.",
      "",
      "Timothy Ferriss"),
  const QuoteData(
      "La seule limite à notre épanouissement de demain sera nos doutes d'aujourd'hui.",
      "",
      "Franklin D.Roosevelt"),
  const QuoteData(
      "Le changement est une porte qui ne s'ouvre que de l'intérieur.",
      "",
      "Tom Peters"),
  const QuoteData("Le bonheur le plus doux est celui que l'on partage.", "",
      "Jacques Delille"),
  const QuoteData(
      "Un enseignant, un livre, un stylo, peuvent changer le monde.",
      "",
      "Malala Yousafzai"),
  const QuoteData(
      "Le monde déteste le changement. C'est pourtant la première chose qui lui a permis de progresser.",
      "",
      "Charles F.Kettering"),
  const QuoteData(
      "L'art de vivre consiste en un subtil mélange entre lâcher prise et tenir bon.",
      "",
      "Henri Lewis"),
  const QuoteData("Il est grand temps de rallumer les étoiles.", "",
      "Guillaume Apollinaire"),
  const QuoteData("Changer, c'est d'abord changer de point de vue.", "",
      "Jean Bertrand Pontalis"),
  const QuoteData(
      "Celui qui déplace une montagne commence par déplacer de petites pierres.",
      "",
      "Confucius"),
  const QuoteData(
      "Attendre d'en savoir assez pour agir en toute lumière, c'est se condamner à l'inaction.",
      "",
      "Jean Rostand"),
  const QuoteData(
      "L'expérience ce n'est pas ce qui nous arrive, c'est ce que nous faisons avec ce qui nous arrive.",
      "",
      "Aldous Huxley"),
  const QuoteData(
      "Le bonheur n'est pas une destination à atteindre, c'est une manière de voyager.",
      "",
      "Margaret Lee Runbeck"),
  const QuoteData(
      "Les mots sont plus forts que les armes.", "", "Laurent Jacqua"),
  const QuoteData(
      "Prendre conscience de son ignorance est un grand pas vers la connaissance.",
      "",
      "Benjamin Disraeli"),
  const QuoteData(
      "On aide plus un être en lui donnant de lui-même une image favorable qu'en le mettant sans cesse en face de ses défauts.",
      "",
      "Albert Camus"),
  const QuoteData(
      "La valeur d'un homme tient dans sa capacité à donner et non dans sa capacité à recevoir.",
      "",
      "Albert Einstein"),
  const QuoteData(
      "Il n'est jamais trop tard pour devenir ce que vous auriez toujours dû être.",
      "",
      "George Eliot"),
  const QuoteData("Les grands changements viennent des petites choses.", "",
      "Paulo Coelho"),
  const QuoteData(
      "Ne dites pas tout ce que vous pensez, mais pensez à tout ce que vous dites.",
      "",
      "Carlos Martinez Vazquez"),
  const QuoteData(
      "Quand un homme dit 'Je suis heureux', il veut tout bonnement dire 'J'ai des ennuis qui ne m'atteignent pas'.",
      "",
      "Jules Renard"),
  const QuoteData(
      "Tu ne peux pas empêcher un oiseau de voler au-dessus de ta tête. Par contre, tu peux l'empêcher d'y établir son nid.",
      "",
      "Auteur inconnu"),
  const QuoteData(
      "Le succès, c'est d'aller d'échec en échec sans perdre son enthousiasme.",
      "",
      "Winston Churchill"),
  const QuoteData(
      "Saviez-vous que le lion, lorsqu'il sort pour chasser, échoue 7 à 10 fois avant de réussir à capturer une proie ? 85% de sa vie est un échec. Alors qu'est-ce qui fait de lui un roi ? Sa persévérance.",
      "",
      ""),
  const QuoteData(
      "Aucun de nous, en agissant seul, ne peut atteindre le succès.",
      "",
      "Nelson Mandela"),
  const QuoteData("On est pas ici-bas pour se faire du tracas", "", ""),
  const QuoteData(
      "Ne regarde pas en arrière, ce n'est pas là que tu vas.", "", "")
];

var _quoteIndex = 0;

void initQuotes() {
  /// to make sure we don't repeat the same quote,
  /// we shuffle the slice at startup, and then simply
  /// pick quotes sequencially
  _quotes.shuffle(Random(DateTime.now().millisecondsSinceEpoch));
  _quoteIndex = 0;
}

QuoteData pickQuote() {
  final out = _quotes[_quoteIndex % _quotes.length];
  _quoteIndex++;
  return out;
}

class Quote extends StatelessWidget {
  final QuoteData data;
  const Quote(this.data, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    const shadows = [
      BoxShadow(color: Colors.white, blurRadius: 6, spreadRadius: 3)
    ];
    const fontSize = 16.0;
    return Container(
      decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(4)),
          border: Border.all(color: Colors.lightBlueAccent),
          color: Theme.of(context).scaffoldBackgroundColor,
          boxShadow: const [
            BoxShadow(
                color: Colors.lightBlueAccent, blurRadius: 3, spreadRadius: 2)
          ]),
      margin: const EdgeInsets.symmetric(vertical: 12, horizontal: 4),
      padding: const EdgeInsets.all(8),
      child: Column(
        children: data.author.isEmpty
            ? [
                Text(
                  data.content,
                  style: const TextStyle(shadows: shadows, fontSize: fontSize),
                )
              ]
            : [
                Text("${data.content} ${data.author}",
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
