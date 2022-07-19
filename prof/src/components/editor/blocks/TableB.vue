<template>
  <small class="text-grey mt-1">
    Chaque case peut être du texte, une formule LaTeX ($$) ou une expression
    (&2x+1&).
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
              :key="index"
              style="text-align: center; width: 40px"
            >
              <text-part-field
                :model-value="props.modelValue.HorizontalHeaders[index]!"
                @update:model-value="
                v => props.modelValue.HorizontalHeaders![index] = v
              "
              >
              </text-part-field>
            </th>
          </tr>
          <tr
            v-for="(row, index) in props.modelValue.Values || []"
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
                v => props.modelValue.VerticalHeaders![index] = v
              "
              >
              </text-part-field>
            </th>
            <td
              v-for="(x, j) in row"
              :key="j"
              style="text-align: center; min-width: 80px"
            >
              <text-part-field
                :model-value="x"
                @update:model-value="v => (row![j] = v)"
              >
              </text-part-field>
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
    <v-col cols="4" align-self="center">
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
    <v-col cols="4">
      <v-switch
        label="Entête horizontal"
        hide-details
        :model-value="props.modelValue.HorizontalHeaders != null"
        @update:model-value="toogleHorizontal"
        color="secondary"
      ></v-switch>
    </v-col>
    <v-col cols="4">
      <v-switch
        label="Entête vertical"
        hide-details
        :model-value="props.modelValue.VerticalHeaders != null"
        @update:model-value="toogleVertical"
        color="secondary"
      ></v-switch>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { TableBlock } from "@/controller/api_gen";
import { TextKind } from "@/controller/api_gen";
import { computed } from "@vue/runtime-core";
import TextPartField from "./TextPartField.vue";

interface Props {
  modelValue: TableBlock;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: TableBlock): void;
}>();

function addColumn() {
  if (props.modelValue.HorizontalHeaders != null) {
    props.modelValue.HorizontalHeaders.push({
      Kind: TextKind.Text,
      Content: "",
    });
  }
  props.modelValue.Values?.forEach((row) =>
    row!.push({
      Kind: TextKind.Text,
      Content: "",
    })
  );

  emit("update:modelValue", props.modelValue);
}

const rowLength = computed(() => {
  let len = props.modelValue.HorizontalHeaders?.length;
  if (!len) {
    if (props.modelValue.Values?.length) {
      len = props.modelValue.Values[0]?.length;
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

  props.modelValue.Values?.push(
    Array.from(new Array(rowLength.value), () => ({
      Kind: TextKind.Text,
      Content: "",
    }))
  );

  emit("update:modelValue", props.modelValue);
}

function removeColumn(index: number) {
  console.log("removecol", index);

  if (props.modelValue.HorizontalHeaders != null) {
    props.modelValue.HorizontalHeaders.splice(index, 1);
  }
  props.modelValue.Values?.forEach((row) => row!.splice(index, 1));

  emit("update:modelValue", props.modelValue);
}

function removeRow(index: number) {
  if (props.modelValue.VerticalHeaders != null) {
    props.modelValue.VerticalHeaders.splice(index, 1);
  }
  props.modelValue.Values?.splice(index, 1);

  emit("update:modelValue", props.modelValue);
}

function toogleHorizontal(b: boolean) {
  if (b) {
    props.modelValue.HorizontalHeaders = Array.from(
      new Array(rowLength.value),
      () => ({
        Kind: TextKind.Text,
        Content: "",
      })
    );
  } else {
    props.modelValue.HorizontalHeaders = null;
  }

  emit("update:modelValue", props.modelValue);
}

function toogleVertical(b: boolean) {
  if (b) {
    props.modelValue.VerticalHeaders = Array.from(
      new Array(props.modelValue.Values?.length || 0),
      () => ({
        Kind: TextKind.Text,
        Content: "",
      })
    );
  } else {
    props.modelValue.VerticalHeaders = null;
  }

  emit("update:modelValue", props.modelValue);
}
</script>

<style>
.centered-input:deep(input) {
  text-align: center;
}
</style>
