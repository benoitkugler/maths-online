<template>
  <v-row no-gutters>
    <v-col cols="auto" align-self="center">
      <v-chip class="mx-1" color="light-green">
        {{ props.success }}
      </v-chip>
    </v-col>
    <v-col align-self="center">
      <v-progress-linear
        color="light-green"
        bg-color="orange-lighten-2"
        bg-opacity="0.8"
        :model-value="value"
        rounded
        :height="height"
      >
        <template v-slot:default="{ value }">
          <strong>{{ Math.ceil(value) }}%</strong>
        </template>
      </v-progress-linear>
    </v-col>
    <v-col cols="auto" align-self="center">
      <v-chip class="mx-1" color="orange-lighten-2">
        {{ props.failure }}
      </v-chip>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { computed } from "vue";

interface Props {
  success: number;
  failure: number;
  height?: number;
}

const props = defineProps<Props>();

const value = computed(
  () => (100 * props.success) / (props.success + props.failure)
);
</script>

<style></style>
