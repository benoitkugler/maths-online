import type { CategoriesQuestions } from "./api_gen";
import { LevelTag } from "./exercice_gen";

export const colorsPerCategorie = [
  "purple",
  "green",
  "orange",
  "#FDD835",
  "blue"
];

/** questionPropositions is a list of question to use in a 
trivial configuration, working nicely with the officialy supported questions
*/
export const questionPropositions: {
  name: string;
  Questions: CategoriesQuestions;
}[] = [
  {
    name: "Pourcentages - 2NDE",
    Questions: {
      Difficulties: [], // all difficulties accepted,
      Tags: [
        [[LevelTag.Seconde, "POURCENTAGES", "EVOLUTION UNIQUE"]],
        [
          [LevelTag.Seconde, "POURCENTAGES", "Taux global"],
          [LevelTag.Seconde, "POURCENTAGES", "Taux réciproque"]
        ],
        [
          [LevelTag.Seconde, "POURCENTAGES", "Proportion"],
          [LevelTag.Seconde, "POURCENTAGES", "Proportion de proportion"],
          [LevelTag.Seconde, "POURCENTAGES", "Pourcentage d'un nombre"]
        ],
        [
          [LevelTag.Seconde, "POURCENTAGES", "Evolutions identiques"],
          [LevelTag.Seconde, "POURCENTAGES", "Evolutions successives"]
        ],
        [
          [LevelTag.Seconde, "POURCENTAGES", "Coefficient multiplicateur"],
          [LevelTag.Seconde, "POURCENTAGES", "Taux d'évolution"]
        ]
      ]
    }
  }
];
