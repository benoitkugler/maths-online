<template>
  <small class="text-grey mt-1">
    {{ props.description }}
  </small>
  <v-row>
    <v-col md="10" align-self="center">
      <v-table style="overflow-x: auto; max-width: 70vh">
        <tr>
          <th></th>
          <td
            v-for="(_, index) in props.modelValue.Xs"
            style="text-align: center; width: 40px"
            :key="index"
          >
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
        </tr>
        <tr>
          <th>x</th>
          <td v-for="(x, index) in props.modelValue.Xs" :key="index">
            <expression-field
              :model-value="x"
              @update:model-value="(s) => updateXs(index, s)"
              center
              width="50px"
            >
            </expression-field>
          </td>
        </tr>
        <tr>
          <td class="px-2" style="width: 40px">
            <v-text-field
              variant="outlined"
              density="compact"
              :model-value="props.modelValue.Label"
              @update:model-value="updateLabel"
              label="Légende"
              hide-details
              class="label-input"
              :color="latexColor"
            ></v-text-field>
          </td>
          <td v-for="(fx, index) in props.modelValue.Fxs" :key="index">
            <expression-field
              :model-value="fx"
              @update:model-value="(s) => updateFXs(index, s)"
              center
              width="50px"
            >
            </expression-field>
          </td>
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
import { TextKind, type VariationTableBlock } from "@/controller/api_gen";
import { colorByKind } from "@/controller/editor";
import ExpressionField from "../utils/ExpressionField.vue";

interface Props {
  modelValue: VariationTableBlock;
  description: string;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: VariationTableBlock): void;
}>();

const latexColor = colorByKind[TextKind.StaticMath];

function addColumn() {
  props.modelValue.Xs?.push("5");
  props.modelValue.Fxs?.push("5");
  emit("update:modelValue", props.modelValue);
}

function removeColumn(index: number) {
  props.modelValue.Xs?.splice(index, 1);
  props.modelValue.Fxs?.splice(index, 1);
  emit("update:modelValue", props.modelValue);
}

function updateXs(index: number, s: string) {
  props.modelValue.Xs![index] = s;
  emit("update:modelValue", props.modelValue);
}
function updateFXs(index: number, s: string) {
  props.modelValue.Fxs![index] = s;
  emit("update:modelValue", props.modelValue);
}
function updateLabel(label: string) {
  props.modelValue.Label = label;
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
