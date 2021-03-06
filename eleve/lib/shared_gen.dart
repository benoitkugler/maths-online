// Code generated by structgen. DO NOT EDIT

import 'questions/types.gen.dart';

typedef JSON = Map<String, dynamic>; // alias to shorten JSON convertors

String stringFromJson(dynamic json) => json == null ? "" : json as String;

String stringToJson(String item) => item;

bool boolFromJson(dynamic json) => json as bool;

bool boolToJson(bool item) => item;

// github.com/benoitkugler/maths-online.CheckExpressionOut
class CheckExpressionOut {
  final String reason;
  final bool isValid;

  const CheckExpressionOut(this.reason, this.isValid);

  @override
  String toString() {
    return "CheckExpressionOut($reason, $isValid)";
  }
}

CheckExpressionOut checkExpressionOutFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return CheckExpressionOut(
      stringFromJson(json['Reason']), boolFromJson(json['IsValid']));
}

JSON checkExpressionOutToJson(CheckExpressionOut item) {
  return {
    "Reason": stringToJson(item.reason),
    "IsValid": boolToJson(item.isValid)
  };
}

int intFromJson(dynamic json) => json as int;

int intToJson(int item) => item;

// github.com/benoitkugler/maths-online/maths/expression.Variable
class Variable {
  final String indice;
  final int name;

  const Variable(this.indice, this.name);

  @override
  String toString() {
    return "Variable($indice, $name)";
  }
}

Variable variableFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return Variable(stringFromJson(json['Indice']), intFromJson(json['Name']));
}

JSON variableToJson(Variable item) {
  return {"Indice": stringToJson(item.indice), "Name": intToJson(item.name)};
}

// github.com/benoitkugler/maths-online/prof/editor.VarEntry
class VarEntry {
  final Variable variable;
  final String resolved;

  const VarEntry(this.variable, this.resolved);

  @override
  String toString() {
    return "VarEntry($variable, $resolved)";
  }
}

VarEntry varEntryFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return VarEntry(
      variableFromJson(json['Variable']), stringFromJson(json['Resolved']));
}

JSON varEntryToJson(VarEntry item) {
  return {
    "Variable": variableToJson(item.variable),
    "Resolved": stringToJson(item.resolved)
  };
}

List<VarEntry> listVarEntryFromJson(dynamic json) {
  if (json == null) {
    return [];
  }
  return (json as List<dynamic>).map(varEntryFromJson).toList();
}

List<dynamic> listVarEntryToJson(List<VarEntry> item) {
  return item.map(varEntryToJson).toList();
}

// github.com/benoitkugler/maths-online/prof/editor.Answer
class Answer {
  final List<VarEntry> params;
  final QuestionAnswersIn answer;

  const Answer(this.params, this.answer);

  @override
  String toString() {
    return "Answer($params, $answer)";
  }
}

Answer answerFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return Answer(listVarEntryFromJson(json['Params']),
      questionAnswersInFromJson(json['Answer']));
}

JSON answerToJson(Answer item) {
  return {
    "Params": listVarEntryToJson(item.params),
    "Answer": questionAnswersInToJson(item.answer)
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

// github.com/benoitkugler/maths-online/prof/editor.Progression
class Progression {
  final int id;
  final int id_exercice;

  const Progression(this.id, this.id_exercice);

  @override
  String toString() {
    return "Progression($id, $id_exercice)";
  }
}

Progression progressionFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return Progression(intFromJson(json['Id']), intFromJson(json['id_exercice']));
}

JSON progressionToJson(Progression item) {
  return {"Id": intToJson(item.id), "id_exercice": intToJson(item.id_exercice)};
}

List<bool> listBoolFromJson(dynamic json) {
  if (json == null) {
    return [];
  }
  return (json as List<dynamic>).map(boolFromJson).toList();
}

List<dynamic> listBoolToJson(List<bool> item) {
  return item.map(boolToJson).toList();
}

// github.com/benoitkugler/maths-online/prof/editor.QuestionHistory
typedef QuestionHistory = List<bool>;

List<QuestionHistory> listListBoolFromJson(dynamic json) {
  if (json == null) {
    return [];
  }
  return (json as List<dynamic>).map(listBoolFromJson).toList();
}

List<dynamic> listListBoolToJson(List<QuestionHistory> item) {
  return item.map(listBoolToJson).toList();
}

// github.com/benoitkugler/maths-online/prof/editor.ProgressionExt
class ProgressionExt {
  final Progression progression;
  final List<QuestionHistory> questions;
  final int nextQuestion;

  const ProgressionExt(this.progression, this.questions, this.nextQuestion);

  @override
  String toString() {
    return "ProgressionExt($progression, $questions, $nextQuestion)";
  }
}

ProgressionExt progressionExtFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return ProgressionExt(
      progressionFromJson(json['Progression']),
      listListBoolFromJson(json['Questions']),
      intFromJson(json['NextQuestion']));
}

