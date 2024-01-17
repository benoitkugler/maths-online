import {
  PublicStatus,
  type ClassroomTravaux,
  type Date_,
  type Origin,
  type SheetExt,
  type Tags,
  type TaskExt,
  type Time,
  Stage,
  Visibility,
  Rank,
} from "./api_gen";

/** copy returns a deep copy of `obj` */
export function copy<T>(obj: T): T {
  return JSON.parse(JSON.stringify(obj));
}

export function formatDate(date: Date_) {
  const d = new Date(date);
  return d.toLocaleDateString();
}

export function formatTime(
  time: Time,
  showMinute = false,
  showWeekday = false
) {
  const ti = new Date(time);
  if (isNaN(ti.valueOf()) || ti.getFullYear() <= 1) {
    return "";
  }
  const s = ti.toLocaleString(undefined, {
    year: "numeric",
    day: "numeric",
    month: "short",
    hour: "2-digit",
    minute: showMinute ? "2-digit" : undefined,
  });
  if (showWeekday) {
    return `${_weekdays[ti.getDay()]} ${s}`;
  }
  return s;
}

const _weekdays = ["Dim.", "Lun.", "Mar.", "Mer.", "Jeu.", "Ven.", "Sam."];

export function onDragListItemStart(payload: DragEvent, index: number) {
  if (payload.dataTransfer == null) return;
  payload.dataTransfer.setData("text/json", JSON.stringify({ index: index }));
  payload.dataTransfer.dropEffect = "move";
}

/** take the block at the index `origin` and insert it right before
  the block at index `target` (which is between 0 and nbBlocks)
   */
export function swapItems<T>(origin: number, target: number, list: T[]) {
  if (target == origin || target == origin + 1) {
    // nothing to do
    return list;
  }

  if (origin < target) {
    const after = list.slice(target);
    const before = list.slice(0, target);
    const originRow = before.splice(origin, 1);
    before.push(...originRow);
    before.push(...after);
    return before;
  } else {
    const before = list.slice(0, target);
    const originRow = list.splice(origin, 1);
    const after = list.slice(target);
    before.push(...originRow);
    before.push(...after);
    return before;
  }
}

export function taskBareme(task: TaskExt) {
  return task.Bareme?.reduce((v, qu) => v + qu, 0) || 0;
}

export function sheetBareme(sheet: SheetExt) {
  return sheet.Tasks?.reduce((v, task) => v + taskBareme(task), 0) || 0;
}

// safer and easier access
export interface HomeworksT {
  Sheets: Map<number, SheetExt>;
  Travaux: ClassroomTravaux[];
}

export interface PrefillTrivialCategorie {
  matiere: string;
  level: string;
  chapter: string;
  sublevels: Tags;
}

/** `visiblityColors` exposes the colors used to differentiate ressource visiblity */
export const ColorAdmin = "yellow-lighten-4";
export const ColorPersonnal = "light-blue-lighten-5";
export const ColorPublic = "blue-lighten-4";

export function colorForOrigin(origin: Origin) {
  if (origin.PublicStatus == PublicStatus.AdminPublic) return ColorPublic;
  return origin.Visibility == Visibility.Personnal
    ? ColorPersonnal
    : ColorAdmin;
}

export const rankColors: { [key in Rank]: string } = {
  [Rank.StartRank]: "",
  [Rank.Blanche]: "white",
  [Rank.Jaune]: "yellow-accent-2",
  [Rank.Orange]: "orange",
  [Rank.Verte]: "light-green",
  [Rank.Bleue]: "blue",
  [Rank.Rouge]: "red",
  [Rank.Marron]: "deep-orange-darken-4",
  [Rank.Noire]: "black",
};

export function sameStage(loc1: Stage, loc2: Stage) {
  return loc1.Domain == loc2.Domain && loc1.Rank == loc2.Rank;
}
