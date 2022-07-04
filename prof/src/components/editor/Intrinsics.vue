<template>
  <v-dialog v-model="showHelp" min-width="1200px" width="max-content">
    <v-card
      title="Fonctions spéciales"
      subtitle="Description"
      style="width: 800px"
    >
      <v-card-text>
        Les fonctions spéciales permettent de définir des paramètres aléatoires
        complexes de manière rapide et simple. La syntaxe d'une définition suit
        le format : <br />
        <p class="my-2">
          <i>a,b,c,d,... </i> = <b>fonction</b>(<i>argument1, argument2, ...</i
          >)
        </p>
        Les fonctions utilisables sont les suivantes :
        <v-list color="info" rounded>
          <v-list-item>
            <v-row>
              <v-col cols="6">
                <v-list-item-title
                  >a, b, c = pythagorians(bound)</v-list-item-title
                >
              </v-col>
              <v-col align-self="center">
                <v-list-item-subtitle
                  >Génère trois entiers <i>a</i>,<i>b</i>,<i>c</i> vérifiant a^2
                  + b^2 = c^2. <i>bound</i> est un argument optionnel qui
                  controle le maximum de <i>a</i> par <i>2 bound^2</i>
                </v-list-item-subtitle>
              </v-col>
            </v-row>
          </v-list-item>
          <v-list-item>
            <v-row>
              <v-col cols="6">
                <v-list-item-title>H = projection(A, B, C)</v-list-item-title>
              </v-col>
              <v-col align-self="center">
                <v-list-item-subtitle>
                  Calcule le projeté orthogonal du point <i>A</i> sur
                  <i>(BC)</i>.
                </v-list-item-subtitle>
              </v-col>
            </v-row>
          </v-list-item>
        </v-list>
      </v-card-text>
    </v-card>
  </v-dialog>

  <v-card>
    <v-row
      :style="{
        'background-color': props.isValidated ? 'lightgreen' : 'lightgray',
      }"
      class="rounded"
      no-gutters
    >
      <v-col md="5" align-self="center">
        <v-card-subtitle class="py-2">Fonctions spéciales</v-card-subtitle>
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
          title="Ajouter une fonction spéciale"
          size="x-small"
          class="mr-2 my-2"
        >
          <v-icon icon="mdi-plus" color="green"></v-icon>
        </v-btn>
        <v-btn
          icon
          @click="showHelp = true"
          title="Aide"
          size="x-small"
          class="mr-2 my-2"
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
            @drop="(origin) => emit('swap', origin, 0)"
          ></drop-zone>
          <div v-for="(param, index) in props.parameters" :key="index">
            <v-list-item class="pr-0 pl-1">
              <v-row no-gutters>
                <v-col cols="1" align-self="center">
                  <v-icon
                    size="large"
                    class="pr-2"
                    style="cursor: grab"
                    @dragstart="(e) => onItemDragStart(e, index)"
                    draggable="true"
                    color="green-lighten-3"
                    icon="mdi-drag-vertical"
                  ></v-icon>
                </v-col>
                <v-col cols="9" class="pl-1">
                  <v-text-field
                    class="small-input"
                    hide-details
                    variant="underlined"
                    density="compact"
                    :model-value="param"
                    @update:model-value="(v) => autocomplete(index, v)"
                    @blur="emit('done')"
                  ></v-text-field>
                </v-col>
                <v-col cols="2" align-self="center">
                  <v-btn
                    icon
                    size="small"
                    flat
                    @click="emit('delete', index)"
                    title="Supprimer cette fonction spéciale"
                  >
                    <v-icon icon="mdi-delete" color="red"></v-icon>
                  </v-btn>
                </v-col>
              </v-row>
            </v-list-item>
            <drop-zone
              v-if="showDropZone"
              @drop="(origin) => emit('swap', origin, index + 1)"
            ></drop-zone>
          </div>
        </v-list>
      </v-col>
    </v-row>
  </v-card>
</template>

<script setup lang="ts">
import { onDragListItemStart } from "@/controller/editor";
import { ref } from "@vue/reactivity";
import { $ref } from "vue/macros";
import DropZone from "./DropZone.vue";

interface Props {
  parameters: string[];
  isLoading: boolean;
  isValidated: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "add"): void;
  (e: "update", index: number, param: string): void;
  (e: "delete", index: number): void;
  (e: "swap", origin: number, target: number): void;
  (e: "done"): void;
}>();

const showHelp = ref(false);

let showDropZone = $ref(false);
function onItemDragStart(payload: DragEvent, index: number) {
  onDragListItemStart(payload, index);
  setTimeout(() => (showDropZone = true), 100); // workaround bug
}

// to keep sync with the server
const intrincics = ["pythagorians", "projection"];

function autocomplete(index: number, text: string) {
  for (const it of intrincics) {
    if (text.endsWith(it.substr(0, 4))) {
      text += it.substring(4) + "()";
      break;
    }
  }

  emit("update", index, text);
}
</script>

<style scoped>
.small-input:deep(input) {
  font-size: 14px;
  width: 100%;
}
</style>
