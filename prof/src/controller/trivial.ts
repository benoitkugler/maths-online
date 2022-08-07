import { CategoriesQuestions } from "./api_gen";
import { controller } from "./controller";
import type { teacherSocketData } from "./trivial_config_socket_gen";

export const colorsPerCategorie = [
  "purple",
  "green",
  "orange",
  "yellow",
  "blue",
];

export class TrivialMonitorController {
  constructor(onServerEvent: (data: teacherSocketData) => void) {
    const url =
      controller.getURL(`/prof/trivial/monitor`).replace("http", "ws") +
      "?token=" +
      controller.getToken();
    const socket = new WebSocket(url);

    // Connection opened
    socket.addEventListener("open", function (event) {
      socket.send("Hello Server!");
    });

    // Listen for messages
    socket.addEventListener("message", function (event) {
      onServerEvent(JSON.parse(event.data));
    });
  }
}

/** questionPropositions is a list of question to use in a 
trivial configuration, working nicely with the officialy supported questions
*/
export const questionPropositions: {
  name: string;
  Questions: CategoriesQuestions;
}[] = [
  {
    name: "Pourcentages",
    Questions: [
      [
        ["Pourcentages", "Valeur initiale"],
        ["Pourcentages", "Valeur finale"],
      ],
      [
        ["Pourcentages", "Taux global"],
        ["Pourcentages", "Taux réciproque"],
      ],
      [
        ["Pourcentages", "Proportion"],
        ["Pourcentages", "Proportion de proportion"],
        ["Pourcentages", "Pourcentage d'un nombre"],
      ],
      [
        ["Pourcentages", "Evolutions identiques"],
        ["Pourcentages", "Evolution unique"],
        ["Pourcentages", "Evolutions successives"],
      ],
      [
        ["Pourcentages", "Coefficient multiplicateur"],
        ["Pourcentages", "Taux d'évolution"],
      ],
    ],
  },
];
