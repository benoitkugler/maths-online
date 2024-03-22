<template>
  <v-row class="my-2" justify="space-between" no-gutters>
    <v-col cols="7" align-self="center">
      <v-text-field
        variant="outlined"
        density="compact"
        :model-value="props.modelValue.Answer"
        @update:model-value="onTypeExpr"
        label="Réponse"
        hint="Expression définissant un ensemble"
        persistent-hint
        :color="ExpressionColor"
        ref="answerField"
      >
      </v-text-field>
    </v-col>
    <v-col cols="auto" align-self="center">
      <v-btn
        size="small"
        label
        class="mx-1"
        @click="insertSpecial('∪')"
        title="Union"
      >
        ∪
      </v-btn>
      <v-btn
        size="small"
        label
        class="mx-1"
        @click="insertSpecial('∩')"
        title="Intersection"
        >∩</v-btn
      >
      <v-btn
        size="small"
        label
        class="mx-1"
        @click="insertSpecial('¬')"
        title="Complémentaire"
        >¬</v-btn
      >
    </v-col>
  </v-row>
  <v-row class="bg-secondary pa-2 rounded" no-gutters>
    <v-col md="9" align-self="center"> Ensemble additionnels proposés</v-col>
    <v-spacer></v-spacer>
    <v-col align-self="center" style="text-align: right">
      <v-btn
        icon
        @click="addAdditionalSet"
        title="Ajouter un ensemble"
        size="x-small"
        class="mr-2"
      >
        <v-icon icon="mdi-plus" color="green" small></v-icon>
      </v-btn>
    </v-col>
  </v-row>
  <v-row no-gutters class="mt-3">
    <v-col>
      <v-row>
        <v-col
          cols="6"
          v-for="(param, index) in props.modelValue.AdditionalSets"
          :key="index"
          class="pr-0 py-1"
        >
          <v-row no-gutters>
            <v-col align-self="center">
              <interpolated-text
                force-latex
                v-model="props.modelValue.AdditionalSets![index]"
              >
              </interpolated-text>
            </v-col>
            <v-col align-self="center" cols="auto">
              <v-btn
                icon
                size="small"
                flat
                @click="removeAdditionalProposal(index)"
                title="Supprimer cet élément"
              >
                <v-icon icon="mdi-delete" color="red"></v-icon>
              </v-btn>
            </v-col>
          </v-row>
        </v-col>
      </v-row>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { type SetFieldBlock, type Variable } from "@/controller/api_gen";
import InterpolatedText from "../utils/InterpolatedText.vue";
import { ExpressionColor } from "@/controller/editor";
import { nextTick } from "vue";
import { ref } from "vue";
import { VTextField } from "vuetify/lib/components/index.mjs";

interface Props {
  modelValue: SetFieldBlock;
  availableParameters: Variable[];
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: SetFieldBlock): void;
}>();

function onTypeExpr(s: string) {
  const l = props.modelValue;
  l.Answer = s;
  emit("update:modelValue", l);
}

function addAdditionalSet() {
  const l = props.modelValue;
  l.AdditionalSets?.push("$y$");
  emit("update:modelValue", l);
}

function removeAdditionalProposal(index: number) {
  const l = props.modelValue;
  l.AdditionalSets?.splice(index, 1);
  emit("update:modelValue", l);
}

const answerField = ref<InstanceType<typeof VTextField> | null>(null);

function insertSpecial(s: string) {
  if (answerField.value == null) return 0;
  const currentText = props.modelValue.Answer;
  const caret = answerField.value?.selectionStart || currentText.length;
  const updated =
    currentText.substring(0, caret) + s + currentText.substring(caret);
  const l = props.modelValue;
  l.Answer = updated;
  emit("update:modelValue", l);
  nextTick(() => {
    answerField.value?.focus();
    answerField.value?.setSelectionRange(caret + 1, caret + 1);
  });
}
</script>

<style></style>
