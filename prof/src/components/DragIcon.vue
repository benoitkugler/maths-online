<template>
  <v-icon
    size="large"
    class="pr-2"
    style="cursor: grab"
    @dragstart="(e) => emit('start', e)"
    draggable="true"
    :color="props.color"
    icon="mdi-drag-vertical"
  ></v-icon>
</template>

<script setup lang="ts">
interface Props {
  color: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "start", ev: DragEvent): void;
}>();
</script>
