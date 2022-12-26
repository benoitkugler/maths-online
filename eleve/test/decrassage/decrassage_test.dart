import 'dart:convert';

import 'package:eleve/types/src_tasks.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  testWidgets('decrassage ...', (tester) async {
    const input = """
[
  {
   "Id": 24,
   "Question": {
    "Enonce": [
     {
      "Data": {
       "Parts": [
        {
         "Text": "Un article est mis en vente aux enchères au prix de ",
         "IsMath": false
        },
        {
         "Text": "2 €",
         "IsMath": true
        },
        {
         "Text": ". Il est vendu quelques minutes plus tard à ",
         "IsMath": false
        },
        {
         "Text": "79",
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
      "Kind": "TextBlock"
     },
     {
      "Data": {
       "Parts": [
        {
         "Text": "Quel est le pourcentage de la hausse ?\\n+",
         "IsMath": false
        }
       ],
       "Bold": true,
       "Italic": false,
       "Smaller": false
      },
      "Kind": "TextBlock"
     },
     {
      "Data": {
       "ID": 0,
       "SizeHint": 4
      },
      "Kind": "NumberFieldBlock"
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
      "Kind": "TextBlock"
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
      "Kind": "TextBlock"
     }
    ]
   },
   "Params": [
    {
     "Variable": {
      "Indice": "f",
      "Name": 86
     },
     "Resolved": "79"
    },
    {
     "Variable": {
      "Indice": "",
      "Name": 116
     },
     "Resolved": "38,5"
    },
    {
     "Variable": {
      "Indice": "i",
      "Name": 86
     },
     "Resolved": "2"
    }
   ]
  },
  {
   "Id": 29,
   "Question": {
    "Enonce": [
     {
      "Data": {
       "Parts": [
        {
         "Text": "Quel est le coefficient multiplicateur associé à une baisse de ",
         "IsMath": false
        },
        {
         "Text": "-66",
         "IsMath": true
        },
        {
         "Text": "% ?\\n\\n",
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
      "Kind": "TextBlock"
     },
     {
      "Data": {
       "ID": 0,
       "SizeHint": 4
      },
      "Kind": "NumberFieldBlock"
     }
    ]
   },
   "Params": [
    {
     "Variable": {
      "Indice": "",
      "Name": 116
     },
     "Resolved": "66"
    }
   ]
  },
  {
   "Id": 37,
   "Question": {
    "Enonce": [
     {
      "Data": {
       "Parts": [
        {
         "Text": "Un enfant possède ",
         "IsMath": false
        },
        {
         "Text": "97 €",
         "IsMath": true
        },
        {
         "Text": " dans sa tirelire le 1er janvier ",
         "IsMath": false
        },
        {
         "Text": "2021",
         "IsMath": true
        },
        {
         "Text": ".\\nLe 31 décembre de la même année, il possède désormais ",
         "IsMath": false
        },
        {
         "Text": "115",
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
      "Kind": "TextBlock"
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
      "Kind": "TextBlock"
     },
     {
      "Data": {
       "ID": 0,
       "SizeHint": 5
      },
      "Kind": "NumberFieldBlock"
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
      "Kind": "TextBlock"
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
      "Kind": "TextBlock"
     }
    ]
   },
   "Params": [
    {
     "Variable": {
      "Indice": "",
      "Name": 114
     },
     "Resolved": "0,1856"
    },
    {
     "Variable": {
      "Indice": "",
      "Name": 97
     },
     "Resolved": "2021"
    },
    {
     "Variable": {
      "Indice": "i",
      "Name": 86
     },
     "Resolved": "97"
    },
    {
     "Variable": {
      "Indice": "f",
      "Name": 86
     },
     "Resolved": "115"
    },
    {
     "Variable": {
      "Indice": "",
      "Name": 116
     },
     "Resolved": "0,1855670103"
    }
   ]
  }
 ]
 """;
    final questions = listInstantiatedQuestionFromJson(jsonDecode(input));
    expect(questions.length, equals(3));
  });
}
