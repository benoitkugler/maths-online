import 'dart:convert';

import 'package:eleve/shared_gen.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  testWidgets('decrassage ...', (tester) async {
    const input = """
[
  {
   "Id": 24,
   "Question": {
    "Title": "Calculer un taux d'évolution",
    "Enonce": [
     {
      "Data": {
       "Parts": [
        {
         "Text": "Un article est mis en vente aux enchères au prix de ",
         "IsMath": false
        },
        {
         "Text": "16 €",
         "IsMath": true
        },
        {
         "Text": ". Il est vendu quelques minutes plus tard à ",
         "IsMath": false
        },
        {
         "Text": "66",
         "IsMath": true
        },
        {
         "Text": " ",
         "IsMath": false
        },
        {
         "Text": "€",
         "IsMath": true
        },
        {
         "Text": ".",
         "IsMath": false
        }
       ],
       "Bold": false,
       "Italic": false,
       "Smaller": false
      },
      "Kind": 15
     },
     {
      "Data": {
       "Parts": [
        {
         "Text": "Quel est le pourcentage de la hausse ?+",
         "IsMath": false
        }
       ],
       "Bold": true,
       "Italic": false,
       "Smaller": false
      },
      "Kind": 15
     },
     {
      "Data": {
       "ID": 0
      },
      "Kind": 9
     },
     {
      "Data": {
       "Parts": [
        {
         "Text": "%",
         "IsMath": false
        }
       ],
       "Bold": false,
       "Italic": false,
       "Smaller": false
      },
      "Kind": 15
     },
     {
      "Data": {
       "Parts": [
        {
         "Text": "On donnera la valeur exacte.",
         "IsMath": false
        }
       ],
       "Bold": false,
       "Italic": true,
       "Smaller": true
      },
      "Kind": 15
     }
    ]
   },
   "Params": [
    {
     "Variable": {
      "Indice": "f",
      "Name": 86
     },
     "Resolved": "66"
    },
    {
     "Variable": {
      "Indice": "",
      "Name": 116
     },
     "Resolved": "3,125"
    },
    {
     "Variable": {
      "Indice": "i",
      "Name": 86
     },
     "Resolved": "16"
    }
   ]
  },
  {
   "Id": 29,
   "Question": {
    "Title": "Calculer un coefficient multiplicateur",
    "Enonce": [
     {
      "Data": {
       "Parts": [
        {
         "Text": "Quel est le coefficient multiplicateur associé à une baisse de ",
         "IsMath": false
        },
        {
         "Text": "-68",
         "IsMath": true
        },
        {
         "Text": "% ?",
         "IsMath": false
        },
        {
         "Text": "CM=",
         "IsMath": true
        }
       ],
       "Bold": false,
       "Italic": false,
       "Smaller": false
      },
      "Kind": 15
     },
     {
      "Data": {
       "ID": 0
      },
      "Kind": 9
     }
    ]
   },
   "Params": [
    {
     "Variable": {
      "Indice": "",
      "Name": 116
     },
     "Resolved": "68"
    }
   ]
  },
  {
   "Id": 37,
   "Question": {
    "Title": "Calculer un taux d'évolution",
    "Enonce": [
     {
      "Data": {
       "Parts": [
        {
         "Text": "Un enfant possède ",
         "IsMath": false
        },
        {
         "Text": "99 €",
         "IsMath": true
        },
        {
         "Text": " dans sa tirelire le 1er janvier ",
         "IsMath": false
        },
        {
         "Text": "2018",
         "IsMath": true
        },
        {
         "Text": ".Le 31 décembre de la même année, il possède désormais ",
         "IsMath": false
        },
        {
         "Text": "71",
         "IsMath": true
        },
        {
         "Text": " ",
         "IsMath": false
        },
        {
         "Text": "€",
         "IsMath": true
        },
        {
         "Text": ".",
         "IsMath": false
        }
       ],
       "Bold": false,
       "Italic": false,
       "Smaller": false
      },
      "Kind": 15
     },
     {
      "Data": {
       "Parts": [
        {
         "Text": "Quel est le pourcentage d'évolution de la somme présente dans sa tirelire au cours de cette année écoulée ?",
         "IsMath": false
        }
       ],
       "Bold": true,
       "Italic": false,
       "Smaller": false
      },
      "Kind": 15
     },
     {
      "Data": {
       "ID": 0
      },
      "Kind": 9
     },
     {
      "Data": {
       "Parts": [
        {
         "Text": "%",
         "IsMath": false
        }
       ],
       "Bold": false,
       "Italic": false,
       "Smaller": false
      },
      "Kind": 15
     },
     {
      "Data": {
       "Parts": [
        {
         "Text": "On arrondira à ",
         "IsMath": false
        },
        {
         "Text": "0,01",
         "IsMath": true
        },
        {
         "Text": "% près.",
         "IsMath": false
        }
       ],
       "Bold": false,
       "Italic": true,
       "Smaller": true
      },
      "Kind": 15
     }
    ]
   },
   "Params": [
    {
     "Variable": {
      "Indice": "f",
      "Name": 86
     },
     "Resolved": "71"
    },
    {
     "Variable": {
      "Indice": "",
      "Name": 116
     },
     "Resolved": "-28 / 99"
    },
    {
     "Variable": {
      "Indice": "",
      "Name": 114
     },
     "Resolved": "-0,2828"
    },
    {
     "Variable": {
      "Indice": "",
      "Name": 97
     },
     "Resolved": "2018"
    },
    {
     "Variable": {
      "Indice": "i",
      "Name": 86
     },
     "Resolved": "99"
    }
   ]
  }
 ]
 """;
    final questions = listInstantiatedQuestionFromJson(jsonDecode(input));
    expect(questions.length, equals(3));
  });
}
