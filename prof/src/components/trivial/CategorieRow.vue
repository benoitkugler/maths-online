<template>
  <v-list-item
    rounded
    :style="{
      'border-color': colors[props.index],
      borderWidth: '2px',
      borderStyle: 'solid',
    }"
    class="my-2"
  >
    <slot></slot>
  </v-list-item>
</template>

<script setup lang="ts">
import { colorsPerCategorie } from "@/controller/trivial";

interface Props {
  index: number;
}

const props = defineProps<Props>();
const colors = colorsPerCategorie;
</script>

<style scoped></style>
