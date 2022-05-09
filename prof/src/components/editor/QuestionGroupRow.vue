<template>
  <v-expansion-panels class="my-1" v-model="state">
    <v-expansion-panel>
      <v-expansion-panel-title class="py-0 bg-lime-lighten-4 rounded">
        <v-row no-gutters cols="auto">
          <v-col style="text-align: left" align-self="center">
            <v-row no-gutters>
              <v-col cols="12">
                {{ props.group.Title }}
              </v-col>
              <v-col cols="12" v-if="state == 0">
                <small>
                  ({{ props.group.Questions?.length || 0 }} /
                  {{ props.group.Size }} questions)
                </small>
              </v-col>
            </v-row>
          </v-col>
          <v-col align-self="center" style="text-align: right">
            <TagChip :tag="tag" :key="tag" v-for="tag in tags"></TagChip>
          </v-col>
        </v-row>
      </v-expansion-panel-title>
      <v-expansion-panel-text>
        <QuestionRow
          :question="question"
          v-for="question in props.group.Questions"
          @clicked="emit('clicked', question)"
          @delete="emit('delete', question)"
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
}>();

const tags = computed(() => commonTags(props.group.Questions || []));

const state = ref<number | null>(null);
</script>
