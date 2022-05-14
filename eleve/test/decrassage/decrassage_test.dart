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
      "Kind": 15
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
      "Indice": "i",
      "Name": 86
     },
     "Resolved": {
      "V": {
       "Indice": "",
       "Name": 0
      },
      "N": 2,
      "IsVariable": false
     }
    },
    {
     "Variable": {
      "Indice": "f",
      "Name": 86
     },
     "Resolved": {
      "V": {
       "Indice": "",
       "Name": 0
      },
      "N": 79,
      "IsVariable": false
     }
    },
    {
     "Variable": {
      "Indice": "",
      "Name": 116
     },
     "Resolved": {
      "V": {
       "Indice": "",
       "Name": 0
      },
      "N": 38.5,
      "IsVariable": false
     }
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
     "Resolved": {
      "V": {
       "Indice": "",
       "Name": 0
      },
      "N": 66,
      "IsVariable": false
     }
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
         "Text": "115 €",
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
         "Text": "97",
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
      "Indice": "i",
      "Name": 86
     },
     "Resolved": {
      "V": {
       "Indice": "",
       "Name": 0
      },
      "N": 115,
      "IsVariable": false
     }
    },
    {
     "Variable": {
      "Indice": "f",
      "Name": 86
     },
     "Resolved": {
      "V": {
       "Indice": "",
       "Name": 0
      },
      "N": 97,
      "IsVariable": false
     }
    },
    {
     "Variable": {
      "Indice": "",
      "Name": 116
     },
     "Resolved": {
      "V": {
       "Indice": "",
       "Name": 0
      },
      "N": -0.1565217391304348,
      "IsVariable": false
     }
    },
    {
     "Variable": {
      "Indice": "",
      "Name": 114
     },
     "Resolved": {
      "V": {
       "Indice": "",
       "Name": 0
      },
      "N": -0.1565,
      "IsVariable": false
     }
    },
    {
     "Variable": {
      "Indice": "",
      "Name": 97
     },
     "Resolved": {
      "V": {
       "Indice": "",
       "Name": 0
      },
      "N": 2021,
      "IsVariable": false
     }
    }
   ]
  }
 ]
 """;
    final questions = listInstantiatedQuestionFromJson(jsonDecode(input));
    expect(questions.length, equals(3));
  });
}
