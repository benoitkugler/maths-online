<template>
  <v-row>
    <v-col md="9" align-self="center">
      Réponses possibles (insérer du code LaTeX avec $$ et une expression avec
      #{}).</v-col
    >
    <v-spacer></v-spacer>
    <v-col align-self="center" style="text-align: right">
      <v-btn
        icon
        @click="addProposal"
        title="Ajouter un choix"
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
          v-for="(param, index) in props.modelValue.Proposals"
          class="pr-0"
        >
          <v-list-item-title style="width: 120px">
            Choix {{ index }} :
          </v-list-item-title>
          <interpolated-text v-model="props.modelValue.Proposals![index]">
          </interpolated-text>

          <v-btn
            icon
            size="small"
            flat
            @click="removeProposal(index)"
            title="Supprimer ce choix"
          >
            <v-icon icon="mdi-delete" color="red"></v-icon>
          </v-btn>
        </v-list-item>
      </v-list>
    </v-col>
  </v-row>
  <v-row class="mt-2 mb-0 pb-0">
    <v-col class="mb-0 pb-0">
      <v-text-field
        variant="outlined"
        density="compact"
        v-model="props.modelValue.Answer"
        label="Réponse"
        :hint="hint"
        :color="color"
      >
      </v-text-field>
    </v-col>
  </v-row>
  <v-row class="mt-0 pt-0">
    <v-col class="mt-0 pt-0">
      <v-switch
        label="Afficher comme menu déroulant"
        v-model="props.modelValue.AsDropDown"
        color="secondary"
        hide-details
      ></v-switch>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { colorByKind } from "@/controller/editor";
import type { RadioFieldBlock } from "@/controller/exercice_gen";
import { TextKind } from "@/controller/exercice_gen";
import { computed } from "@vue/runtime-core";
import InterpolatedText from "../InterpolatedText.vue";

interface Props {
  modelValue: RadioFieldBlock;
}
const props = defineProps<Props>();
const color = colorByKind[TextKind.Expression];

const emit = defineEmits<{
  (event: "update:modelValue", value: RadioFieldBlock): void;
}>();

const hint = computed(
  () =>
    "Expression s'évaluant en un indice dans la liste des choix (de 0 à " +
    (props.modelValue.Proposals?.length || 0 - 1) +
    ")."
);

function addProposal() {
  props.modelValue.Proposals?.push("");
}

function removeProposal(index: number) {
  props.modelValue.Proposals?.splice(index, 1);
}
</script>

<style></style>
