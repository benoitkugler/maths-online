<template>
  <v-row no-gutters class="my-1">
    <v-col cols="12">
      <div
        :style="{
          margin: 'auto',
          display: 'block',
          width: '60px',
          height: '60px',
          'background-image': gradientSpec,
          border: '1px solid black',
          'border-radius': '50%'
        }"
      ></div>
    </v-col>
    <v-col cols="12" style="text-align: center" class="mt-1">
      <v-chip
        size="small"
        :color="props.highlight ? 'yellow-darken-3' : 'grey'"
        >{{ props.label }}</v-chip
      >
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { colorsPerCategorie } from "@/controller/trivial";
import { computed } from "@vue/runtime-core";

interface Props {
  success: boolean[];
  label: string;
  highlight: boolean;
}

const props = defineProps<Props>();

const gradientSpec = computed(() => {
  const args = colorsPerCategorie
    .map((color, index) => {
      const c = props.success[index] ? color : "white";
      const angle = ((index + 1) * 360) / colorsPerCategorie.length - 1;
      return `${c} 0 ${angle}deg, black ${angle}deg ${angle + 2}deg`;
    })
    .join(", ");
  return `conic-gradient(${args})`;
});
</script>

<style scoped></style>
