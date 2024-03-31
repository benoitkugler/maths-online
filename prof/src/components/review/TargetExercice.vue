<template>
  <exercicegroup-pannel
    :group="props.data.Group"
    :all-tags="props.data.AllTags || []"
  ></exercicegroup-pannel>
</template>

<script setup lang="ts">
import type { TargetExercice } from "@/controller/api_gen";
import ExercicegroupPannel from "../editor/exercices/ExercicegroupPannel.vue";

interface Props {
  data: TargetExercice;
}

const props = defineProps<Props>();
</script>
