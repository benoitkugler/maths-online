<template>
  <v-card title="Démarrer la session" width="870px">
    <v-card-text class="mt-2">
      <v-select
        variant="outlined"
        density="comfortable"
        label="Formation des groupes"
        :items="strategyItems"
        :model-value="strategyItems[launchOptions.Kind]"
        @update:model-value="
          (s) => (launchOptions.Kind = strategyItems.indexOf(s))
        "
        :hint="hint"
        persistent-hint
      >
      </v-select>
      <component
        v-model="launchOptions.Data"
        :is="groupStrategyComponent"
      ></component>

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-col cols="auto" class="text-right">
          <v-btn
            block
            @click="emit('launch', launchOptions)"
            :disabled="!isValid"
            color="success"
            variant="contained"
          >
            Lancer la session
          </v-btn>
        </v-col>
      </v-card-actions>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type {
  FixedSizeGroupStrategy,
  GroupStrategy,
  RandomGroupStrategy,
} from "@/controller/api_gen";
import { GroupStrategyKind } from "@/controller/api_gen";
import { computed, markRaw } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import FixedSizeGroupsVue from "./FixedSizeGroups.vue";
import RandomGroupsVue from "./RandomGroups.vue";

// interface Props {}
// const props = defineProps<Props>();
//

const emit = defineEmits<{
  (e: "launch", v: GroupStrategy): void;
}>();

let launchOptions = $ref<GroupStrategy>({
  Kind: GroupStrategyKind.FixedSizeGroupStrategy,
  Data: { Groups: [] },
});
// let launchOptions = $ref<GroupStrategy>({
//   Kind: GroupStrategyKind.RandomGroupStrategy,
//   Data: { MaxPlayersPerGroup: 4, TotalPlayersNumber: 0 }
// });

const groupStrategyComponent = computed(() => {
  switch (launchOptions.Kind) {
    case GroupStrategyKind.RandomGroupStrategy:
      return markRaw(RandomGroupsVue);
    case GroupStrategyKind.FixedSizeGroupStrategy:
      return markRaw(FixedSizeGroupsVue);
  }
});

const strategyItems = [
  "Taille des groupes fixées",
  "Création aléatoire des groupes",
];

const hint = computed(() => {
  return [
    "Le nombre et la taille de chaque groupe est fixée au lancement. Chaque groupe est accessible par un code différent.",
    "Les groupes sont formés au fur et à mesure des connexions. Un seul code permet d'accéder à la partie.",
  ][launchOptions.Kind];
});

const isValid = computed(() => isGroupValid(launchOptions));

function isGroupValid(options: GroupStrategy) {
  switch (options.Kind) {
    case GroupStrategyKind.RandomGroupStrategy: {
      const d = options.Data as RandomGroupStrategy;
      return d.MaxPlayersPerGroup >= 1 && d.TotalPlayersNumber >= 1;
    }
    case GroupStrategyKind.FixedSizeGroupStrategy: {
      const data = options.Data as FixedSizeGroupStrategy;
      return data.Groups?.length && data.Groups.every((v) => v > 0);
    }
  }
}
</script>

<style scoped></style>
