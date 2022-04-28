import { controller } from "./controller";
import type { teacherSocketData } from "./trivial_config_socket_gen";

export const colorsPerCategorie = [
  "purple",
  "green",
  "orange",
  "yellow",
  "blue"
];

export class TrivialMonitorController {
  constructor(
    sessionID: string,
    onServerEvent: (data: teacherSocketData) => void
  ) {
    const url = controller
      .getURL(`/prof/trivial/monitor?session-id=${sessionID}`)
      .replace("http", "ws");
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
