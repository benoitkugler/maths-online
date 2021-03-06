// Code generated by structgen DO NOT EDIT
// github.com/benoitkugler/maths-online/trivial-poursuit/game.Success
export type Success = boolean[];
// github.com/benoitkugler/maths-online/prof/trivial.gamePlayers
export interface gamePlayers {
  Player: string;
  Successes: Success;
}
// github.com/benoitkugler/maths-online/trivial-poursuit.GameID
export type GameID = string;
// github.com/benoitkugler/maths-online/prof/trivial.gameSummary
export interface gameSummary {
  GameID: GameID;
  CurrentPlayer: string;
  Players: gamePlayers[] | null;
  RoomSize: number;
}
// github.com/benoitkugler/maths-online/prof/trivial.teacherSocketData
export interface teacherSocketData {
  Games: gameSummary[] | null;
}
