<template>
  <v-sheet
    class="rounded pa-1"
    :style="props.isRoot ? '' : 'border: 1px solid grey'"
  >
    <div v-for="(part, index) in props.modelValue.Parts" :key="index">
      <v-row>
        <v-col cols="11">
          <ProofAssertionVue
            :model-value="part"
            @update:model-value="(v) => onChange(index, v)"
          >
          </ProofAssertionVue>
        </v-col>
        <v-col md="1" align-self="center" class="px-0">
          <v-btn icon size="x-small" @click="deletePart(index)">
            <v-icon
              :icon="isInvalidAt(index) ? 'mdi-delete' : 'mdi-close'"
              color="red"
            ></v-icon>
          </v-btn>
        </v-col>
      </v-row>

      <div
        v-if="index != (props.modelValue.Parts?.length || 0) - 1"
        class="my-1"
      >
        donc
      </div>
    </div>
    <v-btn block size="small" class="my-2" @click="append">
      <v-icon icon="mdi-plus" color="green"></v-icon>
      Ajouter un bloc
    </v-btn>
  </v-sheet>
</template>

<script setup lang="ts">
import {
  ProofAssertionKind,
  type ProofAssertion,
  type ProofSequence,
} from "@/controller/api_gen";
import { emptyAssertion } from "@/controller/editor";
import ProofAssertionVue from "./ProofAssertion.vue";

interface Props {
  modelValue: ProofSequence;
  isRoot?: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: ProofSequence): void;
}>();

function onChange(index: number, value: ProofAssertion) {
  props.modelValue.Parts![index] = value;
  emit("update:modelValue", props.modelValue);
}

function isInvalidAt(index: number) {
  const kind = props.modelValue.Parts![index].Kind;
  return kind == ProofAssertionKind.ProofInvalid;
}

function deletePart(index: number) {
  const kind = props.modelValue.Parts![index].Kind;
  if (isInvalidAt(index)) {
    props.modelValue.Parts!.splice(index, 1);
  } else {
    props.modelValue.Parts![index] = emptyAssertion();
  }
  emit("update:modelValue", props.modelValue);
}
function append() {
  const l = props.modelValue.Parts || [];
  l.push(emptyAssertion());
  props.modelValue.Parts = l;
  emit("update:modelValue", props.modelValue);
}
</script>

<style></style>
