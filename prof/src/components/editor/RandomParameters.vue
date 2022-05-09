<template>
  <v-dialog v-model="showHelp" min-width="1200px" width="max-content">
    <v-card
      title="Paramètres aléatoires"
      subtitle="Description"
      style="width: 800px; max-height: 70vh"
      class="overflow-y-auto"
    >
      <v-card-text>
        Les paramètres aléatoires sont des variables dont les valeurs sont
        générées à chaque fois qu'une question est posée à l'élève. <br />
        Leur définition se fait avec une syntaxe de type calculatrice. En
        particulier, les fonctions suivantes peuvent être utilisées.

        <v-list color="info" rounded>
          <v-list-item v-for="content in helpContent">
            <v-row>
              <v-col cols="6">
                <v-list-item-title> {{ content[0] }} </v-list-item-title>
              </v-col>
              <v-col align-self="center">
                <v-list-item-subtitle>
                  {{ content[1] }}
                </v-list-item-subtitle>
              </v-col>
            </v-row>
          </v-list-item>
        </v-list>
      </v-card-text>
    </v-card>
  </v-dialog>

  <v-card class="mb-2">
    <v-row
      :style="{
        'background-color': props.isValidated ? 'lightgreen' : 'lightgray'
      }"
      class="rounded"
      no-gutters
    >
      <v-col cols="5" align-self="center">
        <v-card-subtitle class="py-2">Paramètres aléatoires</v-card-subtitle>
      </v-col>
      <v-col align-self="center" style="text-align: right">
        <v-progress-circular
          v-show="isLoading"
          indeterminate
          class="mx-1"
        ></v-progress-circular>
        <v-btn
          icon
          @click="emit('add')"
          title="Ajouter un paramètre"
          size="x-small"
          class="mx-1 my-2"
        >
          <v-icon icon="mdi-plus" color="green" small></v-icon>
        </v-btn>
        <v-btn
          icon
          @click="showHelp = true"
          title="Aide"
          size="x-small"
          class="mx-1 my-2"
        >
          <v-icon icon="mdi-help" color="info"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-row no-gutters>
      <v-col>
        <v-list v-if="props.parameters?.length" @dragend="showDropZone = false">
          <drop-zone
            v-if="showDropZone"
            @drop="origin => emit('swap', origin, 0)"
          ></drop-zone>
          <div v-for="(param, index) in props.parameters">
            <v-list-item
              class="pr-0"
              @dragstart="e => onItemDragStart(e, index)"
              draggable="true"
            >
              <v-row no-gutters>
                <v-col cols="3">
                  <variable-field
                    v-model="param.variable"
                    @update:model-value="emit('update', index, param)"
                    @blur="emit('done')"
                  >
                  </variable-field>
                </v-col>
                <v-col cols="7">
                  <v-text-field
                    class="ml-2 small-input"
                    variant="underlined"
                    density="compact"
                    hide-details
                    :model-value="param.expression"
                    @update:model-value="s => onExpressionChange(s, index)"
                    @blur="emit('done')"
                    :color="expressionColor"
                  ></v-text-field>
                </v-col>
                <v-col cols="2">
                  <v-btn icon size="small" flat @click="emit('delete', index)">
                    <v-icon icon="mdi-delete" color="red"></v-icon>
                  </v-btn>
                </v-col>
              </v-row>
            </v-list-item>
            <drop-zone
              v-if="showDropZone"
              @drop="origin => emit('swap', origin, index + 1)"
            ></drop-zone>
          </div>
        </v-list>
      </v-col>
    </v-row>
  </v-card>
</template>

<script setup lang="ts">
import { ExpressionColor, onDragListItemStart } from "@/controller/editor";
import type {
  randomParameter,
  randomParameters
} from "@/controller/exercice_gen";
import { $ref } from "vue/macros";
import DropZone from "./DropZone.vue";
import VariableField from "./utils/VariableField.vue";

interface Props {
  parameters: randomParameters;
  isLoading: boolean;
  isValidated: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "add"): void;
  (e: "update", index: number, param: randomParameter): void;
  (e: "delete", index: number): void;
  (e: "swap", origin: number, target: number): void;
  (e: "done"): void;
}>();

const expressionColor = ExpressionColor;

let showHelp = $ref(false);

function onExpressionChange(s: string, index: number) {
  const param = props.parameters![index];
  param.expression = s;
  emit("update", index, param);
}

let showDropZone = $ref(false);
function onItemDragStart(payload: DragEvent, index: number) {
  onDragListItemStart(payload, index);
  showDropZone = true;
}

const helpContent = [
  [
    "randChoice(-4;12;99)",
    "Renvoie un nombre aléatoire parmi ceux proposés par l'utilisateur, ici {-4, 12, 99}."
  ],
  [
    "randLetter(A, B, C)",
    "Renvoie une variable dont le nom sera choisi parmi ceux proposés, ici {A, B, C}."
  ],
  ["randPrime(15;28)", "Renvoie un nombre premier entre 15 et 28 (inclus)."],
  [
    "randDecDen()",
    "Renvoie un entier aléatoire parmi 1, 2, 4, 5, 8, 10, 16, 20, 25, 40, 50, 80, 100 (diviser n'importe quel entier par l'un de ces nombres permettra d'obtenir un nombre décimal)"
  ],
  ["round(x; 3)", "Arrondi x à trois chiffres après la virgule"],
  ["isPrime(n)", "Renvoie 1 is n est un nombre premier, 0 sinon"],
  ["sgn(x)", "Renvoie le signe de x : 1 si x > 0, -1 si x < 0, 0 si x = 0"],
  ["isZero(x)", "Renvoie 1 si x vaut 0, 0 sinon"],
  ["exp(x)", "Fonction exponentielle"],
  ["ln(x)", "Fonction logarithme"],
  ["sin(x)", "Fonction sinus"],
  ["cos(x)", "Fonction cosinus"],
  ["tan(x)", "Fonction tangente"],
  ["asin(x)", "Fonction arcsinus"],
  ["acos(x)", "Fonction arccos"],
  ["atan(x)", "Fonction arctan"],
  ["abs(x)", "Fonction valeur absolue"],
  ["sqrt(x)", "Fonction racine carrée"]
] as const;
</script>

<style scoped>
.small-input:deep(input) {
  font-size: 14px;
  width: 100%;
}
</style>
