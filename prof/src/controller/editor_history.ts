import { copy } from "./utils";

export class History<T> {
  // start with initial question in history
  constructor(
    initialSnapshot: T,
    private showMessage: (m: string, color: string) => void,
    /** `restore` is called when a snapshot should be
     * restored
     */
    private restore: (toRestore: T) => void
  ) {
    this.data = [copy(initialSnapshot)];
  }

  // the accumulated snapshot
  private data: T[];
  private index = 0;

  /** `add` adds a copy of `snapshot` in the history list, at the current index */
  add(snapshot: T) {
    // if we are back in time, erase future versions
    this.data.splice(this.index + 1);
    this.data.push(copy(snapshot));
    this.index++;
  }

  private back() {
    if (this.index >= 1) {
      this.index--;
      this.restore(copy(this.data[this.index]));
    } else {
      this.showMessage("L'historique est déjà au point initial.", "warning");
    }
  }

  private next() {
    if (this.index < this.data.length - 1) {
      this.index++;
      this.restore(copy(this.data[this.index]));
    } else {
      this.showMessage("L'historique est déjà au point final.", "warning");
    }
  }

  addListener() {
    document.addEventListener("keydown", this.onKeyDown);
  }

  /** `clearListener` must be called to ensure the event listener is cleared */
  clearListener() {
    document.removeEventListener("keydown", this.onKeyDown);
  }

  private onKeyDown = (e: KeyboardEvent) => {
    if (e.key === "z" && (e.ctrlKey || e.metaKey)) {
      e.preventDefault();
      this.back();
    } else if (e.key === "y" && (e.ctrlKey || e.metaKey)) {
      e.preventDefault();
      this.next();
    }
  };
}
