import 'package:material_ui/material_ui.dart';

const Color darkBlue = Color(0xFF104986);

final theme = ThemeData.dark().copyWith(
  scaffoldBackgroundColor: darkBlue,
  cardTheme: ThemeData.dark().cardTheme.copyWith(
    color: darkBlue.withValues(alpha: 0.7),
  ),
);

const localizations = GlobalMaterialLocalizations.delegates;

const locales = [
  Locale('fr', ''), // French, no country code
  Locale('en', ''), // English, no country code
];
