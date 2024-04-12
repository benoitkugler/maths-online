<template>
  <v-row class="mx-2">
    <v-col cols="10">
      <v-row v-for="(row, index) in modelValue" :key="index">
        <v-col cols="10">
          <tag-list-field
            :model-value="row || []"
            @update:model-value="(r) => updateRow(index, r)"
            :all-tags="allTags"
            :ref="(el:any) => (rows[index] = el as TF)"
            :readonly="false"
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
import {
  Section,
  type Tags,
  type TagsDB,
  type TagSection,
} from "@/controller/api_gen";
import { ref, nextTick } from "vue";
import TagListField from "../editor/TagListField.vue";
import type { PrefillTrivialCategorie } from "@/controller/utils";

interface Props {
  modelValue: Tags[]; // union of intersection
  lastMatiereLevelChapter: PrefillTrivialCategorie;
  allTags: TagsDB;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: Tags[]): void;
}>();

type TF = InstanceType<typeof TagListField>;

const rows = ref<TF[]>([]);

function addIntersection() {
  // defaut to last level and section, if any
  const newTags: Tags = [];
  if (props.lastMatiereLevelChapter.level.length) {
    newTags.push({
      Section: Section.Level,
      Tag: props.lastMatiereLevelChapter.level,
    });
  }
  if (props.lastMatiereLevelChapter.matiere.length) {
    newTags.push({
      Section: Section.Matiere,
      Tag: props.lastMatiereLevelChapter.matiere,
    });
  }
  if (props.lastMatiereLevelChapter.chapter.length) {
    newTags.push({
      Section: Section.Chapter,
      Tag: props.lastMatiereLevelChapter.chapter,
    });
  }
  newTags.push(...(props.lastMatiereLevelChapter.sublevels || []));

  props.modelValue.push(newTags);
  emit("update:model-value", props.modelValue);
  nextTick(() => {
    rows.value[rows.value.length - 1].startEdit();
  });
}

function updateRow(index: number, r: TagSection[]) {
  props.modelValue[index] = r;
  emit("update:model-value", props.modelValue);
}

function deleteRow(index: number) {
  props.modelValue.splice(index, 1);
  emit("update:model-value", props.modelValue);
}
</script>

<style scoped></style>
