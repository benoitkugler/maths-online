import type { Date_, ExerciceHeader, SheetExt, Time } from "./api_gen";

export function copy<T>(obj: T): T {
  return JSON.parse(JSON.stringify(obj));
}

export function formatDate(date: Date_) {
  const d = new Date(date);
  return d.toLocaleDateString();
}

export function formatTime(time: Time) {
  const ti = new Date(time);
  if (isNaN(ti.valueOf()) || ti.getFullYear() <= 1) {
    return "";
  }
  return ti.toLocaleString(undefined, {
    year: "numeric",
    day: "numeric",
    month: "short",
    hour: "2-digit",
  });
}

export function exerciceBareme(ex: ExerciceHeader) {
  return ex.Questions?.reduce((v, qu) => v + qu.bareme, 0) || 0;
}

export function sheetBareme(sheet: SheetExt) {
  return sheet.Exercices?.reduce((v, ex) => v + exerciceBareme(ex), 0) || 0;
}
