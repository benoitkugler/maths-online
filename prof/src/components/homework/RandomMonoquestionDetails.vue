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
      <v-row>
        <v-col cols="12">
          <v-select
            label="Choix de la difficulté"
            :color="DifficultyColor"
            variant="outlined"
            density="compact"
            chips
            multiple
            :items="[
              DifficultyTag.Diff1,
              DifficultyTag.Diff2,
              DifficultyTag.Diff3
            ]"
            :model-value="actualTags"
            @update:model-value="updateTags"
          >
          </v-select>
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
import { copy } from "@/controller/utils";
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
  DifficultyTag.Diff3
];

const actualTags = computed(() => {
  // replace empty list by all tags
  const isEmpty = !props.randomMonoquestion.Difficulty?.length;
  if (isEmpty) return diffChoices;
  const out = copy(props.randomMonoquestion.Difficulty || []);
  out.sort();
  return out;
});

function updateTags(l: DifficultyTag[]) {
  // replace all tags by empty list
  if (l.length == diffChoices.length) l = [];
  props.randomMonoquestion.Difficulty = l;
}

const subtitle = computed(
  () => `Groupe de questions ${props.randomMonoquestion.IdQuestiongroup}`
);
</script>

<style scoped></style>
