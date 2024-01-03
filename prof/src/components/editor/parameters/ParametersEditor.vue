<template>
  <ParametersHelp v-model="showHelp"></ParametersHelp>

  <SnackErrorParameters
    :error="errorParameters"
    @close="errorParameters = null"
  >
  </SnackErrorParameters>

  <div class="px-2">
    <v-row
      :style="{
        'background-color': props.isLoading
          ? 'lightgrey'
          : props.isValidated
          ? 'lightgreen'
          : '#FFB74D',
      }"
      class="rounded"
      no-gutters
    >
      <v-col md="8" align-self="center">
        <v-card-subtitle class="py-2">{{ title }}</v-card-subtitle>
      </v-col>
      <v-spacer></v-spacer>
      <v-col cols="auto" align-self="center" class="px-2">
        <v-btn-toggle
          v-if="props.showSwitch != null"
          mandatory
          rounded
          density="compact"
          :model-value="props.showSwitch ? 0 : 1"
          @update:model-value="(v) => emit('switch', v == 0)"
        >
          <v-tooltip
            location="bottom"
            text="Afficher les paramètres partagés par toutes les questions"
          >
            <template v-slot:activator="{ props: innerProps }">
              <v-btn
                v-bind="innerProps"
                variant="outlined"
                rounded
                icon="mdi-file-multiple"
              ></v-btn>
            </template>
          </v-tooltip>
          <v-tooltip
            location="bottom"
            text="Afficher les paramètres uniques de la question"
          >
            <template v-slot:activator="{ props: innerProps }">
              <v-btn
                v-bind="innerProps"
                variant="outlined"
                rounded
                icon="mdi-file"
              ></v-btn>
            </template>
          </v-tooltip>
        </v-btn-toggle>
      </v-col>
      <v-col cols="auto" align-self="center" style="text-align: right">
        <v-btn
          icon
          @click="showHelp = true"
          title="Aide"
          size="x-small"
          class="mr-2 my-2"
        >
          <v-icon icon="mdi-help" color="info"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-row no-gutters style="height: 61vh; overflow-y: auto">
      <v-col>
        <HiglightedText
          v-model="rawText"
          :tokenizer="tokenize"
          :center="false"
          focus-color="lightgreen"
          @blur="updateParams"
          @focus="errorParameters = null"
          :border-color="errorParameters != null ? '#FFB74D' : undefined"
        ></HiglightedText>
      </v-col>
    </v-row>
  </div>
</template>

<script setup lang="ts">
import type { ErrParameters, Parameters } from "@/controller/api_gen";
import { ref, computed, watch } from "vue";
import HiglightedText from "../utils/HiglightedText.vue";
import {
  parametersToString,
  parseParameters,
  tokenize,
} from "./parameters_editor";
import ParametersHelp from "./ParametersHelp.vue";
import SnackErrorParameters from "./SnackErrorParameters.vue";

interface Props {
  modelValue: Parameters;
  isLoading: boolean;
  isValidated: boolean;
  showSwitch: boolean | null;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:modelValue", ps: Parameters): void;
  (e: "switch", showShared: boolean): void;
}>();

const title = computed(() =>
  props.showSwitch == null
    ? "Paramètres aléatoires"
    : props.showSwitch
    ? "Paramètres partagés"
    : "Paramètres de la question"
);

// to avoid cursor issues, we use a string as source of truth,
// delaying the parsing/validation to the blur event
const rawText = ref(parametersToString(props.modelValue));
watch(props, (newP, oldP) => {
  // do not update the layout if it is triggered by error
  //   if (
  //     JSON.stringify(newP.modelValue || []) !=
  //     JSON.stringify(oldP.modelValue || [])
  //   ) {
  rawText.value = parametersToString(props.modelValue);
  //   }
  //   console.log(
  //     props.modelValue,
  //     JSON.stringify(newP.modelValue),
  //     JSON.stringify(oldP.modelValue)
  //   );
});

const errorParameters = ref<ErrParameters | null>(null);
function updateParams() {
  const { params, error } = parseParameters(rawText.value);
  if (error != null) {
    errorParameters.value = error;
  } else {
    emit("update:modelValue", params);
  }
}

const showHelp = ref(false);
</script>

<style scoped></style>
