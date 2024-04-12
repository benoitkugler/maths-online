<template>
  <small class="text-grey mt-1">
    Les cases de la ligne <b>x</b> sont des expressions. Utiliser inf (ou -inf)
    pour l'infini. <br />
    Le nom d'une fonction peut contenir une expression, marquée par &&.
  </small>
  <v-table style="overflow-x: auto; max-width: 75vh">
    <!-- header -->
    <tr>
      <th style="min-width: 120px" colspan="2">
        <v-btn
          icon
          @click="addColumn"
          title="Ajouter une colonne"
          size="x-small"
          class="mr-2 my-2"
        >
          <v-icon icon="mdi-plus" color="green"></v-icon>
        </v-btn>
      </th>
      <template v-for="(x, index) in props.modelValue.Xs" :key="index">
        <td v-if="index" style="min-width: 60px"></td>
        <td style="text-align: center; min-width: 60px">
          <v-btn
            icon
            size="x-small"
            flat
            @click="removeColumn(index)"
            title="Supprimer la colonne"
          >
            <v-icon icon="mdi-close" color="red"></v-icon>
          </v-btn>
        </td>
      </template>
    </tr>
    <!-- X values -->
    <tr>
      <th></th>
      <th>x</th>
      <template v-for="(x, index) in props.modelValue.Xs" :key="index">
        <td v-if="index"></td>
        <td style="text-align: center">
          <expression-field
            :model-value="x"
            @update:model-value="s => props.modelValue.Xs![index] = s"
            center
            width="50px"
          >
          </expression-field>
        </td>
      </template>
    </tr>
    <!-- Functions -->
    <tr v-for="(fn, index) in props.modelValue.Functions" :key="index">
      <td style="text-align: center; width: 20px">
        <v-btn
          icon
          size="x-small"
          flat
          @click="removeRow(index)"
          title="Supprimer la fonction"
        >
          <v-icon icon="mdi-close" color="red"></v-icon>
        </v-btn>
      </td>
      <td class="px-2">
        <interpolated-text
          v-model="fn.Label"
          force-latex
          center
        ></interpolated-text>
      </td>
      <template v-for="(fx, j) in fn.FxSymbols" :key="j">
        <td v-if="j" style="text-align: center">
          <v-btn
            size="small"
            rounded
            @click="fn.Signs![j - 1] = !fn.Signs![j - 1]"
          >
            <b>
              {{ fn.Signs![j - 1] ? "+" : "-" }}
            </b>
          </v-btn>
        </td>
        <td style="text-align: center">
          <sign-symbol-field
            :model-value="fn.FxSymbols![j]"
            @update:model-value="v => fn.FxSymbols![j]=v"
          >
          </sign-symbol-field>
        </td>
      </template>
    </tr>
  </v-table>
  <v-row>
    <v-col align-self="center">
      <v-btn
        @click="addRow"
        title="Ajouter une fonction"
        size="small"
        class="mr-2 my-2"
      >
        <v-icon icon="mdi-plus" color="green"></v-icon>
        Ajouter une fonction
      </v-btn>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { SignTableBlock } from "@/controller/api_gen";
import { SignSymbol, TextKind } from "@/controller/api_gen";
import ExpressionField from "../utils/ExpressionField.vue";
import SignSymbolField from "./SignSymbolField.vue";
import { colorByKind } from "@/controller/editor";
import InterpolatedText from "../utils/InterpolatedText.vue";

interface Props {
  modelValue: SignTableBlock;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: SignTableBlock): void;
}>();

const latexColor = colorByKind[TextKind.StaticMath];

function addColumn() {
  props.modelValue.Xs?.push("5");
  props.modelValue.Functions?.forEach(fn => {
    if (fn.FxSymbols?.length) {
      // do not add sign if there is only on x
      fn.Signs?.push(true);
    }
    fn.FxSymbols?.push(SignSymbol.Nothing);
  });
  emit("update:modelValue", props.modelValue);
}

function removeColumn(index: number) {
  props.modelValue.Xs?.splice(index, 1);
  props.modelValue.Functions?.forEach(fn => {
    fn.FxSymbols?.splice(index, 1);
    fn.Signs?.splice(index == 0 ? 0 : index - 1, 1);
  });
  emit("update:modelValue", props.modelValue);
}

function addRow() {
  const L = props.modelValue.Xs?.length || 0;
  props.modelValue.Functions = (props.modelValue.Functions || []).concat({
    Label: "f(x)",
    FxSymbols: Array.from(new Array(L), () => SignSymbol.Nothing),
    Signs: L == 0 ? [] : Array.from(new Array(L - 1), () => true)
  });
  emit("update:modelValue", props.modelValue);
}

function removeRow(index: number) {
  props.modelValue.Functions?.splice(index, 1);
}
</script>

<style scoped>
.label-input:deep(input) {
  width: 50px;
  text-align: center;
}
.label-input:deep(.v-field__input) {
  padding-inline: 4px;
}
</style>
