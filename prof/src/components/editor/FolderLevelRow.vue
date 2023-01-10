<template>
  <v-card class="bg-pink-lighten-4">
    <v-card-title> {{ title }} </v-card-title>
    <v-card-text>
      <v-row>
        <v-col
          v-for="(chapter, index) in props.level.Chapters"
          :key="index"
          cols="12"
          sm="6"
          md="4"
          lg="3"
        >
          <v-card
            :subtitle="chapter.Chapter || 'Non classé'"
            @click="emit('clicked', chapter.Chapter)"
          >
            <v-card-text style="text-align: center">
              <v-chip>{{ chapter.GroupCount }} élément(s)</v-chip>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { LevelTagLabels, type LevelItems } from "@/controller/api_gen";
import { computed } from "vue";

interface Props {
  level: LevelItems;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "clicked", chapter: string): void;
}>();

const title = computed(() => LevelTagLabels[props.level.Level] || "Non classé");
</script>

<style scoped></style>
