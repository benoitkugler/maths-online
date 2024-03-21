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
        <v-card-subtitle class="py-2">Paramètres aléatoires</v-card-subtitle>
      </v-col>
      <v-spacer></v-spacer>
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
    <div style="height: 61vh; overflow-y: auto">
      <v-row no-gutters>
        <v-col cols="12">
          <HiglightedText
            v-model="rawParameters"
            :tokenizer="tokenize"
            :center="false"
            focus-color="lightgreen"
            @blur="updateParams"
            @focus="errorParameters = null"
            :border-color="errorParameters != null ? '#FFB74D' : undefined"
          ></HiglightedText>
        </v-col>
        <v-col cols="12" v-if="props.dual">
          <HiglightedText
            v-model="rawSharedParameters"
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
  </div>
</template>

<script setup lang="ts">
import type { ErrParameters, Parameters } from "@/controller/api_gen";
import { ref, watch } from "vue";
import HiglightedText from "../utils/HiglightedText.vue";
import {
  parametersToString,
  parseParameters,
  tokenize,
} from "./parameters_editor";
import ParametersHelp from "./ParametersHelp.vue";
import SnackErrorParameters from "./SnackErrorParameters.vue";

interface Props {
  parameters: Parameters;
  sharedParameters: Parameters;
  dual: boolean;
  isLoading: boolean;
  isValidated: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update", ps: Parameters, shared: Parameters): void;
}>();

// to avoid cursor issues, we use a string as source of truth,
// delaying the parsing/validation to the blur event
const rawParameters = ref(parametersToString(props.parameters));
const rawSharedParameters = ref(parametersToString(props.sharedParameters));
watch(props, (newP, oldP) => {
  rawParameters.value = parametersToString(props.parameters);
  rawSharedParameters.value = parametersToString(props.sharedParameters);
});

const errorParameters = ref<ErrParameters | null>(null);
function updateParams() {
  const { params, error: error1 } = parseParameters(rawParameters.value);
  const { params: shared, error: error2 } = parseParameters(
    rawSharedParameters.value
  );
  if (error1 != null) {
    errorParameters.value = error1;
  } else if (error2 != null) {
    errorParameters.value = error2;
  } else {
    emit("update", params, shared);
  }
}

const showHelp = ref(false);
</script>

<style scoped></style>
