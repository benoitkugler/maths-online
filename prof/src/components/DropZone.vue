<template>
  <div
    @dragover="onDragOver"
    @dragleave="isActive = false"
    @drop="onDrop"
    class="py-3"
  >
    <v-sheet
      :class="
        (isActive ? 'bg-blue-lighten-2' : 'bg-blue-lighten-5') +
        '  rounded mx-1'
      "
      style="height: 8px"
    ></v-sheet>
  </div>
</template>

<script setup lang="ts">
import { $ref } from "vue/macros";

let isActive = $ref(false);

const emit = defineEmits<{
  (e: "drop", origin: number): void;
}>();

function onDrop(ev: DragEvent) {
  const eventData: { index: number } = JSON.parse(
    ev.dataTransfer?.getData("text/json") || ""
  );
  emit("drop", eventData.index);
}

function onDragOver(ev: DragEvent) {
  ev.preventDefault();
  isActive = true;
}
</script>
