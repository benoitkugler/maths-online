<template>
  <v-row dense class="mt-1" style="text-align: center">
    <v-col>
      <v-btn size="small" @click="addStatement"> Champ </v-btn>
    </v-col>
    <v-col>
      <v-btn size="small" @click="addEquality"> Egalité </v-btn>
    </v-col>
    <v-col>
      <v-btn size="small" @click="addNode"> Opérateur logique </v-btn>
    </v-col>
    <v-col>
      <v-btn size="small" @click="addSequence"> Séquence </v-btn>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import {
  Binary,
  ProofAssertionKind,
  type ProofAssertion,
  type ProofEquality,
  type ProofNode,
  type ProofSequence,
  type ProofStatement,
} from "@/controller/api_gen";
import { emptyAssertion } from "@/controller/editor";

// interface Props {}
// const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "add", value: ProofAssertion): void;
}>();

function addStatement() {
  const data: ProofStatement = { Content: "" };
  emit("add", {
    Kind: ProofAssertionKind.ProofStatement,
    Data: data,
  });
}

function addEquality() {
  const data: ProofEquality = { Terms: "m + n = 2k + 2k' = 2(k+k')" };
  emit("add", {
    Kind: ProofAssertionKind.ProofEquality,
    Data: data,
  });
}

function addNode() {
  const data: ProofNode = {
    Op: Binary.And,
    Left: emptyAssertion(),
    Right: emptyAssertion(),
  };
  emit("add", {
    Kind: ProofAssertionKind.ProofNode,
    Data: data,
  });
}
function addSequence() {
  const data: ProofSequence = {
    Parts: [emptyAssertion(), emptyAssertion()],
  };
  emit("add", {
    Kind: ProofAssertionKind.ProofSequence,
    Data: data,
  });
}
</script>

<style></style>
