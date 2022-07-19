<template>
  <small v-if="props.label" class="text-grey">{{ props.label }} </small>
  <QuillEditor
    theme=""
    toolbar=""
    class="__quill-text-field elevation-2"
    content-type="text"
    @update:content="onTextChange"
    @text-change="updateVisual"
    :content="props.modelValue"
    :options="{ formats: ['color', 'bold', 'align'] }"
    ref="quill"
  />
  <small v-if="props.hint" class="text-grey">{{ props.hint }} </small>
</template>

<script setup lang="ts">
import { TextKind } from "@/controller/api_gen";
import { colorByKind, itemize } from "@/controller/editor";
import { computed, onMounted, watch } from "@vue/runtime-core";
import type { Quill } from "@vueup/vue-quill";
import { QuillEditor } from "@vueup/vue-quill";
import "@vueup/vue-quill/dist/vue-quill.snow.css";
import type { Sources } from "quill";
import { $ref } from "vue/macros";

type Props = {
  modelValue: string;
  label?: string;
  hint?: string;
  forceLatex?: boolean;
  transform?: (quill: Quill) => void;
  center?: boolean;
};

const props = defineProps<Props>();
const emit = defineEmits<{
  (event: "update:modelValue", value: string): void;
}>();

let quill = $ref<InstanceType<typeof QuillEditor> | null>();

watch(props, () => {
  const current = quill?.getText().trimEnd(); // quill add a `\n`
  if (current != props.modelValue) {
    quill?.setText(props.modelValue);
  }
  updateVisual({ source: "user" });
});

onMounted(() => setTimeout(() => updateVisual({ source: "user" }), 100));

function onTextChange(text: string) {
  emit("update:modelValue", text.trimEnd());
}

function updateVisual(arg: { source: Sources }) {
  if (arg.source != "user" || quill == null) {
    return;
  }
  const qu = quill?.getQuill() as Quill;

  defaultTransform(qu);
  if (props.transform) {
    props.transform(qu);
  }
}

// colorize $ $ and & &
function defaultTransform(quill: Quill) {
  const text = quill.getText() || "";
  const parts = itemize(text);
  let cursor = 0;
  parts.forEach((p) => {
    quill.formatText(cursor, p.Content.length, {
      color: colorByKind[p.Kind],
      bold: p.Kind == TextKind.Expression,
    });
    cursor += p.Content.length;
  });

  if (props.center) {
    const lineNb = quill.getLines().length;
    quill.formatLine(0, lineNb, { align: "center" });
  }
}

const laTeXColor = colorByKind[TextKind.StaticMath];
const activeColor = computed(() => (props.forceLatex ? laTeXColor : "black"));
</script>

<style>
.__quill-text-field {
  width: 100%;
  border: 2px solid lightgray;
  border-radius: 4px;
  background-color: white;
  padding: 4px;
}

.__quill-text-field:hover {
  border: 2px solid gray;
}

.__quill-text-field:focus-within {
  border: 2px solid v-bind("activeColor");
}

.ql-container {
  height: unset;
}
.ql-editor {
  padding: 4px 4px;
}
</style>
