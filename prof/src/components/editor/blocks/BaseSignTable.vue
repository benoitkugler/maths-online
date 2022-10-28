<template>
  <small class="text-grey mt-1">
    Les cases de la ligne <b>x</b> sont des expressions. Utiliser inf (ou -inf)
    pour l'infini.
  </small>
  <v-row>
    <v-col md="10" align-self="center">
      <v-table style="overflow-x: auto; max-width: 70vh">
        <tr>
          <th style="min-width: 100px"></th>
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
        <tr>
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
        <tr>
          <td class="px-2">
            <v-text-field
              variant="outlined"
              density="compact"
              v-model="props.modelValue.Label"
              label="Légende"
              hide-details
              class="label-input"
            ></v-text-field>
          </td>
          <template
            v-for="(fx, index) in props.modelValue.FxSymbols"
            :key="index"
          >
            <td v-if="index" style="text-align: center">
              <v-btn
                size="small"
                rounded
                @click="
                  props.modelValue.Signs![index - 1] =
                    !props.modelValue.Signs![index - 1]
                "
              >
                {{ props.modelValue.Signs![index - 1] ? "+" : "-" }}
              </v-btn>
            </td>
            <td style="text-align: center">
              <sign-symbol-field
                :model-value="props.modelValue.FxSymbols![index]"
                @update:model-value="v => props.modelValue.FxSymbols![index]=v"
              >
              </sign-symbol-field>
            </td>
          </template>
        </tr>
      </v-table>
    </v-col>
    <v-col md="2" align-self="center">
      <v-btn
        icon
        @click="addColumn"
        title="Ajouter une colonne"
        size="x-small"
        class="mr-2 my-2"
      >
        <v-icon icon="mdi-plus" color="green"></v-icon>
      </v-btn>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { SignTableBlock } from "@/controller/api_gen";
import { SignSymbol } from "@/controller/api_gen";
import ExpressionField from "../utils/ExpressionField.vue";
import SignSymbolField from "./SignSymbolField.vue";

interface Props {
  modelValue: SignTableBlock;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: SignTableBlock): void;
}>();

function addColumn() {
  props.modelValue.Xs?.push("5");
  props.modelValue.FxSymbols?.push(SignSymbol.Zero);
  props.modelValue.Signs?.push(true);
  emit("update:modelValue", props.modelValue);
}

function removeColumn(index: number) {
  props.modelValue.Xs?.splice(index, 1);
  props.modelValue.FxSymbols?.splice(index, 1);
  props.modelValue.Signs?.splice(index == 0 ? 0 : index - 1, 1);
  emit("update:modelValue", props.modelValue);
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
