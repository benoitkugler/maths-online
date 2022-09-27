import type { Date_, Monoquestion, SheetExt, TaskExt, Time } from "./api_gen";

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
    hour: "2-digit"
  });
}

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
  return task.Baremes?.reduce((v, qu) => v + qu, 0) || 0;
}

export function sheetBareme(sheet: SheetExt) {
  return sheet.Tasks?.reduce((v, task) => v + taskBareme(task), 0) || 0;
}

export function monoquestionFromTask(task: TaskExt): Monoquestion {
  return {
    Id: task.IdWork.ID,
    IdQuestion: -1,
    Bareme: (task.Baremes || [])[0],
    NbRepeat: task.Baremes?.length || 0
  };
}
