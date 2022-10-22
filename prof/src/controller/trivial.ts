import type { CategoriesQuestions } from "./api_gen";

export const colorsPerCategorie = [
  "purple",
  "green",
  "orange",
  "yellow-darken-2",
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
    name: "Pourcentages",
    Questions: {
      Difficulties: [], // all difficulties accepted,
      Tags: [
        [["Pourcentages", "EVOLUTION UNIQUE"]],
        [
          ["Pourcentages", "Taux global"],
          ["Pourcentages", "Taux réciproque"]
        ],
        [
          ["Pourcentages", "Proportion"],
          ["Pourcentages", "Proportion de proportion"],
          ["Pourcentages", "Pourcentage d'un nombre"]
        ],
        [
          ["Pourcentages", "Evolutions identiques"],
          ["Pourcentages", "Evolutions successives"]
        ],
        [
          ["Pourcentages", "Coefficient multiplicateur"],
          ["Pourcentages", "Taux d'évolution"]
        ]
      ]
    }
  }
];
