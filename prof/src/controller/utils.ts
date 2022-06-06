import type { Date_ } from "./api_gen";

export function copy<T>(obj: T): T {
  return JSON.parse(JSON.stringify(obj));
}

export function formatDate(date: Date_) {
  const d = new Date(date);
  return d.toLocaleDateString();
}
