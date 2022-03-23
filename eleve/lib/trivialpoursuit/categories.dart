import 'package:flutter/material.dart';

import 'events.gen.dart';

/// map question categories to colors
extension CategorieColor on Categorie {
  Color get color {
    switch (this) {
      case Categorie.purple:
        return Colors.purple;
      case Categorie.green:
        return Colors.green;
      case Categorie.orange:
        return Colors.orange.shade700;
      case Categorie.yellow:
        return Colors.yellow.shade700;
      case Categorie.blue:
        return Colors.blue;
    }
  }
}
