<template>
  <v-card :title="title" :subtitle="subtitle">
    <v-card-text>
      <v-list density="compact">
        <v-list-item v-if="!props.sheet.Tasks?.length">
          Aucune tâche.
        </v-list-item>
        <v-list-item v-for="(task, index) in props.sheet.Tasks" :key="index">
          <v-list-item-title class="mr-2">
            {{ task.Title }}
          </v-list-item-title>
          <template v-slot:append> / {{ taskBareme(task) }} </template>
        </v-list-item>
      </v-list>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { SheetExt } from "@/controller/api_gen";
import { taskBareme } from "@/controller/utils";
import { computed } from "vue";

interface Props {
  sheet: SheetExt;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "delete"): void;
}>();

const title = computed(() => `Feuille : ${props.sheet.Sheet.Title}`);
const subtitle = computed(() => `Niveau : ${props.sheet.Sheet.Level}`);
</script>
