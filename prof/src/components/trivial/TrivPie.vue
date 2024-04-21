<template>
  <v-row no-gutters class="my-1">
    <v-col cols="12">
      <div
        :style="{
          margin: 'auto',
          display: 'block',
          width: '50px',
          height: '50px',
          'background-image': gradientSpec,
          border: '1px solid black',
          'border-radius': '50%',
        }"
      ></div>
    </v-col>
    <v-col
      cols="12"
      style="text-align: center"
      align-self="center"
      class="mt-2"
    >
      <v-chip
        size="small"
        :color="props.highlight ? 'yellow-darken-3' : 'grey'"
        >{{ props.label }}</v-chip
      >
      <v-tooltip
        v-if="props.isWaiting"
        text="Ce joueur est occupé par la question"
      >
        <template v-slot:activator="{ props }">
          <v-icon v-bind="props" size="24" color="info" class="mx-1"
            >mdi-square-edit-outline</v-icon
          >
        </template>
      </v-tooltip>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { colorsPerCategorie } from "@/controller/trivial";
import { computed } from "vue";

interface Props {
  success: boolean[];
  label: string;
  highlight: boolean;
  isWaiting: boolean;
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
