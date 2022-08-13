<template>
  <v-menu offset-y :close-on-content-click="false">
    <template v-slot:activator="{ isActive, props: slotProps }">
      <v-btn
        title="Modifier la couleur"
        v-on="{ isActive }"
        v-bind="slotProps"
        :color="btnParams.color"
        block
      >
        {{ btnParams.text }}
      </v-btn>
    </template>
    <v-color-picker
      hide-inputs
      mode="hex"
      :swatches="swatches"
      show-swatches
      swatches-max-height="200"
      :model-value="props.modelValue"
      @update:model-value="onPick"
    ></v-color-picker>
  </v-menu>
</template>

<script setup lang="ts">
import { lastColorUsed } from "@/controller/editor";
import { computed } from "vue";

interface Props {
  modelValue: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: string): void;
}>();

const swatches = [
  ["#00000000", "#FF0000", "#AA0000", "#550000"],
  ["#000000", "#FFFF00", "#AAAA00", "#555500"],
  ["#fc03df", "#00FF00", "#00AA00", "#005500"],
  ["#852179", "#00FFFF", "#00AAAA", "#005555"],
  ["#fcba03", "#0000FF", "#0000AA", "#000055"],
];

const btnParams = computed(() => {
  if (props.modelValue == "#00000000") {
    return { text: "Aucune", color: "#FFFFFF" };
  } else {
    return { text: "", color: props.modelValue };
  }
});

function onPick(color: string) {
  lastColorUsed.color = color;
  emit("update:model-value", color);
}
</script>

<style scoped>
.centered-input:deep(input) {
  text-align: center;
}

:deep(.v-field__append-inner) {
  padding-top: 4px;
}
</style>
