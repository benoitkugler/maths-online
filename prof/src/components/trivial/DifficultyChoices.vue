<template>
  <v-card title="Ajuster les étiquettes de difficulté.">
    <v-card-text>
      <v-checkbox
        hide-details
        :label="DifficultyTag.Diff1"
        v-model="diffChoices[DifficultyTag.Diff1]"
      >
      </v-checkbox>
      <v-checkbox
        hide-details
        :label="DifficultyTag.Diff2"
        v-model="diffChoices[DifficultyTag.Diff2]"
      >
      </v-checkbox>
      <v-checkbox
        hide-details
        :label="DifficultyTag.Diff3"
        v-model="diffChoices[DifficultyTag.Diff3]"
      >
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

const diffChoices = ref<{ [key in DifficultyTag]: boolean }>(buildCrible());

watch(props, () => (diffChoices.value = buildCrible()));

function buildCrible() {
  return {
    [DifficultyTag.DiffEmpty]: false, // not used
    [DifficultyTag.Diff1]: props.difficulties.includes(DifficultyTag.Diff1),
    [DifficultyTag.Diff2]: props.difficulties.includes(DifficultyTag.Diff2),
    [DifficultyTag.Diff3]: props.difficulties.includes(DifficultyTag.Diff3),
  };
}

function apply() {
  const diffs = Object.entries(diffChoices)
    .filter((e) => e[1])
    .map((e) => e[0] as DifficultyTag);
  emit("update:difficulties", diffs);
}
</script>

<style scoped></style>
