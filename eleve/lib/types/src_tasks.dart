// Code generated by gomacro/generator/dart. DO NOT EDIT

import 'predefined.dart';
import 'src_maths_expression.dart';
import 'src_maths_questions_client.dart';
import 'src_sql_editor.dart';
import 'src_sql_tasks.dart';

// github.com/benoitkugler/maths-online/server/src/tasks.AnswerP
class AnswerP {
  final Params params;
  final QuestionAnswersIn answer;

  const AnswerP(this.params, this.answer);

  @override
  String toString() {
    return "AnswerP($params, $answer)";
  }
}

AnswerP answerPFromJson(dynamic json_) {
  final json = (json_ as Map<String, dynamic>);
  return AnswerP(paramsFromJson(json['Params']),
      questionAnswersInFromJson(json['Answer']));
}

Map<String, dynamic> answerPToJson(AnswerP item) {
  return {
    "Params": paramsToJson(item.params),
    "Answer": questionAnswersInToJson(item.answer)
  };
}

// github.com/benoitkugler/maths-online/server/src/tasks.EvaluateQuestionIn
class EvaluateQuestionIn {
  final AnswerP answer;
  final IdQuestion idQuestion;

  const EvaluateQuestionIn(this.answer, this.idQuestion);

  @override
  String toString() {
    return "EvaluateQuestionIn($answer, $idQuestion)";
  }
}

EvaluateQuestionIn evaluateQuestionInFromJson(dynamic json_) {
  final json = (json_ as Map<String, dynamic>);
  return EvaluateQuestionIn(
      answerPFromJson(json['Answer']), intFromJson(json['IdQuestion']));
}

Map<String, dynamic> evaluateQuestionInToJson(EvaluateQuestionIn item) {
  return {
    "Answer": answerPToJson(item.answer),
    "IdQuestion": intToJson(item.idQuestion)
  };
}

// github.com/benoitkugler/maths-online/server/src/tasks.EvaluateWorkIn
class EvaluateWorkIn {
  final WorkID iD;
  final Map<int, AnswerP> answers;
  final ProgressionExt progression;

  const EvaluateWorkIn(this.iD, this.answers, this.progression);

  @override
  String toString() {
    return "EvaluateWorkIn($iD, $answers, $progression)";
  }
}

EvaluateWorkIn evaluateWorkInFromJson(dynamic json_) {
  final json = (json_ as Map<String, dynamic>);
  return EvaluateWorkIn(
      workIDFromJson(json['ID']),
      dictIntToAnswerPFromJson(json['Answers']),
      progressionExtFromJson(json['Progression']));
}

Map<String, dynamic> evaluateWorkInToJson(EvaluateWorkIn item) {
  return {
    "ID": workIDToJson(item.iD),
    "Answers": dictIntToAnswerPToJson(item.answers),
    "Progression": progressionExtToJson(item.progression)
  };
}

// github.com/benoitkugler/maths-online/server/src/tasks.EvaluateWorkOut
class EvaluateWorkOut {
  final Map<int, QuestionAnswersOut> results;
  final ProgressionExt progression;
  final List<InstantiatedQuestion> newQuestions;

  const EvaluateWorkOut(this.results, this.progression, this.newQuestions);

  @override
  String toString() {
    return "EvaluateWorkOut($results, $progression, $newQuestions)";
  }
}

EvaluateWorkOut evaluateWorkOutFromJson(dynamic json_) {
  final json = (json_ as Map<String, dynamic>);
  return EvaluateWorkOut(
      dictIntToQuestionAnswersOutFromJson(json['Results']),
      progressionExtFromJson(json['Progression']),
      listInstantiatedQuestionFromJson(json['NewQuestions']));
}

Map<String, dynamic> evaluateWorkOutToJson(EvaluateWorkOut item) {
  return {
    "Results": dictIntToQuestionAnswersOutToJson(item.results),
    "Progression": progressionExtToJson(item.progression),
    "NewQuestions": listInstantiatedQuestionToJson(item.newQuestions)
  };
}

