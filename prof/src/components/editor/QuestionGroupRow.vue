<template>
  <v-expansion-panels class="my-1" v-model="state">
    <v-expansion-panel>
      <v-expansion-panel-title class="py-0 bg-lime-lighten-5 rounded">
        <v-row no-gutters justify="space-between">
          <v-col cols="auto" style="text-align: left" align-self="center">
            <v-row no-gutters>
              <v-col cols="12">
                {{ props.group.Title }}
              </v-col>
            </v-row>
            <v-row no-gutters v-if="state == 0">
              <v-col cols="12">
                <small>
                  ({{ props.group.Questions?.length || 0 }} question(s)
                  affichée(s) / {{ props.group.Size }} questions)
                </small>
              </v-col>
            </v-row>
          </v-col>
          <v-col cols="auto" align-self="center">
            <TagListField
              :readonly="!isEditable"
              :model-value="tags"
              @update:model-value="(l) => emit('updateTags', l)"
              :all-tags="props.allTags"
              y-padding
            >
            </TagListField>
          </v-col>
        </v-row>
      </v-expansion-panel-title>
      <v-expansion-panel-text class="px-0">
        <QuestionRow
          :question="question"
          v-for="question in props.group.Questions"
          :key="question.Id"
          @clicked="emit('clicked', question)"
          @delete="emit('delete', question)"
          @update-public="(id, b) => emit('updatePublic', id, b)"
        >
        </QuestionRow>
      </v-expansion-panel-text>
    </v-expansion-panel>
  </v-expansion-panels>
</template>

<script setup lang="ts">
import {
  Visibility,
  type QuestionGroup,
  type QuestionHeader,
} from "@/controller/api_gen";
import { commonGroupTags } from "@/controller/editor";
import { computed, ref } from "@vue/runtime-core";
import QuestionRow from "./QuestionRow.vue";
import TagListField from "./TagListField.vue";

interface Props {
  group: QuestionGroup;
  allTags: string[];
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "delete", question: QuestionHeader): void;
  (e: "clicked", question: QuestionHeader): void;
  (e: "updatePublic", questionID: number, isPublic: boolean): void;
  (e: "updateTags", tags: string[]): void;
}>();

const tags = computed(() => commonGroupTags(props.group.Questions || []));

const isEditable = computed(
  () =>
    props.group.Questions?.every(
      (qu) => qu.Origin.Visibility == Visibility.Personnal
    ) || true
);
const state = ref<number | null>(null);
</script>

<style scoped>
:deep(.v-expansion-panel-text__wrapper) {
  padding-left: 2px;
  padding-right: 2px;
}
</style>