JSON progressionExtToJson(ProgressionExt item) {
  return {
    "Progression": progressionToJson(item.progression),
    "Questions": listListBoolToJson(item.questions),
    "NextQuestion": intToJson(item.nextQuestion)
  };
}

// github.com/benoitkugler/maths-online/prof/editor.EvaluateExerciceIn
class EvaluateExerciceIn {
  final int idExercice;
  final Map<int, Answer> answers;
  final ProgressionExt progression;

  const EvaluateExerciceIn(this.idExercice, this.answers, this.progression);

  @override
  String toString() {
    return "EvaluateExerciceIn($idExercice, $answers, $progression)";
  }
}

EvaluateExerciceIn evaluateExerciceInFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return EvaluateExerciceIn(
      intFromJson(json['IdExercice']),
      dictIntAnswerFromJson(json['Answers']),
      progressionExtFromJson(json['Progression']));
}

JSON evaluateExerciceInToJson(EvaluateExerciceIn item) {
  return {
    "IdExercice": intToJson(item.idExercice),
    "Answers": dictIntAnswerToJson(item.answers),
    "Progression": progressionExtToJson(item.progression)
  };
}

Map<int, QuestionAnswersOut> dictIntQuestionAnswersOutFromJson(dynamic json) {
  if (json == null) {
    return {};
  }
  return (json as JSON)
      .map((k, v) => MapEntry(int.parse(k), questionAnswersOutFromJson(v)));
}

Map<String, dynamic> dictIntQuestionAnswersOutToJson(
    Map<int, QuestionAnswersOut> item) {
  return item.map(
      (k, v) => MapEntry(intToJson(k).toString(), questionAnswersOutToJson(v)));
}

// github.com/benoitkugler/maths-online/prof/editor.InstantiatedQuestion
class InstantiatedQuestion {
  final int id;
  final Question question;
  final List<VarEntry> params;

  const InstantiatedQuestion(this.id, this.question, this.params);

  @override
  String toString() {
    return "InstantiatedQuestion($id, $question, $params)";
  }
}

InstantiatedQuestion instantiatedQuestionFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return InstantiatedQuestion(intFromJson(json['Id']),
      questionFromJson(json['Question']), listVarEntryFromJson(json['Params']));
}

JSON instantiatedQuestionToJson(InstantiatedQuestion item) {
  return {
    "Id": intToJson(item.id),
    "Question": questionToJson(item.question),
    "Params": listVarEntryToJson(item.params)
  };
}

List<InstantiatedQuestion> listInstantiatedQuestionFromJson(dynamic json) {
  if (json == null) {
    return [];
  }
  return (json as List<dynamic>).map(instantiatedQuestionFromJson).toList();
}

List<dynamic> listInstantiatedQuestionToJson(List<InstantiatedQuestion> item) {
  return item.map(instantiatedQuestionToJson).toList();
}

// github.com/benoitkugler/maths-online/prof/editor.EvaluateExerciceOut
class EvaluateExerciceOut {
  final Map<int, QuestionAnswersOut> results;
  final ProgressionExt progression;
  final List<InstantiatedQuestion> newQuestions;

  const EvaluateExerciceOut(this.results, this.progression, this.newQuestions);

  @override
  String toString() {
    return "EvaluateExerciceOut($results, $progression, $newQuestions)";
  }
}

