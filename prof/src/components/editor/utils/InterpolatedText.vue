<template>
  <div class="mx-0">
    <small v-if="props.label" class="ml-2 text-grey">{{ props.label }} </small>
    <div
      class="editor elevation-2 mb-2"
      contenteditable="true"
      spellcheck="false"
      ref="editor"
      @keydown="handleTab"
      @keyup="onKeyUp"
      :style="props.center ? 'text-align:  center' : ''"
    ></div>
    <small v-if="props.hint" class="text-grey">{{ props.hint }} </small>
  </div>
</template>

<script setup lang="ts">
import { colorByKind } from "@/controller/editor";
import { TextKind } from "@/controller/loopback_gen";
import { computed, onMounted, watch } from "vue";
import { $ref } from "vue/macros";
import { defautTokenize, splitNewLines, type Token } from "./interpolated_text";

type Props = {
  modelValue: string;
  label?: string;
  hint?: string;
  forceLatex?: boolean;
  center?: boolean;
  customTokenize?: (input: string) => Token[];
};

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:modelValue", modelValue: string): void;
}>();

let editor = $ref<HTMLDivElement | null>(null);

// since vue may reuse the same component but change
// the modelValue, we have to watch for it
watch(props, () => {
  const currentText = HTMLToText();
  if (props.modelValue != currentText) {
    updateDisplay(props.modelValue);
  }
});

function textToHTML(input: string) {
  // to make the line actually show and take space, we insert a line break
  if (!input.length) return "<div><br/></div>";
  const tokenize = props.customTokenize ? props.customTokenize : defautTokenize;
  const tokens = tokenize(input);
  const withNewLines = splitNewLines(tokens);

  const rows = withNewLines
    .map((tokens) => {
      const fullText = tokens.map((tok) => tok.Content).join("");
      if (!fullText.length) {
        // avoid empty lines which are not show
        return "<br/>";
      }
      return tokens
        .map((token) => `<span style="${token.Kind}">${token.Content}</span>`)
        .join("");
    })
    .map((line) => `<div>${line}</div>`);

  return rows.join("");
}

function HTMLToText() {
  if (editor == null) return "";
  const rows = Array.from(editor.children);
  const text = rows
    .map((row) =>
      // empty div have innerText as "\n"
      (row as HTMLElement).innerText == "\n"
        ? ""
        : (row as HTMLElement).innerText
    )
    .join("\n");
  return text;
}

function updateText() {
  const text = HTMLToText();
  emit("update:modelValue", text);
  return text;
}

function updateDisplay(source: string) {
  if (editor == null) return;
  editor.innerHTML = textToHTML(source);
}

function caret() {
  if (editor == null) return 0;
  const range = window.getSelection()?.getRangeAt(0);
  if (range == undefined) return 0;
  const prefix = range.cloneRange();
  prefix.selectNodeContents(editor);
  prefix.setEnd(range.endContainer, range.endOffset);
  return prefix.toString().length;
}

function setCaret(pos: number, parent: HTMLElement) {
  for (const node of Array.from(parent.childNodes)) {
    if (node.nodeType == Node.TEXT_NODE) {
      const text = node.nodeValue!;
      if (text.length >= pos) {
        const range = document.createRange();
        const sel = window.getSelection()!;
        range.setStart(node, pos);
        range.collapse(true);
        sel.removeAllRanges();
        sel.addRange(range);
        return -1;
      } else {
        pos = pos - text.length;
      }
    } else {
      pos = setCaret(pos, node as HTMLElement);
      if (pos < 0) {
        return pos;
      }
    }
  }
  return pos;
}

const tab = "    ";
function handleTab(e: KeyboardEvent) {
  if (editor == null) return;
  if (e.key === "Tab") {
    const pos = caret() + tab.length;
    const range = window.getSelection()!.getRangeAt(0);
    range.deleteContents();
    range.insertNode(document.createTextNode(tab));
    const text = updateText();
    updateDisplay(text);
    setCaret(pos, editor);
    e.preventDefault();
  }
}

function onKeyUp(e: KeyboardEvent) {
  if (editor == null) return;
  const text = updateText();
  if (e.keyCode >= 0x30 || e.keyCode == 0x20) {
    const pos = caret();
    updateDisplay(text);
    setCaret(pos, editor);
  }
}

onMounted(() => {
  updateDisplay(props.modelValue);
});

const colorLatex = colorByKind[TextKind.StaticMath];
const activeColor = computed(() => (props.forceLatex ? colorLatex : "#444444"));
</script>

<style>
.editor {
  font-family: "Roboto Mono", monospace;
  font-size: 13px;
  outline: none;
  overflow-y: auto;
  counter-reset: line;

  padding: 4px;
  border-radius: 4px;
  border: 2px solid lightgray;
}

.editor:hover {
  border: 2px solid gray;
}

.editor:focus-within {
  border: 2px solid v-bind("activeColor");
}

.editor div {
  display: block;
  position: relative;
  white-space: pre-wrap;
}
</style>
