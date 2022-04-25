import 'dart:math';

const _quotes = [
  "Apprendre par coeur est bien, apprendre par le coeur est mieux. (Paul Masson)",
  "On ne peut faire saigner un coeur sans en être éclaboussé. (Paul Masson)",
  "N'oublie jamais : tu as de la valeur !",
];

String pickQuote() {
  final index = Random().nextInt(_quotes.length);
  return _quotes[index];
}
