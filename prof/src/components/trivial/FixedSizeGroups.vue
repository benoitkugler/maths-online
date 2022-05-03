<template>
  <v-form class="mx-4 my-2">
    <int-list-field
      label="Taille des groupes"
      :model-value="props.modelValue.Groups || []"
      @update:model-value="v => emit('update:model-value', { Groups: v })"
      :sorted="false"
      disallow-intervals
    ></int-list-field>
  </v-form>
</template>

<script setup lang="ts">
import type { FixedSizeGroupStrategy } from "@/controller/api_gen";
import IntListField from "../editor/blocks/IntListField.vue";

interface Props {
  modelValue: FixedSizeGroupStrategy;
}

const emit = defineEmits<{
  (e: "update:model-value", v: FixedSizeGroupStrategy): void;
}>();

const props = defineProps<Props>();
</script>

<style scoped></style>
