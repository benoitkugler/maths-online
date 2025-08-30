<template>
  <v-card title="Ajuster les étiquettes de difficulté.">
    <v-card-text>
      <v-checkbox hide-details :label="DifficultyTag.Diff1" v-model="d1">
      </v-checkbox>
      <v-checkbox hide-details :label="DifficultyTag.Diff2" v-model="d2">
      </v-checkbox>
      <v-checkbox hide-details :label="DifficultyTag.Diff3" v-model="d3">
      </v-checkbox>
    </v-card-text>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn color="success" @click="apply">Appliquer</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import { DifficultyTag } from "@/controller/api_gen";
import { watch } from "vue";
import { ref } from "vue";

interface Props {
  difficulties: DifficultyTag[];
}

const emit = defineEmits<{
  (e: "update:difficulties", v: DifficultyTag[]): void;
}>();

const props = defineProps<Props>();

const d1 = ref(props.difficulties.includes(DifficultyTag.Diff1));
const d2 = ref(props.difficulties.includes(DifficultyTag.Diff2));
const d3 = ref(props.difficulties.includes(DifficultyTag.Diff3));

watch(props, () => {
  d1.value = props.difficulties.includes(DifficultyTag.Diff1);
  d2.value = props.difficulties.includes(DifficultyTag.Diff2);
  d3.value = props.difficulties.includes(DifficultyTag.Diff3);
});

function apply() {
  const out: DifficultyTag[] = [];
  if (d1.value) out.push(DifficultyTag.Diff1);
  if (d2.value) out.push(DifficultyTag.Diff2);
  if (d3.value) out.push(DifficultyTag.Diff3);
  emit("update:difficulties", out);
}
</script>

<style scoped></style>
