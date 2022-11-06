<template>
  <v-card title="Démarrer la session">
    <v-card-text class="mt-2">
      {{ hint }}
      <fixed-size-groups v-model="launchOptions"> </fixed-size-groups>

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-col cols="auto" class="text-right">
          <v-btn
            block
            @click="emit('launch', launchOptions)"
            :disabled="!isValid"
            color="success"
            variant="outlined"
          >
            Lancer la session
          </v-btn>
        </v-col>
      </v-card-actions>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { computed } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import { default as FixedSizeGroups } from "./FixedSizeGroups.vue";

// interface Props {}
// const props = defineProps<Props>();
//

const emit = defineEmits<{
  (e: "launch", groups: number[]): void;
}>();

let launchOptions = $ref<number[]>([]);

const hint =
  "Le nombre et la taille de chaque groupe est fixée au lancement. Chaque groupe est accessible par un code différent.";

const isValid = computed(() => isGroupValid(launchOptions));

function isGroupValid(groups: number[]) {
  return groups?.length && groups.every((v) => v > 0);
}
</script>

<style scoped></style>
