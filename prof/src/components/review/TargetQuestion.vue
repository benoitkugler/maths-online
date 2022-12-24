<template>
  <v-row>
    <v-col>
      <questiongroup-pannel
        :group="props.data.Group"
        :variants="props.data.Variants || []"
        :all-tags="props.data.AllTags || []"
        @preview="(qu) => preview?.showQuestion(qu)"
      ></questiongroup-pannel>
    </v-col>
    <v-col cols="auto">
      <client-preview ref="preview"></client-preview>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { TargetQuestion } from "@/controller/api_gen";
import { $ref } from "vue/macros";
import ClientPreview from "../editor/ClientPreview.vue";
import QuestiongroupPannel from "../editor/questions/QuestiongroupPannel.vue";

interface Props {
  data: TargetQuestion;
}

const props = defineProps<Props>();

let preview = $ref<InstanceType<typeof ClientPreview> | null>(null);
</script>
