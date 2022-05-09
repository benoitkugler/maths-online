<template>
  <v-list-item
    dense
    rounded
    :class="{
      'py-0': true
    }"
    @click="emit('clicked', props.question)"
  >
    <v-list-item-media class="mr-3">
      <v-btn
        size="x-small"
        icon
        @click.stop="emit('delete', props.question)"
        title="Supprimer"
      >
        <v-icon icon="mdi-delete" color="red" size="small"></v-icon>
      </v-btn>
      <v-btn
        v-if="!props.question.IsInGroup"
        class="mx-1"
        size="x-small"
        icon
        @click.stop="emit('duplicate', props.question)"
        title="Dupliquer en différenciant la difficulté"
      >
        <v-icon icon="mdi-content-copy" color="info" size="small"></v-icon>
      </v-btn>
    </v-list-item-media>
    <i>
      <small> ({{ props.question.Id }}) </small>
    </i>
    <div class="ml-2">
      {{ props.question.Title ? props.question.Title : "..." }}
    </div>
    <v-spacer></v-spacer>
    <v-list-item-media>
      <TagChip :tag="tag" :key="tag" v-for="tag in question.Tags"></TagChip>
    </v-list-item-media>
  </v-list-item>
</template>

<script setup lang="ts">
import type { QuestionHeader } from "@/controller/api_gen";
import TagChip from "./utils/TagChip.vue";

interface Props {
  question: QuestionHeader;
}
const props = defineProps<Props>();
const emit = defineEmits<{
  (e: "delete", question: QuestionHeader): void;
  (e: "clicked", question: QuestionHeader): void;
  (e: "duplicate", question: QuestionHeader): void;
}>();
</script>
