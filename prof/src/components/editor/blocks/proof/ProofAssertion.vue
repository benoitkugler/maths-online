<template>
  <ProofInvalidVue v-if="isInvalid" @add="(v) => emit('update:modelValue', v)">
  </ProofInvalidVue>
  <component
    v-else
    :model-value="props.modelValue.Data"
    @update:model-value="(v: any) => emit('update:modelValue', {Kind:props.modelValue.Kind,Data:v})"
    :is="comp"
  ></component>
</template>

<script setup lang="ts">
import { ProofAssertionKind, type ProofAssertion } from "@/controller/api_gen";
import {
  computed,
  defineAsyncComponent,
  markRaw,
  type DefineComponent,
} from "vue";
import ProofInvalidVue from "./ProofInvalid.vue";

const ProofStatementVue = defineAsyncComponent<DefineComponent>(
  () => import("./ProofStatement.vue") as any
);
const ProofEqualityVue = defineAsyncComponent<DefineComponent>(
  () => import("./ProofEquality.vue") as any
);
const ProofNodeVue = defineAsyncComponent<DefineComponent>(
  () => import("./ProofNode.vue") as any
);
const ProofSequenceVue = defineAsyncComponent<DefineComponent>(
  () => import("./ProofSequence.vue") as any
);

interface Props {
  modelValue: ProofAssertion;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: ProofAssertion): void;
}>();

const isInvalid = computed(
  () => props.modelValue.Kind == ProofAssertionKind.ProofInvalid
);

const comp = computed(() => {
  switch (props.modelValue.Kind) {
    case ProofAssertionKind.ProofStatement:
      return markRaw(ProofStatementVue);
    case ProofAssertionKind.ProofEquality:
      return markRaw(ProofEqualityVue);
    case ProofAssertionKind.ProofNode:
      return markRaw(ProofNodeVue);
    case ProofAssertionKind.ProofSequence:
      return markRaw(ProofSequenceVue);
    default:
      throw "exhaustive assertion type switch";
  }
});
</script>

<style></style>
