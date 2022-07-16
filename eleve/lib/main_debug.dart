import 'package:eleve/main_shared.dart';
import 'package:eleve/questions/proof.dart';
import 'package:eleve/questions/proof.gen.dart';
import 'package:eleve/questions/types.gen.dart';
import 'package:flutter/material.dart';

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
      home: Scaffold(
        body: ListView(children: [
          ProofField(
              Colors.purple,
              ProofController(
                  const ProofFieldBlock(
                      Proof(Sequence([
                        Statement(""),
                        Equality(["", "", "", "", ""]),
                        Node(
                            Statement(""),
                            Sequence([
                              Statement(""),
                              Equality(["", "", "", "", ""]),
                            ]),
                            Binary.invalid),
                        Sequence([
                          Statement(""),
                          Statement(""),
                          Statement(""),
                        ])
                      ])),
                      ["test", "a", "b"],
                      0),
                  () => print("ok"))),
        ]),
      ),
    );
  }
}
