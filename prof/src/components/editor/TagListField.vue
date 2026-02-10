<template>
  <v-dialog v-model="isEditing" :retain-focus="false" max-width="800">
    <v-card title="Modifier les étiquettes">
      <v-card-text>
        <tag-list-edit v-model="tmpList" :all-tags="allTags"></tag-list-edit>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn :disabled="!saveEnabled" @click="endEdit">Enregistrer</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-sheet
    variant="outlined"
    rounded
    border="secondary"
    :style="{
      'border-width': '2px',
      cursor: props.readonly ? '' : 'pointer',
      'text-align': 'center',
    }"
    :class="props.yPadding ? 'py-1' : ''"
  >
    <v-btn
      v-if="props.modelValue.length == 0"
      @click.stop="startEdit"
      flat
      size="x-small"
      class="py-0 px-2"
      :disabled="props.readonly"
      block
    >
      {{ props.readonly ? "Aucune étiquette" : "Ajouter une étiquette..." }}
    </v-btn>

    <v-row no-gutters v-else justify="center" @click.stop="startEdit">
      <v-col v-for="(tag, index) in sorted" :key="index" cols="auto">
        <tag-chip :tag="tag" :pointer="!props.readonly"> </tag-chip>
      </v-col>
    </v-row>
  </v-sheet>
</template>

<script setup lang="ts">
import { Section, type TagsDB, type TagSection } from "@/controller/api_gen";
import { areTagsEquals, tagString } from "@/controller/editor";
import TagListEdit from "./TagListEdit.vue";
import TagChip from "./utils/TagChip.vue";
import { copy } from "@/controller/utils";
import { ref } from "vue";
import { computed } from "vue";

interface Props {
  modelValue: TagSection[];
  allTags: TagsDB;
  readonly: boolean;
  label?: string;
  yPadding?: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: TagSection[]): void;
}>();

defineExpose({ startEdit });

const isEditing = ref(false);
const tmpList = ref<TagSection[]>([]);

function startEdit() {
  if (props.readonly) {
    return;
  }
  isEditing.value = true;
  tmpList.value = props.modelValue.map((v) => ({
    Tag: tagString(v.Tag),
    Section: v.Section,
  }));
}

const saveEnabled = computed(() => {
  return !areTagsEquals(tmpList.value, props.modelValue);
});

const sorted = computed(() => {
  const out = copy(props.modelValue);
  out.sort((a, b) => sectionOrder(a.Section) - sectionOrder(b.Section));
  return out;
});

function endEdit() {
  isEditing.value = false;
  emit("update:model-value", tmpList.value);
}

function sectionOrder(s: Section) {
  switch (s) {
    case Section.Level:
      return 1;
    case Section.Chapter:
      return 3;
    case Section.TrivMath:
      return 4;
    case Section.SubLevel:
      return 2;
    case Section.Matiere:
      return 0;
    default:
      return 5;
  }
}
</script>

<style scoped>
.centered-input:deep(input) {
  text-align: center;
}

:deep(.v-field__append-inner) {
  padding-top: 4px;
}
</style>
