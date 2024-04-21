<template>
  <v-row>
    <v-col md="9" align-self="center">
      Réponses possibles (insérer du code LaTeX avec $$ et une expression avec
      &2x+1&).</v-col
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
          :key="index"
          class="pr-0"
        >
          <v-row no-gutters>
            <v-col cols="3" align-self="center">
              <v-list-item-title>
                Choix {{ index + 1 }} :
              </v-list-item-title></v-col
            >
            <v-col>
              <interpolated-text
                v-model="props.modelValue.Proposals![index]"
                @update:model-value="emitUpdate"
              >
              </interpolated-text>
            </v-col>
            <v-col cols="auto">
              <v-btn
                icon
                size="small"
                flat
                @click="removeProposal(index)"
                title="Supprimer ce choix"
              >
                <v-icon icon="mdi-delete" color="red"></v-icon>
              </v-btn>
            </v-col>
          </v-row>
        </v-list-item>
      </v-list>
    </v-col>
  </v-row>
  <v-row class="mt-2 mb-0 pb-0">
    <v-col class="mb-0 pb-0">
      <expression-field
        v-model="props.modelValue.Answer"
        @update:model-value="emitUpdate"
        label="Réponse"
        :hint="hint"
        prefix="Choix"
      >
      </expression-field>
    </v-col>
  </v-row>
  <v-row class="mt-0 pt-0">
    <v-col class="mt-0 pt-0">
      <v-switch
        label="Afficher comme menu déroulant"
        v-model="props.modelValue.AsDropDown"
        @update:model-value="emitUpdate"
        color="secondary"
        hide-details
      ></v-switch>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { RadioFieldBlock, Variable } from "@/controller/api_gen";
import { computed } from "vue";
import ExpressionField from "../utils/ExpressionField.vue";
import InterpolatedText from "../utils/InterpolatedText.vue";

interface Props {
  modelValue: RadioFieldBlock;
  availableParameters: Variable[];
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: RadioFieldBlock): void;
}>();

const hint = computed(
  () =>
    "Expression s'évaluant en un indice dans la liste des choix, de 1 à " +
    (props.modelValue.Proposals?.length || 0) +
    "."
);

function emitUpdate() {
  emit("update:modelValue", props.modelValue);
}

function addProposal() {
  props.modelValue.Proposals?.push("");
  emitUpdate();
}

function removeProposal(index: number) {
  props.modelValue.Proposals?.splice(index, 1);
  emitUpdate();
}
</script>

<style scoped>
.fix-prefix:deep(.v-field__input) {
  margin-bottom: 5px;
}
</style>
