import 'package:eleve/main_shared.dart';
import 'package:flutter/material.dart' hide Flow;

void main() async {
  runApp(const _App());
}

class _App extends StatelessWidget {
  const _App({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Isyro',
      theme: theme,
      debugShowCheckedModeBanner: false,
      localizationsDelegates: localizations,
      supportedLocales: locales,
      home: Container(),
    );
  }
}
