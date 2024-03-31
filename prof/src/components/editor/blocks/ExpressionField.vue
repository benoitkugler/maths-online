<template>
  <v-row class="pb-1 mb-1">
    <v-col cols="5" align-self="center">
      <interpolated-text
        v-model="props.modelValue.Label"
        label="Préfixe (optionnel)"
        hint="Ajouté devant le champ de réponse."
      >
      </interpolated-text>
    </v-col>
    <v-col cols="7" align-self="center">
      <v-text-field
        class="mt-5"
        variant="outlined"
        density="compact"
        v-model="props.modelValue.Expression"
        label="Réponse"
        hint="Expression"
        persistent-hint
        :color="ExpressionColor"
        @blur="emitUpdate"
      >
      </v-text-field>
    </v-col>
    <v-col cols="12" class="mt-1">
      <v-select
        variant="outlined"
        density="compact"
        :messages="[comparisonMessage]"
        :items="comparisonSelectItems"
        label="Mode de comparaison"
        v-model="props.modelValue.ComparisonLevel"
        @update:model-value="emitUpdate()"
      >
      </v-select>
    </v-col>
    <v-col cols="12">
      <v-checkbox
        density="compact"
        label="Afficher une aide pour les fractions"
        v-model="props.modelValue.ShowFractionHelp"
        @update:model-value="emitUpdate()"
        messages="Si la réponse est une fraction, une aide sur la syntaxe est affichée dans le champ de réponse."
      ></v-checkbox>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { ExpressionFieldBlock, Variable } from "@/controller/api_gen";
import { ComparisonLevel } from "@/controller/api_gen";
import { ExpressionColor } from "@/controller/editor";
import { computed } from "@vue/runtime-core";
import InterpolatedText from "../utils/InterpolatedText.vue";

interface Props {
  modelValue: ExpressionFieldBlock;
  availableParameters: Variable[];
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: ExpressionFieldBlock): void;
}>();

function emitUpdate() {
  emit("update:modelValue", props.modelValue);
}

const comparisonMessage = computed(() => {
  switch (props.modelValue.ComparisonLevel) {
    case ComparisonLevel.SimpleSubstitutions:
      return "Les expressions sont peu transformées : (x+1)^2 et x^2 + 2x + 1 ne sont pas considérées comme égales.";
    case ComparisonLevel.ExpandedSubstitutions:
      return "Les formules usuelles de développement et factorisation sont appliquées en évaluant la réponse : (x+1)^2 et x^2 + 2x + 1 sont considérées égales.";
    case ComparisonLevel.AsLinearEquation:
      return "L'expression définit une équation cartésienne, comparée à un facteur près.";
    default:
      return "";
  }
});

const comparisonSelectItems = [
  { title: "Comparaison stricte", value: ComparisonLevel.SimpleSubstitutions },
  { title: "Comparaison large", value: ComparisonLevel.ExpandedSubstitutions },
  { title: "Equation cartésienne", value: ComparisonLevel.AsLinearEquation },
];
</script>

<style></style>
