<template>
  <v-sheet
    :class="(isActive ? 'bg-secondary' : 'bg-grey') + '  rounded-lg pt-0'"
    :style="{ height: isActive ? '30px' : '20px' }"
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
  emit("drop", eventData.index);
}

function onDragOver(ev: DragEvent) {
  ev.preventDefault();
  isActive = true;
}
</script>
