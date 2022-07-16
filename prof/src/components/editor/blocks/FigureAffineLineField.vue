<template>
  <v-dialog v-model="showDocumentation" width="600">
    <v-card title="Syntaxe de la réponse">
      <v-card-text>
        La réponse est une droite définie par son coefficient directeur
        <i>A</i> et son ordonnée à l'origine <i>B</i>. <br />
        Une droite verticale est obtenue en utilisant <i>A=Inf</i> et en fixant
        <i>B</i> à l'abscisse souhaitée. Remarquez que <i>Inf</i> peut s'obtenir
        par exemple avec l'expression <i>A=1/c</i>, où <i>c=0</i>.
      </v-card-text>
    </v-card>
  </v-dialog>
  <v-card class="my-1">
    <v-row no-gutters class="bg-secondary py-2">
      <v-col align-self="center">
        <v-card-subtitle>Définition de la réponse</v-card-subtitle>
      </v-col>
      <v-col align-self="center" cols="auto">
        <v-btn
          class="mr-2"
          icon
          title="Documentation de la syntaxe"
          size="x-small"
        >
          <v-icon small color="info" @click="showDocumentation = true"
            >mdi-help</v-icon
          >
        </v-btn>
      </v-col>
    </v-row>

    <v-card-text>
      <v-row class="fix-input-width">
        <v-col align-self="center" cols="4">
          <v-text-field
            density="compact"
            variant="outlined"
            label="Légende"
            v-model="props.modelValue.Label"
            hide-details
          ></v-text-field>
        </v-col>
        <v-col align-self="center">
          <v-text-field
            density="compact"
            variant="outlined"
            label="A"
            hint="Coefficient directeur, Inf pour une droite verticale."
            v-model="props.modelValue.A"
            :color="expressionColor"
            persistent-hint
            class="no-hint-padding"
          ></v-text-field>
        </v-col>
        <v-col align-self="center">
          <v-text-field
            density="compact"
            variant="outlined"
            label="B"
            v-model="props.modelValue.B"
            hint="Ordonnée à l'origine ou abscisse pour une droite verticale."
            persistent-hint
            :color="expressionColor"
            class="no-hint-padding"
          ></v-text-field>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
  <figure-block-vue
    v-model="props.modelValue.Figure"
    :available-parameters="props.availableParameters"
  ></figure-block-vue>
</template>

<script setup lang="ts">
import type {
FigureAffineLineFieldBlock,
Variable
} from "@/controller/api_gen";
import { TextKind } from "@/controller/api_gen";
import { colorByKind } from "@/controller/editor";
import { $ref } from "vue/macros";
import FigureBlockVue from "./FigureBlock.vue";

interface Props {
  modelValue: FigureAffineLineFieldBlock;
  availableParameters: Variable[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: FigureAffineLineFieldBlock): void;
}>();

const expressionColor = colorByKind[TextKind.Expression];

let showDocumentation = $ref(false);
</script>

<style scoped>
.no-hint-padding:deep(.v-input__details) {
  padding-inline: 0px;
}

.fix-input-width:deep(input) {
  width: 100%;
}
</style>
