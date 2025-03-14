// Code generated by gomacro/generator/dart. DO NOT EDIT

import 'predefined.dart';

// github.com/benoitkugler/maths-online/server/src/sql/events.EventK
enum EventK {
  isyTriv_Create,
  isyTriv_Streak3,
  isyTriv_Win,
  homework_TaskDone,
  homework_TravailDone,
  all_QuestionRight,
  all_QuestionWrong,
  misc_SetPlaylist,
  connectStreak3,
  connectStreak7,
  connectStreak30
}

extension _EventKExt on EventK {
  static EventK fromValue(int i) {
    return EventK.values[i];
  }

  int toValue() {
    return index;
  }
}

String eventKLabel(EventK v) {
  switch (v) {
    case EventK.isyTriv_Create:
      return "Créer une partie d'IsyTriv";
    case EventK.isyTriv_Streak3:
      return "Réussir trois questions IsyTriv d'affilée";
    case EventK.isyTriv_Win:
      return "Remporter une partie IsyTriv";
    case EventK.homework_TaskDone:
      return "Terminer un exercice";
    case EventK.homework_TravailDone:
      return "Terminer une feuille d'exercices";
    case EventK.all_QuestionRight:
      return "Répondre correctement à une question";
    case EventK.all_QuestionWrong:
      return "Répondre incorrectement à une question";
    case EventK.misc_SetPlaylist:
      return "Modifier sa playlist";
    case EventK.connectStreak3:
      return "Se connecter 3 jours de suite";
    case EventK.connectStreak7:
      return "Se connecter 7 jours de suite";
    case EventK.connectStreak30:
      return "Se connecter 30 jours de suite";
  }
}

EventK eventKFromJson(dynamic json) => _EventKExt.fromValue(json as int);

dynamic eventKToJson(EventK item) => item.toValue();

// github.com/benoitkugler/maths-online/server/src/sql/events.EventNotification
class EventNotification {
  final List<EventK> events;
  final int points;

  const EventNotification(this.events, this.points);

  @override
  String toString() {
    return "EventNotification($events, $points)";
  }
}

EventNotification eventNotificationFromJson(dynamic json_) {
  final json = (json_ as Map<String, dynamic>);
  return EventNotification(
      listEventKFromJson(json['Events']), intFromJson(json['Points']));
}

Map<String, dynamic> eventNotificationToJson(EventNotification item) {
  return {
    "Events": listEventKToJson(item.events),
    "Points": intToJson(item.points)
  };
}

// github.com/benoitkugler/maths-online/server/src/sql/events.StudentAdvance
class StudentAdvance {
  final List<int> occurences;
  final int totalPoints;
  final int flames;
  final int rank;
  final int pointsCurrentRank;
  final int pointsNextRank;

  const StudentAdvance(this.occurences, this.totalPoints, this.flames,
      this.rank, this.pointsCurrentRank, this.pointsNextRank);

  @override
  String toString() {
    return "StudentAdvance($occurences, $totalPoints, $flames, $rank, $pointsCurrentRank, $pointsNextRank)";
  }
}

StudentAdvance studentAdvanceFromJson(dynamic json_) {
  final json = (json_ as Map<String, dynamic>);
  return StudentAdvance(
      listIntFromJson(json['Occurences']),
      intFromJson(json['TotalPoints']),
      intFromJson(json['Flames']),
      intFromJson(json['Rank']),
      intFromJson(json['PointsCurrentRank']),
      intFromJson(json['PointsNextRank']));
}

Map<String, dynamic> studentAdvanceToJson(StudentAdvance item) {
  return {
    "Occurences": listIntToJson(item.occurences),
    "TotalPoints": intToJson(item.totalPoints),
    "Flames": intToJson(item.flames),
    "Rank": intToJson(item.rank),
    "PointsCurrentRank": intToJson(item.pointsCurrentRank),
    "PointsNextRank": intToJson(item.pointsNextRank)
  };
}

List<EventK> listEventKFromJson(dynamic json) {
  if (json == null) {
    return [];
  }
  return (json as List<dynamic>).map(eventKFromJson).toList();
}

List<dynamic> listEventKToJson(List<EventK> item) {
  return item.map(eventKToJson).toList();
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
