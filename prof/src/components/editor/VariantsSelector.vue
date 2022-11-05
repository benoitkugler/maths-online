<template>
  <v-menu offset-y :close-on-content-click="true" v-model="showMenu">
    <template v-slot:activator="{ isActive, props: slotProps }">
      <v-btn
        title="Gérer les variantes"
        v-on="{ isActive }"
        v-bind="slotProps"
        block
      >
        Variantes...
      </v-btn>
    </template>
    <v-card width="600">
      <v-card-text>
        <v-list style="max-height: 70vh">
          <v-list-item
            density="comfortable"
            :active="index == props.modelValue"
            rounded
            class="my-1"
            v-for="(variant, index) in variants"
            link
            @click="
              emit('update:model-value', index);
              showMenu = false;
            "
            :key="variant.Id"
          >
            <v-row>
              <v-col align-self="center" cols="auto">
                <v-btn
                  v-if="!props.readonly"
                  class="ml-1"
                  size="small"
                  icon
                  @click.stop="
                    showMenu = false;
                    emit('delete', variant);
                  "
                  title="Supprimer la variante"
                >
                  <v-icon icon="mdi-delete" color="red" size="small"></v-icon>
                </v-btn>

                <v-btn
                  v-if="!props.readonly"
                  class="mx-1"
                  size="small"
                  icon
                  @click.stop="emit('duplicate', variant)"
                  title="Dupliquer la variante"
                >
                  <v-icon
                    icon="mdi-content-copy"
                    color="info"
                    size="small"
                  ></v-icon>
                </v-btn>
              </v-col>
              <v-col align-self="center" class="my-4" cols="7">
                ({{ variant.Id }})
                <template v-if="variant.Subtitle">{{
                  variant.Subtitle
                }}</template>
                <i v-else>(Sans titre)</i>
              </v-col>
              <v-col align-self="center" style="text-align: right">
                <TagChip
                  v-if="variant.Difficulty"
                  :tag="variant.Difficulty"
                ></TagChip>
                <v-chip v-else size="small" label title="Difficulté"
                  >Aucune</v-chip
                >
              </v-col>
            </v-row>
          </v-list-item>
        </v-list>
      </v-card-text>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import type { VariantG } from "@/controller/editor";
import { $ref } from "vue/macros";
import TagChip from "./utils/TagChip.vue";

interface Props {
  variants: VariantG[];
  readonly: boolean;
  modelValue: number; // index into variants
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: number): void;
  (e: "delete", variant: VariantG): void;
  (e: "duplicate", variant: VariantG): void;
}>();

const showMenu = $ref(false);
</script>

<style scoped></style>
