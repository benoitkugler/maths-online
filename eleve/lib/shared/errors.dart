import 'dart:convert';

import 'package:flutter/material.dart';

dynamic checkServerError(String source) {
  final json = jsonDecode(source);
  if (json is Map<String, dynamic>) {
    if (json.length == 1 && json.containsKey("message")) {
      throw json["message"] as String;
    }
  }
  return json;
}

void showError(String kind, dynamic error, BuildContext context) {
  ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      duration: const Duration(seconds: 6),
      backgroundColor: Theme.of(context).colorScheme.error,
      content: RichText(
          text: TextSpan(children: [
        TextSpan(
            text: kind + "\n",
            style: const TextStyle(fontWeight: FontWeight.bold)),
        const TextSpan(text: "Détails : "),
        TextSpan(
            text: "$error",
            style: const TextStyle(fontStyle: FontStyle.italic)),
      ]))));
}
