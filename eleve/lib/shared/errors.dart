import 'dart:convert';

import 'package:flutter/material.dart';

/// checkServerError throws if the serveur returns an error message
/// It should be called on the response body, before deserializing,
/// as a replacement for [jsonDecode]
dynamic checkServerError(String source) {
  final json = jsonDecode(source);
  if (json is Map<String, dynamic>) {
    if (json.length == 1 && json.containsKey("message")) {
      throw json["message"] as String;
    }
  }
  return json;
}

class ErrorBar extends StatelessWidget {
  final String topic;
  final dynamic error;

  const ErrorBar(this.topic, this.error, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return RichText(
        text: TextSpan(children: [
      TextSpan(
          text: "$topic \n",
          style: const TextStyle(fontWeight: FontWeight.bold)),
      const TextSpan(text: "Détails : "),
      TextSpan(
          text: "$error", style: const TextStyle(fontStyle: FontStyle.italic)),
    ]));
  }
}

void showError(String kind, dynamic error, BuildContext context) {
  ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      duration: const Duration(seconds: 6),
      backgroundColor: Theme.of(context).colorScheme.error,
      content: ErrorBar(kind, error)));
}

class ErrorCard extends StatelessWidget {
  final String message;
  final dynamic error;
  const ErrorCard(this.message, this.error, {super.key});

  @override
  Widget build(BuildContext context) {
    return Center(
        child: Card(
            child: Padding(
      padding: const EdgeInsets.all(8.0),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(
            message,
            style: Theme.of(context).textTheme.titleMedium,
          ),
          const SizedBox(height: 20),
          Text("$error", style: const TextStyle(fontStyle: FontStyle.italic)),
        ],
      ),
    )));
  }
}
