<template>
  <v-row no-gutters class="px-2">
    <v-col>
      <v-btn icon flat title="Supprimer" small>
        <v-icon small color="red" @click="emit('delete')">mdi-close</v-icon>
      </v-btn>

      <v-menu
        offset-y
        close-on-content-click
        v-if="props.nbBlocks >= 2"
        @update:model-value="onClosePositionner"
      >
        <template v-slot:activator="{ isActive, props }">
          <v-btn
            icon
            flat
            title="Positionner"
            size="small"
            v-on="{ isActive }"
            v-bind="props"
          >
            <v-icon small color="secondary">mdi-sort</v-icon>
          </v-btn>
        </template>
        <v-card>
          <v-card-text>
            <v-slider
              class="small-slider"
              direction="vertical"
              step="1"
              :min="0"
              :max="props.nbBlocks - 1"
              show-ticks="always"
              :tick-size="6"
              v-model="initialIndex"
              hide-details
            ></v-slider>
          </v-card-text>
        </v-card>
      </v-menu>
    </v-col>
    <v-col md="auto" style="text-align: right">
      <slot name="toolbar"></slot>
    </v-col>
  </v-row>
  <div class="px-4 my-0 py-0">
    <slot></slot>
  </div>
</template>

<script setup lang="ts">
import { ref } from "@vue/reactivity";
import { watch } from "@vue/runtime-core";

const emit = defineEmits<{
  (e: "delete"): void;
  (e: "swap", origin: number, target: number): void;
}>();

export interface ContainerProps {
  index: number;
  nbBlocks: number;
}

const props = defineProps<ContainerProps>();

let initialIndex = ref(props.nbBlocks - 1 - props.index);

watch(props, () => (initialIndex.value = props.nbBlocks - 1 - props.index));

function onClosePositionner(isOpen: boolean) {
  if (!isOpen) {
    // commit the changes
    console.log(props.index, initialIndex.value);

    emit("swap", props.index, props.nbBlocks - 1 - initialIndex.value);
  }
}
</script>

<style scoped>
.small-slider:deep(.v-input__control) {
  min-height: 200px !important;
}
</style>
