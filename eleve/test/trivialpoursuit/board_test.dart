import 'package:eleve/activities/trivialpoursuit/board.dart';
import 'package:eleve/activities/trivialpoursuit/categories.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  testWidgets('board ...', (tester) async {
    expect(Board.shapes.length, equals(19));
    expect(Board.shapes[0].categorie.color, equals(Colors.purple));
    expect(Board.shapes[16].categorie.color, equals(Colors.blue));
  });
}
