<template>
  <v-chip
    size="x-small"
    label
    class="mx-1"
    :color="tagColor(tag)"
    :style="{ cursor: props.pointer ? 'pointer' : '' }"
    >{{ tagString(tag) }}</v-chip
  >
</template>

<script setup lang="ts">
import { tagColor, tagString } from "@/controller/editor";
interface Props {
  tag: string;
  pointer: boolean;
}
const props = defineProps<Props>();
</script>
