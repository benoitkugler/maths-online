import 'package:eleve/main_shared.dart';
import 'package:eleve/questions/debug.dart';
import 'package:eleve/questions/fields.dart';
import 'package:material_ui/material_ui.dart';

void main() async {
  runApp(const _URLApp());
}

class _URLApp extends StatelessWidget {
  const _URLApp();

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Isyro',
      theme: theme,
      debugShowCheckedModeBanner: false,
      localizationsDelegates: localizations,
      supportedLocales: locales,
      home: Builder(
        builder: (context) => Scaffold(
          body: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              TextRow(
                buildText(
                  [
                    T(
                      "une url https://youtube.com/shorts/Cb-tYUzhRew?si=x5bTOxSMfju3VuUj",
                    ),
                  ],
                  TextS(),
                  12,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
