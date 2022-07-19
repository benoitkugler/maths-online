<template>
  <v-form class="mx-4 my-2">
    <int-list-field
      label="Taille des groupes"
      :model-value="props.modelValue || []"
      @update:model-value="(v) => emit('update:model-value', v)"
      :sorted="false"
      disallow-intervals
    ></int-list-field>
  </v-form>
</template>

<script setup lang="ts">
import IntListField from "../editor/utils/IntListField.vue";

interface Props {
  modelValue: number[];
}

const emit = defineEmits<{
  (e: "update:model-value", v: number[]): void;
}>();

const props = defineProps<Props>();
</script>

<style scoped></style>
