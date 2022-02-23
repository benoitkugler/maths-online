// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'maths.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Obj _$ObjFromJson(Map<String, dynamic> json) => Obj(
      (json['l'] as List<dynamic>)
          .map((e) => Part.fromJson(e as Map<String, dynamic>))
          .toList(),
    );

Map<String, dynamic> _$ObjToJson(Obj instance) => <String, dynamic>{
      'l': instance.l.map((e) => e.toJson()).toList(),
    };

TextBlock _$TextBlockFromJson(Map<String, dynamic> json) => TextBlock(
      json['text'] as String,
    );

Map<String, dynamic> _$TextBlockToJson(TextBlock instance) => <String, dynamic>{
      'text': instance.text,
    };

InlineFormula _$InlineFormulaFromJson(Map<String, dynamic> json) =>
    InlineFormula(
      json['latex'] as String,
    );

Map<String, dynamic> _$InlineFormulaToJson(InlineFormula instance) =>
    <String, dynamic>{
      'latex': instance.latex,
    };

DisplayFormula _$DisplayFormulaFromJson(Map<String, dynamic> json) =>
    DisplayFormula(
      json['latex'] as String,
    );

Map<String, dynamic> _$DisplayFormulaToJson(DisplayFormula instance) =>
    <String, dynamic>{
      'latex': instance.latex,
    };

FormulaAnswer _$FormulaAnswerFromJson(Map<String, dynamic> json) =>
    FormulaAnswer();

Map<String, dynamic> _$FormulaAnswerToJson(FormulaAnswer instance) =>
    <String, dynamic>{};

ListAnswser _$ListAnswserFromJson(Map<String, dynamic> json) => ListAnswser();

Map<String, dynamic> _$ListAnswserToJson(ListAnswser instance) =>
    <String, dynamic>{};
