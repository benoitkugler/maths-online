<template>
  <v-card> {{ props.data }} </v-card>
</template>

<script setup lang="ts">
import type { TargetTrivial } from "@/controller/api_gen";

interface Props {
  data: TargetTrivial;
}

const props = defineProps<Props>();
</script>
