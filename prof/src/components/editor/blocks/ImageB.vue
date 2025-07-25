<template>
  <v-row class="my-1">
    <v-col cols="8">
      <v-text-field
        density="compact"
        variant="outlined"
        label="URL"
        :model-value="props.modelValue.URL"
        @update:model-value="
          (v) =>
            emit('update:modelValue', { URL: v, Scale: props.modelValue.Scale })
        "
        hide-details
      ></v-text-field>
    </v-col>
    <v-col>
      <v-text-field
        type="number"
        step="1"
        min="1"
        max="100"
        suffix="%"
        label="Echelle"
        variant="outlined"
        density="compact"
        :model-value="props.modelValue.Scale"
        @update:model-value="
          (v) =>
            emit('update:modelValue', {
              URL: props.modelValue.URL,
              Scale: Number(v) as Int,
            })
        "
      ></v-text-field>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { ImageBlock, Int, Variable } from "@/controller/api_gen";

interface Props {
  modelValue: ImageBlock;
  availableParameters: Variable[];
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: ImageBlock): void;
}>();
</script>

<style></style>