// github.com/benoitkugler/maths-online/server/src/tasks.InstantiateQuestionsOut
typedef InstantiateQuestionsOut = List<InstantiatedQuestion>;

InstantiateQuestionsOut instantiateQuestionsOutFromJson(dynamic json) {
  return listInstantiatedQuestionFromJson(json);
}

dynamic instantiateQuestionsOutToJson(InstantiateQuestionsOut item) {
  return listInstantiatedQuestionToJson(item);
}

// github.com/benoitkugler/maths-online/server/src/tasks.InstantiatedQuestion
class InstantiatedQuestion {
  final IdQuestion id;
  final Question question;
  final Params params;

  const InstantiatedQuestion(this.id, this.question, this.params);

  @override
  String toString() {
    return "InstantiatedQuestion($id, $question, $params)";
  }
}

InstantiatedQuestion instantiatedQuestionFromJson(dynamic json_) {
  final json = (json_ as Map<String, dynamic>);
  return InstantiatedQuestion(intFromJson(json['Id']),
      questionFromJson(json['Question']), paramsFromJson(json['Params']));
}

Map<String, dynamic> instantiatedQuestionToJson(InstantiatedQuestion item) {
  return {
    "Id": intToJson(item.id),
    "Question": questionToJson(item.question),
    "Params": paramsToJson(item.params)
  };
}

// github.com/benoitkugler/maths-online/server/src/tasks.InstantiatedWork
class InstantiatedWork {
  final WorkID iD;
  final String title;
  final Flow flow;
  final List<InstantiatedQuestion> questions;
  final List<int> baremes;

  const InstantiatedWork(
      this.iD, this.title, this.flow, this.questions, this.baremes);

  @override
  String toString() {
    return "InstantiatedWork($iD, $title, $flow, $questions, $baremes)";
  }
}

InstantiatedWork instantiatedWorkFromJson(dynamic json_) {
  final json = (json_ as Map<String, dynamic>);
  return InstantiatedWork(
      workIDFromJson(json['ID']),
      stringFromJson(json['Title']),
      flowFromJson(json['Flow']),
      listInstantiatedQuestionFromJson(json['Questions']),
      listIntFromJson(json['Baremes']));
}

Map<String, dynamic> instantiatedWorkToJson(InstantiatedWork item) {
  return {
    "ID": workIDToJson(item.iD),
    "Title": stringToJson(item.title),
    "Flow": flowToJson(item.flow),
    "Questions": listInstantiatedQuestionToJson(item.questions),
    "Baremes": listIntToJson(item.baremes)
  };
}

// github.com/benoitkugler/maths-online/server/src/tasks.Params
typedef Params = List<VarEntry>;

Params paramsFromJson(dynamic json) {
  return listVarEntryFromJson(json);
}

dynamic paramsToJson(Params item) {
  return listVarEntryToJson(item);
}

// github.com/benoitkugler/maths-online/server/src/tasks.ProgressionExt
class ProgressionExt {
  final List<QuestionHistory> questions;
  final int nextQuestion;

  const ProgressionExt(this.questions, this.nextQuestion);

  @override
  String toString() {
    return "ProgressionExt($questions, $nextQuestion)";
  }
}

ProgressionExt progressionExtFromJson(dynamic json_) {
  final json = (json_ as Map<String, dynamic>);
  return ProgressionExt(listQuestionHistoryFromJson(json['Questions']),
      intFromJson(json['NextQuestion']));
}

Map<String, dynamic> progressionExtToJson(ProgressionExt item) {
  return {
    "Questions": listQuestionHistoryToJson(item.questions),
    "NextQuestion": intToJson(item.nextQuestion)
  };
}

// github.com/benoitkugler/maths-online/server/src/tasks.TaskProgressionHeader
class TaskProgressionHeader {
  final IdTask id;
  final String title;
  final bool hasProgression;
  final ProgressionExt progression;
  final int mark;
  final int bareme;

