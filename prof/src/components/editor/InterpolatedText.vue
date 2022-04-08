<template>
  <QuillEditor
    theme=""
    toolbar=""
    class="text-field"
    content-type="text"
    @update:content="onTextChange"
    @text-change="colorize"
    :content="props.modelValue"
    ref="quill"
  />
</template>

<script setup lang="ts">
import { itemize } from "@/controller/editor";
import { TextKind } from "@/controller/exercice_gen";
import { onMounted, watch } from "@vue/runtime-core";
import type { Quill } from "@vueup/vue-quill";
import { QuillEditor } from "@vueup/vue-quill";
import "@vueup/vue-quill/dist/vue-quill.snow.css";
import type { Sources } from "quill";
import { $ref } from "vue/macros";

type Props = {
  modelValue: string;
};

const props = defineProps<Props>();
const emit = defineEmits<{
  (event: "update:modelValue", value: string): void;
}>();

const colorByKind = ["", "green", "orange"];

let quill = $ref<InstanceType<typeof QuillEditor> | null>();

watch(props, () => {
  const current = quill?.getText().trimRight(); // quill add a `\n`
  if (current != props.modelValue) {
    quill?.setText(props.modelValue);
  }
  colorize({ source: "user" });
});

onMounted(() => colorize({ source: "user" }));

function onTextChange(text: string) {
  emit("update:modelValue", text.trimRight());
}

function colorize(arg: { source: Sources }) {
  if (arg.source != "user" || quill == null) {
    return;
  }
  const text = quill?.getText() || "";
  const qu = quill?.getQuill() as Quill;
  const parts = itemize(text);
  let cursor = 0;
  parts.forEach(p => {
    qu.formatText(cursor, p.Content.length, {
      color: colorByKind[p.Kind],
      bold: p.Kind == TextKind.Expression
    });
    cursor += p.Content.length;
  });
}
</script>

<style>
.text-field {
  width: 100%;
  border: 1px solid black;
  border-radius: 5px;
  background-color: #d4d4d4a8;
}

.ql-editor {
  padding: 4px 2px;
}
</style>
