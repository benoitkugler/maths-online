<template>
  <v-row class="mt-2">
    <v-col cols="5" class="pb-0">
      <text-part-field
        v-model="props.modelValue.Label"
        label="Préfixe"
        hint="Ajouté devant le champ de réponse. Optionnel"
      >
      </text-part-field>
    </v-col>
    <v-col cols="7" class="pb-0">
      <v-text-field
        variant="outlined"
        density="compact"
        v-model="props.modelValue.Expression"
        label="Réponse"
        hint="Expression"
        :color="color"
      >
      </v-text-field>
    </v-col>
    <v-col cols="12" class="pt-0">
      <v-checkbox
        :model-value="isComparaisonStrict"
        @update:model-value="changeComparaison"
        color="secondary"
        label="Comparaison stricte"
        density="compact"
        :messages="[comparaisonMessage]"
      ></v-checkbox>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { colorByKind } from "@/controller/editor";
import type { ExpressionFieldBlock } from "@/controller/exercice_gen";
import { ComparisonLevel, TextKind } from "@/controller/exercice_gen";
import { computed } from "@vue/runtime-core";
import { $computed } from "vue/macros";
import TextPartField from "./TextPartField.vue";

interface Props {
  modelValue: ExpressionFieldBlock;
}
const props = defineProps<Props>();
const color = colorByKind[TextKind.Expression];

const emit = defineEmits<{
  (event: "update:modelValue", value: ExpressionFieldBlock): void;
}>();

let isComparaisonStrict = $computed(
  () =>
    props.modelValue.ComparisonLevel != ComparisonLevel.ExpandedSubstitutions
);

const comparaisonMessage = computed(() => {
  return isComparaisonStrict
    ? "Les expressions sont peu transformées : (x+1)^2 et x^2 + 2x + 1 ne sont pas considérées comme égales."
    : "Les formules usuelles de développement et factorisation sont appliquées en évaluant la réponse : (x+1)^2 et x^2 + 2x + 1 sont considérées égales.";
});

function changeComparaison(b: boolean) {
  props.modelValue.ComparisonLevel = b
    ? ComparisonLevel.SimpleSubstitutions
    : ComparisonLevel.ExpandedSubstitutions;
}
</script>

<style></style>
