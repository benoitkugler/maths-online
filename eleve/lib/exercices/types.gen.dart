// Code generated by structgen. DO NOT EDIT

typedef JSON = Map<String, dynamic>; // alias to shorten JSON convertors

String stringFromJson(dynamic json) => json as String;

String stringToJson(String item) => item;

// github.com/benoitkugler/maths-online/maths/exercice/client.ExpressionAnswer
class ExpressionAnswer implements Answer {
  final String expression;

  const ExpressionAnswer(this.expression);

  @override
  String toString() {
    return "ExpressionAnswer($expression)";
  }
}

ExpressionAnswer expressionAnswerFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return ExpressionAnswer(stringFromJson(json['Expression']));
}

JSON expressionAnswerToJson(ExpressionAnswer item) {
  return {"Expression": stringToJson(item.expression)};
}

double doubleFromJson(dynamic json) => json as double;

double doubleToJson(double item) => item;

// github.com/benoitkugler/maths-online/maths/exercice/client.NumberAnswer
class NumberAnswer implements Answer {
  final double value;

  const NumberAnswer(this.value);

  @override
  String toString() {
    return "NumberAnswer($value)";
  }
}

NumberAnswer numberAnswerFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return NumberAnswer(doubleFromJson(json['Value']));
}

JSON numberAnswerToJson(NumberAnswer item) {
  return {"Value": doubleToJson(item.value)};
}

int intFromJson(dynamic json) => json as int;

int intToJson(int item) => item;

// github.com/benoitkugler/maths-online/maths/exercice/client.RadioAnswer
class RadioAnswer implements Answer {
  final int index;

  const RadioAnswer(this.index);

  @override
  String toString() {
    return "RadioAnswer($index)";
  }
}

RadioAnswer radioAnswerFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return RadioAnswer(intFromJson(json['Index']));
}

JSON radioAnswerToJson(RadioAnswer item) {
  return {"Index": intToJson(item.index)};
}

abstract class Answer {}

Answer answerFromJson(dynamic json_) {
  final json = json_ as JSON;
  final kind = json['Kind'] as int;
  final data = json['Data'];
  switch (kind) {
    case 0:
      return expressionAnswerFromJson(data);
    case 1:
      return numberAnswerFromJson(data);
    case 2:
      return radioAnswerFromJson(data);
    default:
      throw ("unexpected type");
  }
}

JSON answerToJson(Answer item) {
  if (item is ExpressionAnswer) {
    return {'Kind': 0, 'Data': expressionAnswerToJson(item)};
  } else if (item is NumberAnswer) {
    return {'Kind': 1, 'Data': numberAnswerToJson(item)};
  } else if (item is RadioAnswer) {
    return {'Kind': 2, 'Data': radioAnswerToJson(item)};
  } else {
    throw ("unexpected type");
  }
}

// github.com/benoitkugler/maths-online/maths/exercice/client.ExpressionFieldBlock
class ExpressionFieldBlock implements Block {
  final String label;
  final int iD;

  const ExpressionFieldBlock(this.label, this.iD);

  @override
  String toString() {
    return "ExpressionFieldBlock($label, $iD)";
  }
}

ExpressionFieldBlock expressionFieldBlockFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return ExpressionFieldBlock(
      stringFromJson(json['Label']), intFromJson(json['ID']));
}

JSON expressionFieldBlockToJson(ExpressionFieldBlock item) {
  return {"Label": stringToJson(item.label), "ID": intToJson(item.iD)};
}

// github.com/benoitkugler/maths-online/maths/exercice/client.FormulaBlock
class FormulaBlock implements Block {
  final String formula;

  const FormulaBlock(this.formula);

  @override
  String toString() {
    return "FormulaBlock($formula)";
  }
}

FormulaBlock formulaBlockFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return FormulaBlock(stringFromJson(json['Formula']));
}

JSON formulaBlockToJson(FormulaBlock item) {
  return {"Formula": stringToJson(item.formula)};
}

List<String> listStringFromJson(dynamic json) {
  if (json == null) {
    return [];
  }
  return (json as List<dynamic>).map(stringFromJson).toList();
}

List<dynamic> listStringToJson(List<String> item) {
  return item.map(stringToJson).toList();
}

// github.com/benoitkugler/maths-online/maths/exercice/client.ListFieldBlock
class ListFieldBlock implements Block {
  final List<String> choices;
  final int iD;

  const ListFieldBlock(this.choices, this.iD);

  @override
  String toString() {
    return "ListFieldBlock($choices, $iD)";
  }
}

ListFieldBlock listFieldBlockFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return ListFieldBlock(
      listStringFromJson(json['Choices']), intFromJson(json['ID']));
}

