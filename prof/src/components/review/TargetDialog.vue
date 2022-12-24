<template>
  <v-card>
    <v-row class="pa-2" no-gutters>
      <v-col>
        <v-card-title>Contenu de la publication</v-card-title>
      </v-col>
      <v-col cols="auto">
        <v-btn icon flat @click="emit('back')"
          ><v-icon>mdi-close</v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-card-text>
      <v-row justify="center" v-if="content == null">
        <v-col cols="auto" align-self="center">
          <v-progress-circular indeterminate></v-progress-circular>
        </v-col>
        <v-col cols="auto" align-self="center"> Chargement... </v-col>
      </v-row>
      <v-row v-else>
        <v-col>
          <component
            :data="content.Props.Data"
            :is="content.Component"
          ></component>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import {
  TargetContentKind,
  type ReviewHeader,
  type TargetContent,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { markRaw, onMounted, type Component } from "vue";
import { $ref } from "vue/macros";
import TargetExerciceVue from "./TargetExercice.vue";
import TargetQuestionVue from "./TargetQuestion.vue";
import TargetTrivialVue from "./TargetTrivial.vue";

interface Props {
  review: ReviewHeader;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "back"): void;
}>();

let content = $ref<targetContent | null>(null);

onMounted(loadTarget);

async function loadTarget() {
  const res = await controller.ReviewLoadTarget({
    "id-review": props.review.Id,
  });
  if (res === undefined) return;
  content = newContent(res.Content);
}

interface targetContent {
  Props: TargetContent;
  Component: Component;
}

function newContent(data: TargetContent): targetContent {
  switch (data.Kind) {
    case TargetContentKind.TargetQuestion:
      return { Props: data, Component: markRaw(TargetQuestionVue) };
    case TargetContentKind.TargetExercice:
      return { Props: data, Component: markRaw(TargetExerciceVue) };
    case TargetContentKind.TargetTrivial:
      return { Props: data, Component: markRaw(TargetTrivialVue) };
  }
}
</script>