  const TaskProgressionHeader(this.id, this.title, this.hasProgression,
      this.progression, this.mark, this.bareme);

  @override
  String toString() {
    return "TaskProgressionHeader($id, $title, $hasProgression, $progression, $mark, $bareme)";
  }
}

TaskProgressionHeader taskProgressionHeaderFromJson(dynamic json_) {
  final json = (json_ as Map<String, dynamic>);
  return TaskProgressionHeader(
      intFromJson(json['Id']),
      stringFromJson(json['Title']),
      boolFromJson(json['HasProgression']),
      progressionExtFromJson(json['Progression']),
      intFromJson(json['Mark']),
      intFromJson(json['Bareme']));
}

Map<String, dynamic> taskProgressionHeaderToJson(TaskProgressionHeader item) {
  return {
    "Id": intToJson(item.id),
    "Title": stringToJson(item.title),
    "HasProgression": boolToJson(item.hasProgression),
    "Progression": progressionExtToJson(item.progression),
    "Mark": intToJson(item.mark),
    "Bareme": intToJson(item.bareme)
  };
}

// github.com/benoitkugler/maths-online/server/src/tasks.VarEntry
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
  final json = (json_ as Map<String, dynamic>);
  return VarEntry(
      variableFromJson(json['Variable']), stringFromJson(json['Resolved']));
}

Map<String, dynamic> varEntryToJson(VarEntry item) {
  return {
    "Variable": variableToJson(item.variable),
    "Resolved": stringToJson(item.resolved)
  };
}

// github.com/benoitkugler/maths-online/server/src/tasks.WorkID
class WorkID {
  final int iD;
  final bool isExercice;

  const WorkID(this.iD, this.isExercice);

  @override
  String toString() {
    return "WorkID($iD, $isExercice)";
  }
}

WorkID workIDFromJson(dynamic json_) {
  final json = (json_ as Map<String, dynamic>);
  return WorkID(intFromJson(json['ID']), boolFromJson(json['IsExercice']));
}

Map<String, dynamic> workIDToJson(WorkID item) {
  return {"ID": intToJson(item.iD), "IsExercice": boolToJson(item.isExercice)};
}

Map<int, AnswerP> dictIntToAnswerPFromJson(dynamic json) {
  if (json == null) {
    return {};
  }
  return (json as Map<String, dynamic>)
      .map((k, v) => MapEntry(int.parse(k), answerPFromJson(v)));
}

Map<String, dynamic> dictIntToAnswerPToJson(Map<int, AnswerP> item) {
  return item
      .map((k, v) => MapEntry(intToJson(k).toString(), answerPToJson(v)));
}

Map<int, QuestionAnswersOut> dictIntToQuestionAnswersOutFromJson(dynamic json) {
  if (json == null) {
    return {};
  }
  return (json as Map<String, dynamic>)
      .map((k, v) => MapEntry(int.parse(k), questionAnswersOutFromJson(v)));
}

Map<String, dynamic> dictIntToQuestionAnswersOutToJson(
    Map<int, QuestionAnswersOut> item) {
  return item.map(
      (k, v) => MapEntry(intToJson(k).toString(), questionAnswersOutToJson(v)));
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

List<int> listIntFromJson(dynamic json) {
  if (json == null) {
    return [];
  }
  return (json as List<dynamic>).map(intFromJson).toList();
}

List<dynamic> listIntToJson(List<int> item) {
  return item.map(intToJson).toList();
}

List<QuestionHistory> listQuestionHistoryFromJson(dynamic json) {
  if (json == null) {
    return [];
  }
  return (json as List<dynamic>).map(questionHistoryFromJson).toList();
}

List<dynamic> listQuestionHistoryToJson(List<QuestionHistory> item) {
  return item.map(questionHistoryToJson).toList();
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