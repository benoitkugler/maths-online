<template>
  <div>
    <v-row>
      <v-spacer></v-spacer>
      <v-col md="auto" style="text-align: right">
        <v-checkbox
          label="Comme conseil"
          :model-value="props.data.IsHint"
        ></v-checkbox>
      </v-col>
    </v-row>
    <v-textarea :model-value="props.data.Parts"> </v-textarea>
  </div>
</template>

<script setup lang="ts">
import type { TextBlock } from "@/controller/exercice_gen";
export interface TextProps {
  data: TextBlock;
}

const props = defineProps<TextProps>();
</script>

<style scoped></style>
