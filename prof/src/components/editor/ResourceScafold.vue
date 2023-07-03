<template>
  <v-dialog
    :model-value="variantToDelete != null"
    @update:model-value="variantToDelete = null"
    max-width="700"
  >
    <v-card title="Confirmer">
      <v-card-text
        >Etes-vous certain de vouloir supprimer la variante
        <i>{{ variantToDelete?.Id }} - {{ variantToDelete?.Subtitle }}</i> ?
        <br />
        Cette opération est irréversible.

        <div v-if="props.resource.Variants.length == 1" class="mt-2">
          Le groupe associé sera aussi supprimé.
        </div>
      </v-card-text>
      <v-card-actions>
        <v-btn @click="variantToDelete = null">Retour</v-btn>
        <v-spacer></v-spacer>
        <v-btn
          color="red"
          @click="
            emit('deleteVariant', variantToDelete!);
            variantToDelete = null;
          "
          variant="outlined"
        >
          Supprimer
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog
    :model-value="variantToEdit != null"
    @update:model-value="variantToEdit = null"
    max-width="800"
  >
    <ResourceVariantEdit
      v-if="variantToEdit != null"
      :variant="variantToEdit"
      :readonly="props.readonly"
      @save="
        (v) => {
          variantToEdit = null;
          emit('updateVariant', v);
        }
      "
      @close="variantToEdit = null"
    ></ResourceVariantEdit>
  </v-dialog>

  <v-dialog
    :model-value="tagsToEdit != null"
    @update:model-value="tagsToEdit = null"
    max-width="800"
  >
    <TagListEdit
      v-if="tagsToEdit != null"
      v-model="tagsToEdit"
      :all-tags="props.allTags"
      :save-enabled="areEditedTagsDistinct"
      @save="
        emit('updateTags', tagsToEdit!);
        tagsToEdit = null;
      "
    ></TagListEdit>
  </v-dialog>

  <v-card class="mt-2">
    <v-row no-gutters class="pl-2 pt-1">
      <v-col cols="auto" align-self="center">
        <v-btn
          size="small"
          icon
          title="Retour à la liste"
          @click="emit('back')"
          class="mr-1"
        >
          <v-icon icon="mdi-arrow-left"></v-icon>
        </v-btn>

        <v-tooltip content-class="bg-grey-lighten-4">
          <template v-slot:activator="{ props: innerProps }">
            <v-badge
              v-bind="innerProps"
              :content="props.resource.Tags?.length || 0"
              :color="ChapterColor"
              class="mr-3"
            >
              <v-btn
                icon
                size="x-small"
                :color="LevelColor"
                @click="tagsToEdit = copy(props.resource.Tags || [])"
                :disabled="props.readonly"
              >
                <v-icon size="small">mdi-tag</v-icon>
              </v-btn>
            </v-badge>
          </template>
          <v-row no-gutters justify="center" v-if="props.resource.Tags?.length">
            <v-col
              v-for="(tag, index) in props.resource.Tags"
              :key="index"
              cols="auto"
            >
              <tag-chip :tag="tag"> </tag-chip>
            </v-col>
          </v-row>
          <span v-else>Aucune étiquette</span>
        </v-tooltip>
      </v-col>

      <v-col align-self="center" :cols="titleToEdit == null ? 4 : 8">
        <v-hover
          v-slot="{ isHovering, props: innerProps }"
          v-if="titleToEdit == null"
        >
          <v-card-subtitle
            v-bind="innerProps"
            :style="{
              'white-space': 'unset',
              cursor: props.readonly ? 'unset' : 'text',
              border:
                isHovering && !props.readonly
                  ? '1px solid grey'
                  : '1px solid transparent',
            }"
            :class="{ 'text-subtitle-1': true, rounded: true, 'py-3': true }"
            @click="props.readonly ? {} : (titleToEdit = props.resource.Title)"
            >{{ props.resource.Title || "Aucun titre" }}
          </v-card-subtitle>
        </v-hover>
        <v-text-field
          v-else
          color="grey"
          class="mt-1"
          label="Titre de la ressource"
          variant="outlined"
          density="compact"
          v-model="titleToEdit"
          autofocus
          @focus="($event.target as HTMLInputElement)?.select()"
          hide-details
          @blur="onDoneEditTitle"
        ></v-text-field>
      </v-col>

      <v-spacer></v-spacer>
      <v-col cols="6" align-self="center" v-if="titleToEdit == null">
        <v-tabs
          style="max-width: 90vh"
          density="compact"
          show-arrows
          color="grey"
          :model-value="props.modelValue"
          @update:model-value="(i) => emit('update:model-value', i as number)"
          align-tabs="end"
        >
          <v-tooltip
            v-for="(variant, index) in props.resource.Variants"
            :key="index"
            :text="tabTooltip(variant)"
          >
            <template v-slot:activator="{ props: innerProps }">
              <v-tab
                v-bind="innerProps"
                class="text-subtitle-2 font-weight-light"
              >
                {{ tabTitle(index) }}
                <TagChip
                  v-if="variant.Difficulty && index == props.modelValue"
                  :tag="{ Tag: variant.Difficulty }"
                ></TagChip>
                <v-menu
                  offset-y
                  close-on-content-click
                  v-if="index == props.modelValue"
                >
                  <template v-slot:activator="{ isActive, props: innerProps2 }">
                    <v-btn
                      v-on="{ isActive }"
                      v-bind="innerProps2"
                      icon
                      size="x-small"
                      flat
                      class="pr-0 mr-0"
                    >
                      <v-icon>mdi-dots-vertical</v-icon>
                    </v-btn>
                  </template>
                  <v-list>
                    <v-list-item @click="variantToEdit = copy(variant)" link>
                      <template v-slot:prepend>
                        <v-icon
                          icon="mdi-pencil"
                          color="info"
                          size="small"
                        ></v-icon>
                      </template>
                      <v-list-item-title> Détails </v-list-item-title>
                    </v-list-item>
                    <v-list-item
                      @click="
                        props.readonly ? {} : emit('duplicateVariant', variant)
                      "
                      :link="!props.readonly"
                    >
                      <template v-slot:prepend>
                        <v-icon
                          icon="mdi-content-copy"
                          color="info"
                          size="small"
                        ></v-icon>
                      </template>
                      Dupliquer
                    </v-list-item>
                    <v-list-item
                      @click="props.readonly ? {} : (variantToDelete = variant)"
                      :link="!props.readonly"
                    >
                      <template v-slot:prepend>
                        <v-icon
                          icon="mdi-delete"
                          color="red"
                          size="small"
                        ></v-icon>
                      </template>

                      Supprimer</v-list-item
                    >
                  </v-list>
                </v-menu>
              </v-tab>
            </template>
          </v-tooltip>
        </v-tabs>
      </v-col>
    </v-row>

    <slot></slot>
  </v-card>
