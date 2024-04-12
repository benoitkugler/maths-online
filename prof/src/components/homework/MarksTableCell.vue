<template>
  <v-tooltip>
    <template v-slot:activator="{ isActive, props: innerProps }">
      <span v-on="{ isActive }" v-bind="innerProps" :style="{ color: color }">
        {{ formattedMark }}
      </span>
    </template>
    {{ props.data.NbTries }} essais
  </v-tooltip>
</template>

<script setup lang="ts">
import type { StudentTravailMark } from "@/controller/api_gen";
import { computed } from "vue";

interface Props {
  data: StudentTravailMark;
}

const props = defineProps<Props>();

const formattedMark = computed(() => {
  if (props.data.Dispensed) {
    return `${formatFloat(props.data.Mark)} (*)`;
  }
  return formatFloat(props.data.Mark);
});

const color = computed(() =>
  props.data.NbTries == 0 ? "red" : props.data.Mark == 0 ? "orange" : "black"
);

function formatFloat(v: number) {
  if (Math.round(v) == v) {
    return v.toFixed(0);
  }
  return v.toFixed(1);
}
</script>
