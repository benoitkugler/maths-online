<template>
  <v-menu offset-y :close-on-content-click="false" v-model="showMenu">
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
            v-for="(question, index) in variants"
            @click="
              emit('update:model-value', index);
              showMenu = false;
            "
            :key="question.Id"
          >
            <v-row>
              <v-col align-self="center" cols="auto">
                <v-btn
                  v-if="!props.readonly"
                  class="ml-1"
                  size="small"
                  icon
                  @click.stop="emit('delete', question)"
                  title="Supprimer la variante"
                >
                  <v-icon icon="mdi-delete" color="red" size="small"></v-icon>
                </v-btn>

                <v-btn
                  v-if="!props.readonly"
                  class="mx-1"
                  size="small"
                  icon
                  @click.stop="emit('duplicate', question)"
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
                ({{ question.Id }}) {{ question.Subtitle || "..." }}
              </v-col>
              <v-col align-self="center" style="text-align: right">
                <TagChip
                  v-if="question.Difficulty"
                  :tag="question.Difficulty"
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
import type { Question } from "@/controller/api_gen";
import { $ref } from "vue/macros";
import TagChip from "../utils/TagChip.vue";

interface Props {
  variants: Question[];
  readonly: boolean;
  modelValue: number; // index into variants
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: number): void;
  (e: "delete", question: Question): void;
  (e: "duplicate", question: Question): void;
}>();

const showMenu = $ref(false);
</script>

<style scoped></style>
