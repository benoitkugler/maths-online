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
    :model-value="titleAndTagsToEdit != null"
    @update:model-value="titleAndTagsToEdit = null"
    max-width="800"
  >
    <v-card
      title="Modifier le titre et les étiquettes"
      v-if="titleAndTagsToEdit != null"
    >
      <v-card-text class="my-1">
        <v-row>
          <v-col>
            <v-text-field
              density="compact"
              variant="outlined"
              v-model="titleAndTagsToEdit.title"
              label="Titre"
            ></v-text-field>
          </v-col>
        </v-row>
        <TagListEdit
          v-model="titleAndTagsToEdit.tags"
          :all-tags="props.allTags"
        ></TagListEdit>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          :disabled="
            !areEditedTagsDistinct &&
            props.resource.Title == titleAndTagsToEdit.title
          "
          @click="
            emit(
              'updateTitleAndTags',
              titleAndTagsToEdit.title,
              titleAndTagsToEdit.tags
            );
            titleAndTagsToEdit = null;
          "
          >Enregistrer</v-btn
        >
      </v-card-actions>
    </v-card>
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
      </v-col>

      <v-col :cols="5" class="px-2">
        <v-card
          :subtitle="props.resource.Title || 'Aucun titre'"
          @[!props.readonly&&`click`]="startEdit"
        >
          <v-card-text class="pa-2 pt-0">
            <v-row no-gutters justify="start"
              ><v-col>
                <v-row no-gutters v-if="props.resource.Tags?.length">
                  <v-col
                    v-for="(tag, index) in props.resource.Tags"
                    :key="index"
                    cols="auto"
                  >
                    <tag-chip :tag="tag"> </tag-chip>
                  </v-col>
                </v-row>
                <span v-else>Aucune étiquette</span>
              </v-col></v-row
            >
          </v-card-text>
        </v-card>
      </v-col>

      <v-spacer></v-spacer>
      <!-- Variants tabs -->
      <v-col cols="6" align-self="center">
        <v-tabs
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
                    <v-list-item @click="variantToEdit = copy(variant)">
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
  type ResourceGroup,
  type VariantG,
} from "@/controller/editor";

import TagChip from "./utils/TagChip.vue";
import ResourceVariantEdit from "./ResourceVariantEdit.vue";
import { copy } from "@/controller/utils";
import type { Tags, TagsDB, TagSection } from "@/controller/api_gen";
import TagListEdit from "./TagListEdit.vue";
import { computed, ref } from "vue";

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
  (e: "updateTitleAndTags", title: string, tags: Tags): void;
  (e: "updateVariant", v: VariantG): void;
  (e: "duplicateVariant", v: VariantG): void;
  (e: "deleteVariant", v: VariantG): void;
}>();

defineExpose({ showEditVariant, startEdit });

function startEdit() {
  titleAndTagsToEdit.value = {
    title: props.resource.Title,
    tags: copy(props.resource.Tags || []),
  };
}

function showEditVariant(variant: VariantG) {
  variantToEdit.value = copy(variant);
}

function tabTitle(index: number) {
  if (index == props.modelValue)
    return props.resource.Variants[index].Subtitle || `Variante ${index + 1}`;
  return `Var. ${index + 1}`;
}

function tabTooltip(variant: VariantG) {
  return `(${variant.Id}) -  ${variant.Subtitle || "Sans titre"}`;
}

const variantToDelete = ref<VariantG | null>(null);

const variantToEdit = ref<VariantG | null>(null);

const titleAndTagsToEdit = ref<{ tags: TagSection[]; title: string } | null>(
  null
);
const areEditedTagsDistinct = computed(
  () =>
    !areTagsEquals(props.resource.Tags, titleAndTagsToEdit.value?.tags || [])
);
</script>

<style scoped></style>