JSON listFieldBlockToJson(ListFieldBlock item) {
  return {"Choices": listStringToJson(item.choices), "ID": intToJson(item.iD)};
}

// github.com/benoitkugler/maths-online/maths/exercice/client.NumberFieldBlock
class NumberFieldBlock implements Block {
  final int iD;

  const NumberFieldBlock(this.iD);

  @override
  String toString() {
    return "NumberFieldBlock($iD)";
  }
}

NumberFieldBlock numberFieldBlockFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return NumberFieldBlock(intFromJson(json['ID']));
}

JSON numberFieldBlockToJson(NumberFieldBlock item) {
  return {"ID": intToJson(item.iD)};
}

bool boolFromJson(dynamic json) => json as bool;

bool boolToJson(bool item) => item;

// github.com/benoitkugler/maths-online/maths/exercice/client.TextOrMath
class TextOrMath {
  final String text;
  final bool isMath;

  const TextOrMath(this.text, this.isMath);

  @override
  String toString() {
    return "TextOrMath($text, $isMath)";
  }
}

TextOrMath textOrMathFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return TextOrMath(stringFromJson(json['Text']), boolFromJson(json['IsMath']));
}

JSON textOrMathToJson(TextOrMath item) {
  return {"Text": stringToJson(item.text), "IsMath": boolToJson(item.isMath)};
}

List<TextOrMath> listTextOrMathFromJson(dynamic json) {
  if (json == null) {
    return [];
  }
  return (json as List<dynamic>).map(textOrMathFromJson).toList();
}

List<dynamic> listTextOrMathToJson(List<TextOrMath> item) {
  return item.map(textOrMathToJson).toList();
}

// github.com/benoitkugler/maths-online/maths/exercice/client.ListFieldProposal
class ListFieldProposal {
  final List<TextOrMath> content;

  const ListFieldProposal(this.content);

  @override
  String toString() {
    return "ListFieldProposal($content)";
  }
}

ListFieldProposal listFieldProposalFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return ListFieldProposal(listTextOrMathFromJson(json['Content']));
}

JSON listFieldProposalToJson(ListFieldProposal item) {
  return {"Content": listTextOrMathToJson(item.content)};
}

List<ListFieldProposal> listListFieldProposalFromJson(dynamic json) {
  if (json == null) {
    return [];
  }
  return (json as List<dynamic>).map(listFieldProposalFromJson).toList();
}

List<dynamic> listListFieldProposalToJson(List<ListFieldProposal> item) {
  return item.map(listFieldProposalToJson).toList();
}

// github.com/benoitkugler/maths-online/maths/exercice/client.RadioFieldBlock
class RadioFieldBlock implements Block {
  final List<ListFieldProposal> proposals;
  final int iD;

  const RadioFieldBlock(this.proposals, this.iD);

  @override
  String toString() {
    return "RadioFieldBlock($proposals, $iD)";
  }
}

RadioFieldBlock radioFieldBlockFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return RadioFieldBlock(listListFieldProposalFromJson(json['Proposals']),
      intFromJson(json['ID']));
}

JSON radioFieldBlockToJson(RadioFieldBlock item) {
  return {
    "Proposals": listListFieldProposalToJson(item.proposals),
    "ID": intToJson(item.iD)
  };
}

// github.com/benoitkugler/maths-online/maths/exercice/client.TextBlock
class TextBlock implements Block {
  final List<TextOrMath> parts;

  const TextBlock(this.parts);

  @override
  String toString() {
    return "TextBlock($parts)";
  }
}

TextBlock textBlockFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return TextBlock(listTextOrMathFromJson(json['Parts']));
}

JSON textBlockToJson(TextBlock item) {
  return {"Parts": listTextOrMathToJson(item.parts)};
}

abstract class Block {}

Block blockFromJson(dynamic json_) {
  final json = json_ as JSON;
  final kind = json['Kind'] as int;
  final data = json['Data'];
  switch (kind) {
    case 0:
      return expressionFieldBlockFromJson(data);
    case 1:
      return formulaBlockFromJson(data);
    case 2:
      return listFieldBlockFromJson(data);
    case 3:
      return numberFieldBlockFromJson(data);
    case 4:
      return radioFieldBlockFromJson(data);
    case 5:
      return textBlockFromJson(data);
    default:
      throw ("unexpected type");
  }
}

JSON blockToJson(Block item) {
  if (item is ExpressionFieldBlock) {
    return {'Kind': 0, 'Data': expressionFieldBlockToJson(item)};
  } else if (item is FormulaBlock) {
    return {'Kind': 1, 'Data': formulaBlockToJson(item)};
  } else if (item is ListFieldBlock) {
    return {'Kind': 2, 'Data': listFieldBlockToJson(item)};
  } else if (item is NumberFieldBlock) {
    return {'Kind': 3, 'Data': numberFieldBlockToJson(item)};
  } else if (item is RadioFieldBlock) {
    return {'Kind': 4, 'Data': radioFieldBlockToJson(item)};
  } else if (item is TextBlock) {
    return {'Kind': 5, 'Data': textBlockToJson(item)};
  } else {
    throw ("unexpected type");
  }
}

