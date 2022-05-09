<template>
  <v-chip size="x-small" label class="ml-1" :color="tagColor(tag)">{{
    tagString(tag)
  }}</v-chip>
</template>

<script setup lang="ts">
import { tagColor, tagString } from "@/controller/editor";
interface Props {
  tag: string;
}
const props = defineProps<Props>();
</script>
