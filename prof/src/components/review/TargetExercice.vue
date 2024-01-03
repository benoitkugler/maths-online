<template>
  <v-row>
    <v-col>
      <exercicegroup-pannel
        :group="props.data.Group"
        :all-tags="props.data.AllTags || []"
        @preview="(ex) => preview?.showExercice(ex)"
      ></exercicegroup-pannel>
    </v-col>
    <v-col cols="auto">
      <client-preview ref="preview"></client-preview>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { TargetExercice } from "@/controller/api_gen";
import { ref } from "vue";
import ClientPreview from "../editor/ClientPreview.vue";
import ExercicegroupPannel from "../editor/exercices/ExercicegroupPannel.vue";

interface Props {
  data: TargetExercice;
}

const props = defineProps<Props>();

const preview = ref<InstanceType<typeof ClientPreview> | null>(null);
</script>
