import 'package:eleve/trivialpoursuit/game.dart';
import 'package:flutter/material.dart';
import 'package:otp_text_field/otp_text_field.dart';
import 'package:otp_text_field/style.dart';

/// Loggin is an introduction screen to access
/// a TrivialPoursuit game
class TrivialPoursuitLoggin extends StatefulWidget {
  const TrivialPoursuitLoggin({Key? key}) : super(key: key);

  @override
  _TrivialPoursuitLogginState createState() => _TrivialPoursuitLogginState();
}

class _TrivialPoursuitLogginState extends State<TrivialPoursuitLoggin> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Card(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Padding(
              padding: EdgeInsets.symmetric(vertical: 20),
              child: Text(
                "Code de la partie",
                style: TextStyle(fontSize: 20),
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 20),
              child: OTPTextField(
                length: 6,
                width: MediaQuery.of(context).size.width,
                textFieldAlignment: MainAxisAlignment.center,
                otpFieldStyle: OtpFieldStyle(
                    enabledBorderColor:
                        Theme.of(context).colorScheme.secondary),
                fieldStyle: FieldStyle.box,
                fieldWidth: 45,
                onCompleted: _launchTrivialPoursuit,
                onChanged: (_) {},
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _launchTrivialPoursuit(String code) {
    const host = "localhost:1323";
    final uri = Uri.parse('ws://$host/trivial/game/$code');

    Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (_) => Scaffold(body: TrivialPoursuitController(60, uri))));
  }
}
