<template>
  <v-row style="flex-wrap: nowrap">
    <v-col
      class="mt-2 mb-1"
      align-self="center"
      style="min-width: fit-content"
      v-if="!props.isRoot"
    >
      <v-select
        variant="outlined"
        density="compact"
        label="Evènement"
        :items="items"
        v-model="props.modelValue.Value"
        hide-details
      ></v-select>
    </v-col>
    <v-col align-self="center" style="min-width: fit-content">
      <div v-for="(item, index) in props.modelValue.Children" :key="index">
        <v-divider v-if="index != 0"></v-divider>

        <div style="display: flex; flex-direction: row; flex-wrap: nowrap">
          <v-col align-self="center" style="min-width: 150px">
            <v-text-field
              variant="outlined"
              density="compact"
              v-model="props.modelValue.Probabilities![index]"
              hide-details
              label="Probabilité"
              :color="ExpressionColor"
              class="mr-2"
            >
            </v-text-field>
          </v-col>
          <TreeNode
            v-model="props.modelValue.Children![index]"
            @update:model-value="emit('update:modelValue', props.modelValue)"
            :event-proposals="props.eventProposals"
            :is-root="false"
          ></TreeNode>
        </div>
      </div>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { TreeNodeAnswer } from "@/controller/api_gen";
import { ExpressionColor } from "@/controller/editor";
import { computed } from "@vue/runtime-core";

interface Props {
  modelValue: TreeNodeAnswer;
  eventProposals: string[];
  isRoot: boolean; // if false, only shows the label
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: TreeNodeAnswer): void;
}>();

const items = computed(() =>
  (props.eventProposals || []).map((e, index) => ({ title: e, value: index }))
);
</script>

<style></style>
