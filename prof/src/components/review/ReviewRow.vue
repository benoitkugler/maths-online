<template>
  <v-list-item
    link
    class="py-1 my-2"
    :title="`(${props.review.Id}) ${props.review.Title}`"
    :subtitle="labels[props.review.Kind]"
  >
    <template #prepend>
      <v-chip class="mx-1" size="small" label :color="MatiereColor">{{
        props.review.Matiere
      }}</v-chip>
      <v-chip class="mr-4" size="small" label :color="LevelColor">{{
        props.review.Level
      }}</v-chip>
    </template>
    <template #append>
      <v-chip prepend-icon="mdi-comment">
        {{ props.review.NbComments }}
      </v-chip>
    </template>
  </v-list-item>
</template>

<script setup lang="ts">
import { ReviewKindLabels, type ReviewHeader } from "@/controller/api_gen";
import TagChip from "../editor/utils/TagChip.vue";
import { LevelColor, MatiereColor } from "@/controller/editor";

interface Props {
  review: ReviewHeader;
}

const props = defineProps<Props>();

// const emit = defineEmits<{}>();

const labels = ReviewKindLabels;
</script>
