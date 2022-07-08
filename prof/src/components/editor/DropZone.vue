<template>
  <v-sheet
    :class="
      (isActive ? 'bg-blue-lighten-2' : 'bg-blue-lighten-4') +
      '  rounded pt-0 mx-1'
    "
    :style="{ height: isActive ? '20px' : '20px' }"
    @dragover="onDragOver"
    @dragleave="isActive = false"
    @drop="onDrop"
  ></v-sheet>
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
  console.log("drop");

  emit("drop", eventData.index);
}

function onDragOver(ev: DragEvent) {
  console.log("drag over");

  ev.preventDefault();
  isActive = true;
}
</script>
