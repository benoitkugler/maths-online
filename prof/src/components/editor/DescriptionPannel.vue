<template>
  <v-card
    width="600px"
    title="Description"
    subtitle="Note interne affichée aux professeurs."
  >
    <v-card-text>
      <v-textarea
        @update:model-value="(s) => emit('update:model-value', s)"
        :model-value="props.modelValue"
        placeholder="Commentaire..."
        :readonly="props.readonly"
      ></v-textarea>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
interface Props {
  modelValue: string;
  readonly: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", description: string): void;
}>();
</script>

<style scoped></style>
