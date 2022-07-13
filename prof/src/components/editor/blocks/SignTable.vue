<template>
  <small class="text-grey mt-1">
    Utiliser &2a - 4& pour insérer une expression dans la ligne x.
  </small>
  <v-row>
    <v-col md="10" align-self="center">
      <v-table style="overflow-x: auto; max-width: 70vh">
        <tr>
          <th style="width: 60px"></th>
          <template v-for="(x, index) in props.modelValue.Xs">
            <td v-if="index"></td>
            <td style="text-align: center; width: 40px">
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
          <th style="width: 60px">x</th>
          <template v-for="(x, index) in props.modelValue.Xs">
            <td v-if="index"></td>
            <td style="text-align: center; width: 80px">
              <interpolated-text
                :model-value="x"
                @update:model-value="s => props.modelValue.Xs![index] = s"
              >
              </interpolated-text>
            </td>
          </template>
        </tr>
        <tr>
          <td class="px-2" style="width: 60px">
            <v-text-field
              variant="outlined"
              density="compact"
              v-model="props.modelValue.Label"
              label="Légende"
              hide-details
              class="fix-input-width"
            ></v-text-field>
          </td>
          <template v-for="(fx, index) in props.modelValue.FxSymbols">
            <td v-if="index">
              <v-btn
                size="small"
                rounded
                @click="props.modelValue.Signs![index - 1] = !props.modelValue.Signs![index - 1]"
              >
                {{ props.modelValue.Signs![index - 1] ? "+" : "-" }}
              </v-btn>
            </td>
            <td style="text-align: center; width: 80px">
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
import InterpolatedText from "../utils/InterpolatedText.vue";
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
}

function removeColumn(index: number) {
  props.modelValue.Xs?.splice(index, 1);
  props.modelValue.FxSymbols?.splice(index, 1);
  props.modelValue.Signs?.splice(index == 0 ? 0 : index - 1, 1);
}
</script>

<style scoped>
.small-select:deep(.v-field__input) {
  padding-right: 0;
}

:deep(.v-select__selections input) {
  padding: 0px;
}

.fix-input-width:deep(input) {
  width: 45px;
}
</style>
