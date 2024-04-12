<template>
  <questiongroup-pannel
    :group="props.data.Group"
    :variants="props.data.Variants || []"
    :all-tags="props.data.AllTags || []"
  ></questiongroup-pannel>
</template>

<script setup lang="ts">
import type { TargetQuestion } from "@/controller/api_gen";
import QuestiongroupPannel from "../editor/questions/QuestiongroupPannel.vue";

interface Props {
  data: TargetQuestion;
}

const props = defineProps<Props>();
</script>
