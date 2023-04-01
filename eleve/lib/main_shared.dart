import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';

const Color darkBlue = Color(0xFF104986);

final theme = ThemeData.dark().copyWith(
  scaffoldBackgroundColor: darkBlue,
  cardTheme:
      ThemeData.dark().cardTheme.copyWith(color: darkBlue.withOpacity(0.7)),
);

const localizations = [
  GlobalMaterialLocalizations.delegate,
  GlobalWidgetsLocalizations.delegate,
  GlobalCupertinoLocalizations.delegate,
];

const locales = [
  Locale('fr', ''), // French, no country code
  Locale('en', ''), // English, no country code
];
