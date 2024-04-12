<template>
  <v-card
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
              <v-btn
                icon
                size="x-small"
                title="Copier"
                class="my-1"
                @click="copyAndClose(command.command)"
              >
                <v-icon icon="mdi-content-copy" size="x-small"></v-icon>
              </v-btn>
            </v-col>
          </v-row>
        </v-list-item>
      </v-list>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { controller } from "@/controller/controller";
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
    description: "Inférieur strict (<)",
    command: `\\lt `,
  },
  {
    description: "Supérieur strict (>)",
    command: `\\gt `,
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
  {
    description: "Angle",
    command: `\\widehat{BAC}`,
  },
  {
    description: "Limite",
    command: `\\lim_{x \\to 0} f(x)`,
  },
  {
    description: "Barre (moyenne, contraire)",
    command: `\\overline{x}`,
  },
  {
    description: "Symbole parallèle",
    command: `\\parallel`,
  },
  {
    description: "Symbole non parallèle",
    command: `\\nparallel`,
  },
  {
    description: "Symbole équivalent",
    command: `\\iff`,
  },
  {
    description: "Système de deux équations",
    command: `\\begin{cases}  x+y &=2 \\\\   x-3y &=4 \\\\ \\end{cases}`,
  },
  {
    description: "Flèche à l'envers (attribution d'une variable)",
    command: `\\leftarrow`,
  },
];

async function copyAndClose(command: string) {
  await navigator.clipboard.writeText(command);
  if (controller.showMessage)
    controller.showMessage("Commande copiée dans le presse-papier.");
  emit("close");
}
</script>

<style scoped></style>
