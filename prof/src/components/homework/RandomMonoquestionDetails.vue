<template>
  <v-card title="Paramètres de la question" :subtitle="subtitle">
    <v-card-text class="mt-2">
      <v-row>
        <v-col>
          <v-text-field
            label="Nombre de répétitions"
            density="compact"
            variant="outlined"
            type="number"
            min="1"
            hide-details
            v-model.number="props.randomMonoquestion.NbRepeat"
          ></v-text-field>
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-text-field
            label="Barème"
            density="compact"
            variant="outlined"
            type="number"
            min="0"
            hint="Points attribués pour une question"
            v-model.number="props.randomMonoquestion.Bareme"
          >
          </v-text-field>
        </v-col>
      </v-row>
      <v-row justify="center">
        <v-col cols="auto">
          <label class="text-grey">Choix de la difficulté : </label>
          <v-btn-toggle
            density="compact"
            :color="DifficultyColor"
            multiple
            variant="outlined"
            rounded
            :model-value="difficultiesIndices"
            @update:model-value="updateDiffs"
          >
            <v-btn>{{ DifficultyTag.Diff1 }}</v-btn>
            <v-btn>{{ DifficultyTag.Diff2 }}</v-btn>
            <v-btn>{{ DifficultyTag.Diff3 }}</v-btn>
          </v-btn-toggle>
        </v-col>
      </v-row>
    </v-card-text>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn color="success" @click="emit('update', props.randomMonoquestion)">
        Enregistrer
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import { DifficultyTag, type RandomMonoquestion } from "@/controller/api_gen";
import { DifficultyColor } from "@/controller/editor";
import { computed } from "vue";

interface Props {
  randomMonoquestion: RandomMonoquestion;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update", randomMonoquestion: RandomMonoquestion): void;
}>();

const diffChoices = [
  DifficultyTag.Diff1,
  DifficultyTag.Diff2,
  DifficultyTag.Diff3,
];

const difficultiesIndices = computed(() => {
  const isEmpty = !props.randomMonoquestion.Difficulty?.length;
  const out: number[] = [];
  diffChoices.forEach((v, i) => {
    if (isEmpty || props.randomMonoquestion.Difficulty?.includes(v))
      out.push(i);
  });
  return out;
});

function updateDiffs(l: number[]) {
  if (l.length == diffChoices.length) l = [];
  props.randomMonoquestion.Difficulty = l.map((i) => diffChoices[i]);
}

const subtitle = computed(
  () => `Groupe de questions ${props.randomMonoquestion.IdQuestiongroup}`
);
</script>

<style scoped></style>
