<template>
  <v-card class="mb-2">
    <v-row class="bg-secondary pa-2">
      <v-col md="9" align-self="center">
        Réponse ordonnée attendue.
        <small>
          Les champs sont interprétés comme du code LaTeX (utiliser #{} pour une
          expression).
        </small>
      </v-col>
      <v-spacer></v-spacer>
      <v-col align-self="center" style="text-align: right">
        <v-btn
          icon
          @click="addAnswer"
          title="Ajouter un élément"
          size="x-small"
          class="mr-2 my-2"
        >
          <v-icon icon="mdi-plus" color="green" small></v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-row no-gutters>
      <v-col>
        <v-list>
          <v-list-item
            v-for="(param, index) in props.modelValue.Answer"
            class="pr-0"
          >
            <text-part-field
              v-model="props.modelValue.Answer![index]"
              force-latex
            >
            </text-part-field>

            <v-btn
              icon
              size="small"
              flat
              @click="removeAnswer(index)"
              title="Supprimer cet élément"
            >
              <v-icon icon="mdi-delete" color="red"></v-icon>
            </v-btn>
          </v-list-item>
        </v-list>
      </v-col>
    </v-row>
  </v-card>

  <v-card>
    <v-row class="bg-secondary pa-2">
      <v-col md="9" align-self="center">
        Champs additionnels.
        <small>
          Les champs sont interprétés comme du code LaTeX (utiliser #{} pour une
          expression).
        </small>
      </v-col>
      <v-spacer></v-spacer>
      <v-col align-self="center" style="text-align: right">
        <v-btn
          icon
          @click="addAdditionalProposal"
          title="Ajouter un élément"
          size="x-small"
          class="mr-2 my-2"
        >
          <v-icon icon="mdi-plus" color="green" small></v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-row no-gutters>
      <v-col>
        <v-list>
          <v-list-item
            v-for="(param, index) in props.modelValue.AdditionalProposals"
            class="pr-0"
          >
            <text-part-field
              v-model="props.modelValue.AdditionalProposals![index]"
              force-latex
            >
            </text-part-field>

            <v-btn
              icon
              size="small"
              flat
              @click="removeAdditionalProposal(index)"
              title="Supprimer cet élément"
            >
              <v-icon icon="mdi-delete" color="red"></v-icon>
            </v-btn>
          </v-list-item>
        </v-list>
      </v-col>
    </v-row>
  </v-card>

  <v-row class="my-2">
    <v-col>
      <v-text-field
        label="Préfixe"
        hint="Code LaTeX ajouté devant le champ de réponse. (Optionnel)."
        v-model="props.modelValue.Label"
        variant="outlined"
        density="compact"
      >
      </v-text-field>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { OrderedListFieldBlock } from "@/controller/exercice_gen";
import { TextKind } from "@/controller/exercice_gen";
import TextPartField from "./TextPartField.vue";

interface Props {
  modelValue: OrderedListFieldBlock;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: OrderedListFieldBlock): void;
}>();

function addAnswer() {
  props.modelValue.Answer?.push({ Kind: TextKind.StaticMath, Content: "x" });
}

function removeAnswer(index: number) {
  props.modelValue.Answer?.splice(index, 1);
}

function addAdditionalProposal() {
  props.modelValue.AdditionalProposals?.push({
    Kind: TextKind.StaticMath,
    Content: "x"
  });
}

function removeAdditionalProposal(index: number) {
  props.modelValue.AdditionalProposals?.splice(index, 1);
}
</script>

<style></style>
