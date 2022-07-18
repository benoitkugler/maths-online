<template>
  <div>
    <v-row>
      <v-col cols="11">
        <ProofAssertionVue
          :model-value="props.modelValue.Left"
          @update:model-value="onChangeLeft"
        ></ProofAssertionVue
      ></v-col>
      <v-col md="1" align-self="center" class="px-0">
        <v-btn icon size="x-small" @click="deleteLeft()">
          <v-icon icon="mdi-arrow-left" color="red"></v-icon>
        </v-btn>
      </v-col>
    </v-row>

    <v-row>
      <v-spacer></v-spacer>
      <v-col>
        <v-select
          class="mt-3 mb-2"
          density="compact"
          variant="outlined"
          label="Connecteur"
          hide-details
          :model-value="BinaryLabels[props.modelValue.Op]"
          @update:model-value="onChangeBinary"
          :items="binaryItems"
        ></v-select
      ></v-col>
      <v-spacer></v-spacer>
    </v-row>

    <v-row>
      <v-col cols="11">
        <ProofAssertionVue
          :model-value="props.modelValue.Right"
          @update:model-value="onChangeRight"
        ></ProofAssertionVue
      ></v-col>
      <v-col md="1" align-self="center" class="px-0">
        <v-btn icon size="x-small" @click="deleteRight()">
          <v-icon icon="mdi-arrow-left" color="red"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
  </div>
</template>

<script setup lang="ts">
import {
  Binary,
  BinaryLabels,
  type ProofAssertion,
  type ProofNode,
} from "@/controller/api_gen";
import { emptyAssertion } from "@/controller/editor";
import ProofAssertionVue from "./ProofAssertion.vue";

interface Props {
  modelValue: ProofNode;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: ProofNode): void;
}>();

const binaryItems = [BinaryLabels[Binary.And], BinaryLabels[Binary.Or]];

function onChangeBinary(s: string) {
  const op = Number(
    Object.entries(BinaryLabels).find((v) => v[1] == s)?.[0]
  ) as Binary;
  props.modelValue.Op = op;
  emit("update:modelValue", props.modelValue);
}
function onChangeLeft(v: ProofAssertion) {
  props.modelValue.Left = v;
  emit("update:modelValue", props.modelValue);
}
function onChangeRight(v: ProofAssertion) {
  props.modelValue.Right = v;
  emit("update:modelValue", props.modelValue);
}

function deleteLeft() {
  props.modelValue.Left = emptyAssertion();
  emit("update:modelValue", props.modelValue);
}
function deleteRight() {
  props.modelValue.Right = emptyAssertion();
  emit("update:modelValue", props.modelValue);
}
</script>

<style></style>
