<template>
  <v-row class="fill-width">
    <v-col v-for="(travail, index) in props.travaux" :key="index" cols="6">
      <MarksStatsColumn
        :data="(props.data.Marks || {})[travail.Id]"
        :sheet="props.sheets.get(travail.IdSheet)!"
      ></MarksStatsColumn>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { HomeworkMarksOut, SheetExt, Travail } from "@/controller/api_gen";
import MarksStatsColumn from "./MarksStatsColumn.vue";

interface Props {
  data: HomeworkMarksOut;
  travaux: Travail[];
  sheets: Map<number, SheetExt>;
}

const props = defineProps<Props>();
</script>

<style>
.fill-width {
  overflow-x: auto;
  flex-wrap: nowrap;
}
</style>