EvaluateExerciceOut evaluateExerciceOutFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return EvaluateExerciceOut(
      dictIntQuestionAnswersOutFromJson(json['Results']),
      progressionExtFromJson(json['Progression']),
      listInstantiatedQuestionFromJson(json['NewQuestions']));
}

JSON evaluateExerciceOutToJson(EvaluateExerciceOut item) {
  return {
    "Results": dictIntQuestionAnswersOutToJson(item.results),
    "Progression": progressionExtToJson(item.progression),
    "NewQuestions": listInstantiatedQuestionToJson(item.newQuestions)
  };
}

// github.com/benoitkugler/maths-online.EvaluateQuestionIn
class EvaluateQuestionIn {
  final QuestionAnswersIn answer;
  final List<VarEntry> params;
  final int idQuestion;

  const EvaluateQuestionIn(this.answer, this.params, this.idQuestion);

  @override
  String toString() {
    return "EvaluateQuestionIn($answer, $params, $idQuestion)";
  }
}

EvaluateQuestionIn evaluateQuestionInFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return EvaluateQuestionIn(questionAnswersInFromJson(json['Answer']),
      listVarEntryFromJson(json['Params']), intFromJson(json['IdQuestion']));
}

JSON evaluateQuestionInToJson(EvaluateQuestionIn item) {
  return {
    "Answer": questionAnswersInToJson(item.answer),
    "Params": listVarEntryToJson(item.params),
    "IdQuestion": intToJson(item.idQuestion)
  };
}

// github.com/benoitkugler/maths-online/prof/editor.Flow
enum Flow { parallel, sequencial }

extension _FlowExt on Flow {
  static Flow fromValue(int i) {
    return Flow.values[i];
  }

  int toValue() {
    return index;
  }
}

Flow flowFromJson(dynamic json) => _FlowExt.fromValue(json as int);

dynamic flowToJson(Flow item) => item.toValue();

List<int> listIntFromJson(dynamic json) {
  if (json == null) {
    return [];
  }
  return (json as List<dynamic>).map(intFromJson).toList();
}

List<dynamic> listIntToJson(List<int> item) {
  return item.map(intToJson).toList();
}

// github.com/benoitkugler/maths-online/prof/editor.InstantiatedExercice
class InstantiatedExercice {
  final int id;
  final String title;
  final Flow flow;
  final List<InstantiatedQuestion> questions;
  final List<int> baremes;

  const InstantiatedExercice(
      this.id, this.title, this.flow, this.questions, this.baremes);

  @override
  String toString() {
    return "InstantiatedExercice($id, $title, $flow, $questions, $baremes)";
  }
}

InstantiatedExercice instantiatedExerciceFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return InstantiatedExercice(
      intFromJson(json['Id']),
      stringFromJson(json['Title']),
      flowFromJson(json['Flow']),
      listInstantiatedQuestionFromJson(json['Questions']),
      listIntFromJson(json['Baremes']));
}

JSON instantiatedExerciceToJson(InstantiatedExercice item) {
  return {
    "Id": intToJson(item.id),
    "Title": stringToJson(item.title),
    "Flow": flowToJson(item.flow),
    "Questions": listInstantiatedQuestionToJson(item.questions),
    "Baremes": listIntToJson(item.baremes)
  };
}

// github.com/benoitkugler/maths-online.Exercice
class Exercice {
  final InstantiatedExercice exercice;
  final ProgressionExt progression;

  const Exercice(this.exercice, this.progression);

  @override
  String toString() {
    return "Exercice($exercice, $progression)";
  }
}

Exercice exerciceFromJson(dynamic json_) {
  final json = (json_ as JSON);
  return Exercice(instantiatedExerciceFromJson(json['Exercice']),
      progressionExtFromJson(json['Progression']));
}

JSON exerciceToJson(Exercice item) {
  return {
    "Exercice": instantiatedExerciceToJson(item.exercice),
    "Progression": progressionExtToJson(item.progression)
  };
}

// github.com/benoitkugler/maths-online/prof/editor.InstantiateQuestionsOut
typedef InstantiateQuestionsOut = List<InstantiatedQuestion>;
