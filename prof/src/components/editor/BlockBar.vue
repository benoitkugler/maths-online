<template>
  <v-list style="max-height: 70vh">
    <v-list-subheader><h3 class="text-purple">Enoncés</h3></v-list-subheader>
    <v-list-item
      rounded
      dense
      class="py-0 bg-purple-lighten-4 ma-1"
      v-for="kind in staticKinds"
      link
      @click="emit('add', kind)"
      :key="kind"
    >
      {{ labels[kind].label }}
    </v-list-item>

    <template v-if="!props.hideAnswerFields">
      <v-divider></v-divider>
      <v-list-subheader
        ><h3 class="text-pink">Champs de réponse</h3></v-list-subheader
      >
      <v-list-item
        rounded
        dense
        class="py-0 bg-pink-lighten-4 ma-1"
        v-for="kind in fieldKinds"
        link
        :key="kind"
        @click="emit('add', kind)"
      >
        {{ labels[kind].label }}
      </v-list-item>
    </template>
  </v-list>
</template>

<script setup lang="ts">
import { BlockKind } from "@/controller/api_gen";
import { BlockKindLabels, sortedBlockKindLabels } from "@/controller/editor";

interface Props {
  simplified: boolean;
  hideAnswerFields: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "add", kind: BlockKind): void;
}>();

const labels = BlockKindLabels;

function isSimplified(k: BlockKind) {
  return (
    k == BlockKind.TextBlock ||
    k == BlockKind.TableBlock ||
    k == BlockKind.OrderedListFieldBlock ||
    k == BlockKind.RadioFieldBlock
  );
}

const staticKinds = sortedBlockKindLabels
  .filter((k) => !props.simplified || isSimplified(k[0]))
  .filter((k) => !k[1].isAnswerField)
  .map((k) => k[0]);
const fieldKinds = sortedBlockKindLabels
  .filter((k) => !props.simplified || isSimplified(k[0]))
  .filter((k) => k[1].isAnswerField)
  .map((k) => k[0]);
</script>
