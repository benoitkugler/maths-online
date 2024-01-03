<template>
  <v-alert
    class="py-2 px-3"
    variant="outlined"
    v-if="props.pattern?.length"
    :color="isComplete ? 'success' : 'info'"
    closable
  >
    <div v-if="isComplete">
      Toutes les ressources inclusent dans
      <span>
        <TagChip
          :tag="tag"
          :key="index"
          v-for="(tag, index) in props.pattern"
          :pointer="false"
        ></TagChip
      ></span>
      sont utilisées.
    </div>

    <div
      v-if="!isComplete"
      style="max-height: 50vh"
      class="overflow-y-auto mr-2"
    >
      <div v-if="missingExercices.length">
        Des exercices dans les catégories suivantes ne sont pas utilisées :
        <v-list>
          <v-list-item v-for="(tags, index) in missingExercices" :key="index">
            <TagChip
              :tag="tag"
              :key="index"
              v-for="(tag, index) in tags || []"
              :pointer="false"
            ></TagChip>
          </v-list-item>
        </v-list>
      </div>
      <div v-if="missingQuestions.length">
        Des questions dans les catégories suivantes ne sont pas utilisées :
        <v-list>
          <v-list-item v-for="(tags, index) in missingQuestions" :key="index">
            <TagChip
              :tag="tag"
              :key="index"
              v-for="(tag, index) in tags || []"
              :pointer="false"
            ></TagChip>
          </v-list-item>
        </v-list>
      </div>
    </div>
  </v-alert>
</template>

<script setup lang="ts">
import type { Tags } from "@/controller/api_gen";
import { computed } from "vue";
import TagChip from "./editor/utils/TagChip.vue";
interface Props {
  pattern: Tags;
  missingExercices: Tags[];
  missingQuestions: Tags[];
}

const props = defineProps<Props>();

const emit = defineEmits<{}>();

const isComplete = computed(
  () => props.missingExercices.length + props.missingQuestions.length == 0
);
</script>