List<Block> listBlockFromJson(dynamic json) {
  if (json == null) {
    return [];
  }
  return (json as List<dynamic>).map(blockFromJson).toList();
}

List<dynamic> listBlockToJson(List<Block> item) {
  return item.map(blockToJson).toList();
}

// github.com/benoitkugler/maths-online/maths/exercice/client.Enonce
typedef Enonce = List<Block>;

// github.com/benoitkugler/maths-online/maths/exercice/client.Question
class Question {
  final String title;
  final Enonce enonce;

  const Question(this.title, this.enonce);

  @override
  String toString() {
    return "Question($title, $enonce)";
  }
}

Question questionFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return Question(
      stringFromJson(json['Title']), listBlockFromJson(json['Enonce']));
}

JSON questionToJson(Question item) {
  return {
    "Title": stringToJson(item.title),
    "Enonce": listBlockToJson(item.enonce)
  };
}

Map<int, Answer> dictIntAnswerFromJson(dynamic json) {
  if (json == null) {
    return {};
  }
  return (json as JSON)
      .map((k, v) => MapEntry(int.parse(k), answerFromJson(v)));
}

Map<String, dynamic> dictIntAnswerToJson(Map<int, Answer> item) {
  return item.map((k, v) => MapEntry(intToJson(k).toString(), answerToJson(v)));
}

// github.com/benoitkugler/maths-online/maths/exercice/client.QuestionAnswersIn
class QuestionAnswersIn {
  final Map<int, Answer> data;

  const QuestionAnswersIn(this.data);

  @override
  String toString() {
    return "QuestionAnswersIn($data)";
  }
}

QuestionAnswersIn questionAnswersInFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return QuestionAnswersIn(dictIntAnswerFromJson(json['Data']));
}

JSON questionAnswersInToJson(QuestionAnswersIn item) {
  return {"Data": dictIntAnswerToJson(item.data)};
}

Map<int, bool> dictIntBoolFromJson(dynamic json) {
  if (json == null) {
    return {};
  }
  return (json as JSON).map((k, v) => MapEntry(int.parse(k), boolFromJson(v)));
}

Map<String, dynamic> dictIntBoolToJson(Map<int, bool> item) {
  return item.map((k, v) => MapEntry(intToJson(k).toString(), boolToJson(v)));
}

// github.com/benoitkugler/maths-online/maths/exercice/client.QuestionAnswersOut
class QuestionAnswersOut {
  final Map<int, bool> data;

  const QuestionAnswersOut(this.data);

  @override
  String toString() {
    return "QuestionAnswersOut($data)";
  }
}

QuestionAnswersOut questionAnswersOutFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return QuestionAnswersOut(dictIntBoolFromJson(json['Data']));
}

JSON questionAnswersOutToJson(QuestionAnswersOut item) {
  return {"Data": dictIntBoolToJson(item.data)};
}

// github.com/benoitkugler/maths-online/maths/exercice/client.QuestionSyntaxCheckIn
class QuestionSyntaxCheckIn {
  final Answer answer;
  final int iD;

  const QuestionSyntaxCheckIn(this.answer, this.iD);

  @override
  String toString() {
    return "QuestionSyntaxCheckIn($answer, $iD)";
  }
}

QuestionSyntaxCheckIn questionSyntaxCheckInFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return QuestionSyntaxCheckIn(
      answerFromJson(json['Answer']), intFromJson(json['ID']));
}

JSON questionSyntaxCheckInToJson(QuestionSyntaxCheckIn item) {
  return {"Answer": answerToJson(item.answer), "ID": intToJson(item.iD)};
}

// github.com/benoitkugler/maths-online/maths/exercice/client.QuestionSyntaxCheckOut
class QuestionSyntaxCheckOut {
  final String reason;
  final bool isValid;

  const QuestionSyntaxCheckOut(this.reason, this.isValid);

  @override
  String toString() {
    return "QuestionSyntaxCheckOut($reason, $isValid)";
  }
}

QuestionSyntaxCheckOut questionSyntaxCheckOutFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return QuestionSyntaxCheckOut(
      stringFromJson(json['Reason']), boolFromJson(json['IsValid']));
}

JSON questionSyntaxCheckOutToJson(QuestionSyntaxCheckOut item) {
  return {
    "Reason": stringToJson(item.reason),
    "IsValid": boolToJson(item.isValid)
  };
}