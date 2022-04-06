<template>
  <container @delete="emit('delete')">
    <template v-slot:toolbar>
      <v-checkbox
        class="pr-2"
        hide-details
        label="Comme conseil"
        v-model="props.Data.IsHint"
      ></v-checkbox>
    </template>
    <v-textarea
      class="px-2"
      dense
      hint="Insérer du code LaTeX avec : $\frac{a}{b}$. Insérer une expression avec : #{2x + 1}"
      v-model="props.Data.Parts"
    ></v-textarea>
  </container>
</template>

<script setup lang="ts">
import { itemize } from "@/controller/editor";
import type { TextBlock } from "@/controller/exercice_gen";
import { computed } from "@vue/runtime-core";
import Container from "./Container.vue";

interface Props {
  Data: TextBlock;
}
const props = defineProps<Props>();

const colorByKind = ["", "green", "orange"];

const emit = defineEmits<{
  (e: "delete"): void;
}>();

const spans = computed(() =>
  itemize(props.Data.Parts).map(chunk => ({
    color: colorByKind[chunk.Kind],
    text: chunk.Content
  }))
);
</script>

<style>
.v-field__field {
  padding-top: 0;
}

.v-field__input {
  padding-left: 4px;
}
</style>
