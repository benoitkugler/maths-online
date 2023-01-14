<template>
  <v-chip
    size="x-small"
    label
    class="mx-1"
    :color="tagColor(tag)"
    :style="{ cursor: props.pointer ? 'pointer' : '' }"
    >{{ tagString(tag.Tag) }}</v-chip
  >
</template>

<script setup lang="ts">
import type { TagSection } from "@/controller/api_gen";
import { tagColor, tagString } from "@/controller/editor";
interface Props {
  tag: TagSection;
  pointer?: boolean;
}
const props = defineProps<Props>();
</script>
