<template>
  <v-card class="my-2" @dragstart="onDragStart" draggable="true">
    <v-row no-gutters class="px-2">
      <v-col align-self="center" cols="6">
        <v-card-subtitle>
          {{ kindLabels[props.kind] }}
        </v-card-subtitle>
      </v-col>
      <v-col cols="6" style="text-align: right">
        <!-- <v-btn icon flat title="Masquer" small>
          <v-icon
            small
            @click="hidden = !hidden"
            :icon="hidden ? 'mdi-eye' : 'mdi-eye-off'"
          ></v-icon>
        </v-btn> -->

        <v-btn icon flat title="Supprimer" small>
          <v-icon small color="red" @click="emit('delete')">mdi-close</v-icon>
        </v-btn>
      </v-col>
    </v-row>

    <v-card-text class="pt-0 pb-2" :hidden="hidden">
      <slot></slot>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { BlockKindLabels } from "@/controller/editor";
import type { BlockKind } from "@/controller/exercice_gen";
import { ref } from "@vue/reactivity";
import { watch } from "@vue/runtime-core";

const emit = defineEmits<{
  (e: "delete"): void;
  (e: "swap", origin: number, target: number): void;
  (e: "dragStart"): void;
}>();

const hidden = ref(false);

export interface ContainerProps {
  index: number;
  nbBlocks: number;
  kind: BlockKind;
}

const props = defineProps<ContainerProps>();

const kindLabels = BlockKindLabels;

let initialIndex = ref(props.nbBlocks - props.index);

watch(props, () => (initialIndex.value = props.nbBlocks - props.index));

function onClosePositionner(isOpen: boolean) {
  if (!isOpen) {
    // commit the changes
    console.log(props.index, initialIndex.value);

    emit("swap", props.index, props.nbBlocks - initialIndex.value);
  }
}

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
