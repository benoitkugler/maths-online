<template>
  <small class="text-grey mt-1">
    Chaque case est une expression. Les cases d'entête peuvent être du texte,
    une formule LaTeX ($$) ou une expression (&2x+1&).
  </small>
  <v-row>
    <v-col cols="10" align-self="center">
      <div style="overflow-x: auto; max-width: 70vh">
        <table style="table-layout: fixed">
          <tr>
            <td></td>
            <th v-if="props.modelValue.VerticalHeaders != null"></th>
            <td
              v-for="index in rowLength"
              style="text-align: center"
              :key="index"
            >
              <v-btn
                icon
                size="x-small"
                flat
                @click="removeColumn(index - 1)"
                title="Supprimer la colonne"
              >
                <v-icon icon="mdi-close" color="red"></v-icon>
              </v-btn>
            </td>
          </tr>
          <tr v-if="props.modelValue.HorizontalHeaders != null">
            <td></td>
            <th v-if="props.modelValue.VerticalHeaders != null"></th>
            <th
              v-for="(_, index) in props.modelValue.HorizontalHeaders"
              style="text-align: center; width: 40px"
              :key="index"
            >
              <text-part-field
                :model-value="props.modelValue.HorizontalHeaders[index]!"
                @update:model-value="
                v => {
                    props.modelValue.HorizontalHeaders![index] = v;
                emitUpdate();
                }
              "
              >
              </text-part-field>
            </th>
          </tr>
          <tr
            v-for="(row, index) in props.modelValue.Answer || []"
            :key="index"
          >
            <td style="text-align: center; width: 20px">
              <v-btn
                icon
                size="x-small"
                flat
                @click="removeRow(index)"
                title="Supprimer la ligne"
              >
                <v-icon icon="mdi-close" color="red"></v-icon>
              </v-btn>
            </td>
            <th
              v-if="props.modelValue.VerticalHeaders != null"
              style="min-width: 80px"
            >
              <text-part-field
                :model-value="props.modelValue.VerticalHeaders[index]!"
                @update:model-value="
                v => {props.modelValue.VerticalHeaders![index] = v;
emitUpdate();
                }
              "
              >
              </text-part-field>
            </th>
            <td
              v-for="(x, j) in row"
              :key="j"
              style="text-align: center; min-width: 80px"
            >
              <v-text-field
                variant="underlined"
                density="compact"
                hide-details
                :color="expressionColor"
                :model-value="x"
                @update:model-value="v => {
                    row![j] = v;
                    emitUpdate();
                }"
              >
              </v-text-field>
            </td>
          </tr>
        </table>
      </div>
    </v-col>
    <v-col cols="2" align-self="center">
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
  <v-row>
    <v-col align-self="center">
      <v-btn
        @click="addRow"
        title="Ajouter une ligne"
        size="small"
        class="mr-2 my-2"
      >
        <v-icon icon="mdi-plus" color="green"></v-icon>
        Nouvelle ligne
      </v-btn>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { TableFieldBlock, Variable } from "@/controller/api_gen";
import { TextKind } from "@/controller/api_gen";
import { ExpressionColor } from "@/controller/editor";
import { computed } from "vue";
import TextPartField from "./TextPartField.vue";

interface Props {
  modelValue: TableFieldBlock;
  availableParameters: Variable[];
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: TableFieldBlock): void;
}>();

const expressionColor = ExpressionColor;

function emitUpdate() {
  emit("update:modelValue", props.modelValue);
}

function addColumn() {
  if (props.modelValue.HorizontalHeaders != null) {
    props.modelValue.HorizontalHeaders.push({
      Kind: TextKind.Text,
      Content: "",
    });
  }
  props.modelValue.Answer?.forEach((row) => row!.push("1"));
  emitUpdate();
}

const rowLength = computed(() => {
  let len = props.modelValue.HorizontalHeaders?.length;
  if (!len) {
    if (props.modelValue.Answer?.length) {
      len = props.modelValue.Answer[0]?.length;
    } else {
      len = 0;
    }
  }

  return len;
});

function addRow() {
  if (props.modelValue.VerticalHeaders != null) {
    props.modelValue.VerticalHeaders.push({
      Kind: TextKind.Text,
      Content: "",
    });
  }

  props.modelValue.Answer?.push(
    Array.from(new Array(rowLength.value), () => "1")
  );
  emitUpdate();
}

function removeColumn(index: number) {
  if (props.modelValue.HorizontalHeaders != null) {
    props.modelValue.HorizontalHeaders.splice(index, 1);
  }
  props.modelValue.Answer?.forEach((row) => row!.splice(index, 1));

  emitUpdate();
}

function removeRow(index: number) {
  if (props.modelValue.VerticalHeaders != null) {
    props.modelValue.VerticalHeaders.splice(index, 1);
  }
  props.modelValue.Answer?.splice(index, 1);

  emitUpdate();
}
</script>

<style>
.centered-input:deep(input) {
  text-align: center;
}
</style>
