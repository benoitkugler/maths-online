<template>
  <v-chip size="small" label class="mx-1" :color="tagColor(tag)">{{
    tag.toUpperCase()
  }}</v-chip>
</template>

<script setup lang="ts">
import { tagColor } from "@/controller/editor";
interface Props {
  tag: string;
}
const props = defineProps<Props>();
</script>
