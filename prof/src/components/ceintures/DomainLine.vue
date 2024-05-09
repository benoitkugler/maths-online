<template>
  <div>
    <template v-for="(stage, i) in stages" :key="i">
      <v-row no-gutters>
        <v-col align-self="center" cols="3" class="text-center">
          <RankIcon :rank="stage.Rank"></RankIcon>
        </v-col>
        <v-col align-self="center" cols="5">
          <div
            v-for="(question, i) in props.stages[stage.Rank].Questions"
            :key="i"
          >
            {{ question.Title }}
          </div>
        </v-col>
        <v-col align-self="center" cols="1">
          <v-icon v-if="needs(stage).length">mdi-chevron-left</v-icon>
        </v-col>
        <v-col align-self="center">
          <stage-chip
            class="my-1"
            v-for="(need, i) in needs(stage)"
            :key="i"
            :stage="need"
            link
            @click="emit('goTo', need.Domain)"
            small
          ></stage-chip>
        </v-col>
      </v-row>
      <v-row no-gutters v-if="stage.Rank != Rank.Noire">
        <v-col class="text-center" cols="3">
          <v-icon>mdi-chevron-down</v-icon></v-col
        >
        <v-col cols="1"></v-col>
        <v-col></v-col>
      </v-row>
    </template>
  </div>
</template>

<script setup lang="ts">
import {
  Ar11_StageHeader,
  Domain,
  Rank,
  Scheme,
  Stage,
} from "@/controller/api_gen";
import { computed } from "vue";
import StageChip from "./StageChip.vue";
import { sameStage } from "@/controller/utils";
import RankIcon from "./RankIcon.vue";

interface Props {
  domain: Domain;
  scheme: Scheme;
  stages: Ar11_StageHeader;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "goTo", d: Domain): void;
}>();

const stages = computed<Stage[]>(() =>
  Object.values(Rank)
    .filter((r) => r != Rank.StartRank)
    .map((r) => ({ Domain: props.domain, Rank: r }))
);

function needs(stage: Stage) {
  const out = [];
  for (const link of props.scheme.Ps || []) {
    if (sameStage(link.Pending, stage)) {
      out.push(link.Need);
    }
  }
  return out;
}
</script>
