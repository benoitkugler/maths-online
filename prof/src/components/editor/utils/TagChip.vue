<template>
  <div class="bg-white pb-1 rounded mx-1 d-inline">
    <v-chip
      elevation="0"
      size="x-small"
      label
      :color="tagColor(tag.Section)"
      :style="{ cursor: props.pointer ? 'pointer' : '' }"
      >{{ tagString(tag.Tag) }}</v-chip
    >
  </div>
</template>

<script setup lang="ts">
import type { TagSection } from "@/controller/api_gen";
import { tagColor, tagString } from "@/controller/editor";
interface Props {
  tag: Optional<TagSection, "Section">;
  pointer?: boolean;
}
const props = defineProps<Props>();

type Optional<T, K extends keyof T> = Pick<Partial<T>, K> & Omit<T, K>;
</script>
