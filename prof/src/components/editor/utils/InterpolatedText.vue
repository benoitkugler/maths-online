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
import { computed, onMounted } from "vue";
import { $ref } from "vue/macros";
import { defautTokenize, type Token } from "./interpolated_text";

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

function textToHTML(input: string) {
  // to make the line actually show and take space, we insert an invisble char
  if (!input.length) return "<div>&#8203</div>";
  const tokenize = props.customTokenize ? props.customTokenize : defautTokenize;
  return tokenize(input)
    .map((part) => `<span style="${part.Kind}">${part.Content}</span>`)
    .join("")
    .split("\n")
    .map((line) => `<div>${line}</div>`)
    .join("");
}

function HTMLToText() {
  if (editor == null) return "";
  const text = Array.from(editor.children)
    .map((row) => (row as HTMLElement).innerText)
    .join("\n");

  return text;
}

function updateText() {
  const text = HTMLToText();
  // remove the invisible white space inserted
  if (text.charCodeAt(0) == 8203) {
    emit("update:modelValue", text.substring(1));
  } else {
    emit("update:modelValue", text);
  }
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
