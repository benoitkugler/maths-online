<template>
  <v-card
    width="1000px"
    title="Syntaxe LaTeX"
    subtitle="La syntaxe LaTeX permet d'utiliser entre autres les commandes suivantes."
  >
    <v-card-text>
      <v-list color="info" rounded>
        <v-list-item
          v-for="(command, index) in commands"
          :key="index"
          :class="{ rounded: true, 'bg-grey-lighten-4': index % 2 == 0 }"
        >
          <v-row>
            <v-col cols="2" align-self="center" style="text-align: center">
              <div v-html="commandToHTML(command.command)"></div>
            </v-col>
            <v-col cols="6" align-self="center">
              <v-list-item-title>{{ command.description }}</v-list-item-title>
            </v-col>
            <v-col align-self="center" class="text-grey">
              <div v-html="command.command.replace('\n', '<br />')"></div>
            </v-col>
            <v-col cols="1" align-self="center">
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
</template>

<script setup lang="ts">
import katex from "katex";
import "katex/dist/katex.min.css";

const emit = defineEmits<{
  (e: "close"): void;
}>();

// interface ContainerProps {

// }

// const props = defineProps<ContainerProps>();

function commandToHTML(latex: string) {
  return katex.renderToString(latex);
}

const commands = [
  { description: "Vecteur", command: "\\overrightarrow{AB}" },
  { description: "Accolade ouverte", command: "\\{" },
  {
    description: "Coordonées d'un vecteur (colonne)",
    command: `\\begin{pmatrix} x \\\\ y \\end{pmatrix}`,
  },
  {
    description: "Fraction",
    command: "\\frac{a}{10^p}",
  },
  {
    description: "Ensemble des naturels",
    command: `\\mathbb{N} `,
  },
  {
    description: "Ensemble des relatifs",
    command: `\\mathbb{Z} `,
  },
  {
    description: "Ensemble des décimaux",
    command: `\\mathbb{D} `,
  },
  {
    description: "Ensemble des rationnels",
    command: `\\mathbb{Q} `,
  },
  {
    description: "Ensemble des réels",
    command: `\\mathbb{R} `,
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
    description: "Appartient",
    command: `\\in`,
  },
  {
    description: "N'appartient pas",
    command: `\\notin`,
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
    description: "Environ égal",
    command: `\\approx `,
  },
  {
    description: "Différent de",
    command: `\\neq`,
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
  await navigator.clipboard.writeText(command);
  emit("close");
}
</script>

<style scoped></style>