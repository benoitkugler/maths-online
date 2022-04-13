<template>
  <v-menu offset-y close-on-content-click>
    <template v-slot:activator="{ isActive, props: slotProps }">
      <v-btn
        title="Modifier la couleur"
        v-on="{ isActive }"
        v-bind="slotProps"
        :color="props.modelValue"
      >
        Couleur
      </v-btn>
    </template>
    <v-color-picker
      hide-inputs
      mode="hex"
      :swatches="swatches"
      show-swatches
      :model-value="props.modelValue"
      @update:model-value="s => emit('update:model-value', s)"
    ></v-color-picker>
  </v-menu>
</template>

<script setup lang="ts">
interface Props {
  modelValue: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: string): void;
}>();

const swatches = [
  ["#FF0000", "#AA0000", "#550000"],
  ["#FFFF00", "#AAAA00", "#555500"],
  ["#00FF00", "#00AA00", "#005500"],
  ["#00FFFF", "#00AAAA", "#005555"],
  ["#0000FF", "#0000AA", "#000055"]
];
</script>

<style scoped>
.centered-input:deep(input) {
  text-align: center;
}

:deep(.v-field__append-inner) {
  padding-top: 4px;
}
</style>
