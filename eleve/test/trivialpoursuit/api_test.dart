import 'dart:convert';

import 'package:eleve/trivialpoursuit/events.gen.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

void main() async {
  final channel = WebSocketChannel.connect(
    Uri.parse('ws://localhost:8080/trivial-poursuit'),
  );

  channel.stream.listen((event) {
    try {
      print("as JSON: ${jsonDecode(event as String)}");
    } catch (e) {
      print("as text $event");
    }
  }, onError: (dynamic error) => print("$error"), onDone: () => print("done"));

  channel.sink.add("BAD");

  final jsonMessage =
      jsonEncode(clientEventToJson(const ClientEvent(Ping("info"), 0)));
  channel.sink.add(jsonMessage);

  await Future<void>.delayed(const Duration(seconds: 2));
  channel.sink.close(1000, "Bye");
}
