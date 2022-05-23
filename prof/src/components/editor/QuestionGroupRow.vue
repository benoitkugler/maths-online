<template>
  <v-expansion-panels class="my-1" v-model="state">
    <v-expansion-panel>
      <v-expansion-panel-title class="py-0 bg-lime-lighten-5 rounded">
        <v-row no-gutters>
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
          <v-col align-self="center" style="text-align: right">
            <TagChip :tag="tag" :key="tag" v-for="tag in tags"></TagChip>
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
import type { QuestionGroup, QuestionHeader } from "@/controller/api_gen";
import { commonTags } from "@/controller/editor";
import { computed, ref } from "@vue/runtime-core";
import QuestionRow from "./QuestionRow.vue";
import TagChip from "./utils/TagChip.vue";

interface Props {
  group: QuestionGroup;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "delete", question: QuestionHeader): void;
  (e: "clicked", question: QuestionHeader): void;
  (e: "updatePublic", questionID: number, isPublic: boolean): void;
}>();

const tags = computed(() => commonTags(props.group.Questions || []));

const state = ref<number | null>(null);
</script>

<style scoped>
:deep(.v-expansion-panel-text__wrapper) {
  padding-left: 2px;
  padding-right: 2px;
}
</style>
