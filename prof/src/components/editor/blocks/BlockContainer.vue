<template>
  <v-dialog v-model="showDocumentation">
    <v-card
      title="Syntaxe LaTeX"
      subtitle="La syntaxe LaTeX permet d'utiliser entre autres les commandes suivantes."
    >
      <v-card-text>
        <v-list color="info" rounded>
          <v-list-item v-for="(command, index) in commands" :key="index">
            <v-row>
              <v-col cols="6">
                <v-list-item-title>{{ command.description }}</v-list-item-title>
              </v-col>
              <v-col align-self="center" class="text-grey">
                <div v-html="command.command.replace('\n', '<br />')"></div>
              </v-col>
              <v-col cols="2" align-self="center">
                <v-btn icon size="x-small" title="Copier">
                  <v-icon
                    icon="mdi-content-copy"
                    size="x-small"
                    @click="copyAndClose(command.command)"
                  ></v-icon>
                </v-btn>
              </v-col>
            </v-row>
          </v-list-item>
        </v-list>
      </v-card-text>
    </v-card>
  </v-dialog>

  <v-card class="my-2" elevation="3">
    <v-row
      no-gutters
      :class="'px-2 rounded ' + colorClass"
      @dragstart="onDragStart"
      draggable="true"
      style="cursor: grab"
    >
      <v-col align-self="center" cols="8">
        <v-card-subtitle>
          <b>{{ kindLabels[props.kind].label }}</b>
          <span v-if="isAnswer" class="ml-1">(Champ de réponse)</span>
        </v-card-subtitle>
      </v-col>
      <v-col cols="4" style="text-align: right" class="my-2">
        <v-btn
          v-if="showLaTeXDoc"
          class="mr-2"
          icon
          title="Documentation de la syntaxe LaTeX"
          size="x-small"
        >
          <v-icon small color="info" @click="showDocumentation = true"
            >mdi-help</v-icon
          >
        </v-btn>
        <v-btn icon title="Supprimer" size="x-small">
          <v-icon small color="red" @click="emit('delete')">mdi-close</v-icon>
        </v-btn>
      </v-col>
    </v-row>

    <v-card-text class="pt-1 pb-2" :hidden="hideContent">
      <slot></slot>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { BlockKindLabels, onDragListItemStart } from "@/controller/editor";
import { BlockKind } from "@/controller/exercice_gen";
import { computed } from "@vue/runtime-core";
import { $ref } from "vue/macros";

const emit = defineEmits<{
  (e: "delete"): void;
}>();

interface ContainerProps {
  index: number;
  kind: BlockKind;
  hideContent: boolean;
  hasError: boolean;
}

const props = defineProps<ContainerProps>();

const kindLabels = BlockKindLabels;

const isAnswer = computed(() => kindLabels[props.kind].isAnswerField);

const colorClass = computed(() => {
  if (props.hasError) {
    return "bg-red";
  }
  return isAnswer.value ? "bg-pink-lighten-3" : "bg-purple-lighten-3";
});

const showLaTeXDoc = computed(() => {
  return ![BlockKind.ExpressionFieldBlock, BlockKind.NumberFieldBlock].includes(
    props.kind
  );
});

function onDragStart(payload: DragEvent) {
  onDragListItemStart(payload, props.index);
}

let showDocumentation = $ref(false);
const commands = [
  { description: "Vecteur", command: "\\overrightarrow{AB}" },
  { description: "Accolade ouverte", command: "\\{" },
  {
    description: "Coordonées d'un vecteur (colonne)",
    command: `\\begin{pmatrix} x \\\\ y \\end{pmatrix}`,
  },
  {
    description: "Ensemble des réels",
    command: `\\R`,
  },
  {
    description: "Ensemble vide",
    command: `\\empty`,
  },
  {
    description: "Inclus",
    command: `\\subset`,
  },
  {
    description: "Non inclus",
    command: `\\not\\subset`,
  },
  {
    description: "Antislash",
    command: `\\backslash`,
  },
  {
    description: "Infini",
    command: `\\infty`,
  },
  {
    description: "Inférieur ou égal (<=)",
    command: `\\leq `,
  },
  {
    description: "Supérieur ou égal (>=)",
    command: `\\geq `,
  },
  {
    description: "Union",
    command: `\\cup `,
  },
  {
    description: "Intersection",
    command: `\\cap `,
  },
];

async function copyAndClose(command: string) {
  showDocumentation = false;
  await navigator.clipboard.writeText(command);
}
</script>

<style scoped>
.small-slider:deep(.v-input__control) {
  min-height: 200px !important;
}
</style>
