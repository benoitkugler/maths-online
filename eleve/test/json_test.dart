import 'dart:developer';

import "package:eleve/maths/maths.dart";
import 'package:test/test.dart';

void main() {
  test("JSON", () {
    final l = Obj([
      FormulaAnswer(),
      TextBlock("Essai"),
      InlineFormula(r"\frac{a}{b}"),
      DisplayFormula("x + y"),
      ListAnswser(),
    ]);
    final m = l.toJson();
    print(m);

    final l2 = Obj.fromJson(m);
    print(l2.l);
    l2.l.forEach(inspect);

    expect(l2, equals(l));
  });
  // testWidgets('Counter increments smoke test', (WidgetTester tester) async {
  //   // Build our app and trigger a frame.
  //   await tester.pumpWidget(const MyApp());

  //   // Verify that our counter starts at 0.
  //   expect(find.text('0'), findsOneWidget);
  //   expect(find.text('1'), findsNothing);

  //   // Tap the '+' icon and trigger a frame.
  //   await tester.tap(find.byIcon(Icons.add));
  //   await tester.pump();

  //   // Verify that our counter has incremented.
  //   expect(find.text('0'), findsNothing);
  //   expect(find.text('1'), findsOneWidget);

  //   // Tap the '+' icon and trigger a frame.
  //   await tester.tap(find.byIcon(Icons.add));
  //   await tester.pump();

  //   // Verify that our counter has incremented.
  //   expect(find.text('0'), findsNothing);
  //   expect(find.text('1'), findsNothing);
  //   expect(find.text('2'), findsOneWidget);
  // });
}
