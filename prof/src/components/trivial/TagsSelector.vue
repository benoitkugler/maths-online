<template>
  <v-row class="mx-2">
    <v-col cols="10">
      <v-row v-for="(row, index) in modelValue">
        <v-col cols="10">
          <tag-list-field
            :model-value="row || []"
            @update:model-value="r => (modelValue[index] = r)"
            :all-tags="allTags"
            :ref="(el:any) => (rows[index] = el as TF)"
          ></tag-list-field>
        </v-col>
        <v-col cols="2" align-self="center" style="text-align: center">
          <v-btn
            class="ml-2"
            size="x-small"
            icon
            @click="deleteRow(index)"
            title="Supprimer ce critère"
          >
            <v-icon icon="mdi-delete" color="red"></v-icon>
          </v-btn>
        </v-col>
      </v-row>
    </v-col>
    <v-col cols="2" align-self="center">
      <v-btn
        icon
        size="x-small"
        title="Ajouter une intersection de catégories"
        @click="addIntersection"
      >
        <v-icon icon="mdi-plus" color="green"></v-icon>
      </v-btn>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { nextTick } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import TagListField from "../editor/TagListField.vue";

interface Props {
  modelValue: (string[] | null)[]; // union of intersection
  allTags: string[];
  label?: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: (string[] | null)[]): void;
}>();

type TF = InstanceType<typeof TagListField>;

let rows = $ref<TF[]>([]);

function addIntersection() {
  props.modelValue.push([]);
  emit("update:model-value", props.modelValue);
  nextTick(() => {
    rows[rows.length - 1].startEdit();
  });
}

function deleteRow(index: number) {
  props.modelValue.splice(index, 1);
  emit("update:model-value", props.modelValue);
}
</script>

<style scoped></style>
