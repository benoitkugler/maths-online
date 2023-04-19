<template>
  <v-menu>
    <template v-slot:activator="{ isActive, props: innerP }">
      <v-chip v-bind="innerP" v-on="{ isActive }" color="info" size="small">
        {{ props.variants.length || 0 }}</v-chip
      >
    </template>
    <v-card max-width="500" subtitle="Variantes" class="px-0">
      <v-card-text class="pa-0">
        <v-list max-height="60vh" class="py-0 overflow-y-auto">
          <v-list-item
            density="compact"
            rounded
            class="my-1"
            v-for="(question, index) in props.variants"
            :key="index"
          >
            <v-row no-gutters>
              <v-col align-self="center"> ({{ question.Id }}) </v-col>
              <v-col align-self="center" class="my-4 px-3" cols="auto">
                <template v-if="question.Subtitle">{{
                  question.Subtitle
                }}</template>
                <i v-else>(Sans titre)</i>
              </v-col>
              <v-col align-self="center" style="text-align: right">
                <TagChip
                  v-if="question.Difficulty"
                  :tag="{ Tag: question.Difficulty, Section: 0 }"
                ></TagChip>
                <v-chip v-else size="small" label title="Difficulté"
                  >Aucune</v-chip
                >
              </v-col>
              <v-col align-self="center">
                <v-chip
                  class="mx-2"
                  label
                  v-if="question.HasCorrection"
                  size="small"
                  color="green"
                  >Avec correction</v-chip
                >
              </v-col>
            </v-row>
          </v-list-item>
        </v-list>
      </v-card-text>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import type { VariantG } from "@/controller/editor";
import TagChip from "./utils/TagChip.vue";

interface Props {
  variants: VariantG[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: string[]): void;
}>();
</script>

<style scoped></style>
