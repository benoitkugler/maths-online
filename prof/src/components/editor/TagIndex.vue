<template>
  <v-row no-gutters>
    <v-col v-for="(tag, index) in tags" :key="index">
      <TagChip :tag="tag"></TagChip>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { Section, type TagSection } from "@/controller/api_gen";
import { computed } from "vue";
import TagChip from "./utils/TagChip.vue";

interface Props {
  tags: TagSection[];
}

const props = defineProps<Props>();

const tags = computed(() =>
  props.tags.filter(
    (ts) => ts.Section == Section.Level || ts.Section == Section.Chapter
  )
);
</script>

<style scoped></style>