</template>

<script setup lang="ts">
import {
  areTagsEquals,
  ChapterColor,
  LevelColor,
  type ResourceGroup,
  type VariantG,
} from "@/controller/editor";
import { $ref } from "vue/macros";

import TagChip from "./utils/TagChip.vue";
import ResourceVariantEdit from "./ResourceVariantEdit.vue";
import { copy } from "@/controller/utils";
import type { Tags, TagsDB, TagSection } from "@/controller/api_gen";
import TagListEdit from "./TagListEdit.vue";
import { computed } from "vue";

interface Props {
  resource: ResourceGroup;
  readonly: boolean;
  modelValue: number; // index into variants
  allTags: TagsDB;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "back"): void;
  (e: "update:model-value", v: number): void;
  (e: "updateTitle", title: string): void;
  (e: "updateTags", tags: Tags): void;
  (e: "updateVariant", v: VariantG): void;
  (e: "duplicateVariant", v: VariantG): void;
  (e: "deleteVariant", v: VariantG): void;
}>();

function tabTitle(index: number) {
  if (index == props.modelValue)
    return props.resource.Variants[index].Subtitle || `Variante ${index + 1}`;
  return `Var. ${index + 1}`;
}

function tabTooltip(variant: VariantG) {
  return `(${variant.Id}) -  ${variant.Subtitle || "Sans titre"}`;
}

let variantToDelete = $ref<VariantG | null>(null);

let variantToEdit = $ref<VariantG | null>(null);

let titleToEdit = $ref<string | null>(null);
function onDoneEditTitle() {
  const newTitle = titleToEdit || "";
  // avoid useless query
  if (newTitle != props.resource.Title) {
    emit("updateTitle", newTitle);
  }
  titleToEdit = null;
}

let tagsToEdit = $ref<TagSection[] | null>(null);
const areEditedTagsDistinct = computed(
  () => !areTagsEquals(props.resource.Tags, tagsToEdit)
);
</script>

<style scoped></style>
