<template>
  <PreviewSheet :sheet="props.data.Sheet"></PreviewSheet>
</template>

<script setup lang="ts">
import type { TargetSheet } from "@/controller/api_gen";
import PreviewSheet from "../homework/PreviewSheet.vue";

interface Props {
  data: TargetSheet;
}

const props = defineProps<Props>();
</script>
