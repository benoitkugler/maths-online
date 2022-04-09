<template>
  <v-card class="my-2" elevation="3">
    <v-row
      no-gutters
      class="px-2 bg-purple-lighten-3"
      @dragstart="onDragStart"
      draggable="true"
    >
      <v-col align-self="center" cols="6">
        <v-card-subtitle>
          {{ kindLabels[props.kind] }}
        </v-card-subtitle>
      </v-col>
      <v-col cols="6" style="text-align: right" class="my-2">
        <!-- <v-btn icon flat title="Masquer" small>
          <v-icon
            small
            @click="hidden = !hidden"
            :icon="hidden ? 'mdi-eye' : 'mdi-eye-off'"
          ></v-icon>
        </v-btn> -->

        <v-btn icon flat title="Supprimer" size="small">
          <v-icon small color="red" @click="emit('delete')">mdi-close</v-icon>
        </v-btn>
      </v-col>
    </v-row>

    <v-card-text class="pt-1 pb-2">
      <slot></slot>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { BlockKindLabels } from "@/controller/editor";
import type { BlockKind } from "@/controller/exercice_gen";

const emit = defineEmits<{
  (e: "delete"): void;
  (e: "swap", origin: number, target: number): void;
  (e: "dragStart"): void;
}>();

export interface ContainerProps {
  index: number;
  kind: BlockKind;
}

const props = defineProps<ContainerProps>();

const kindLabels = BlockKindLabels;

function onDragStart(payload: DragEvent) {
  payload.dataTransfer?.setData(
    "text/json",
    JSON.stringify({ index: props.index })
  );
  payload.dataTransfer!.dropEffect = "move";
}
</script>

<style scoped>
.small-slider:deep(.v-input__control) {
  min-height: 200px !important;
}
</style>
