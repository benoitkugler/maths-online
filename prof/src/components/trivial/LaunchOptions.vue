<template>
  <v-card title="Démarrer la session">
    <v-card-text class="mt-2">
      <v-row>
        <v-col>
          <v-select
            variant="outlined"
            density="compact"
            :items="strategyItems"
            label="Type de lancement"
            :hint="hint"
            v-model="launchOptions.Kind"
            persistent-hint
          ></v-select>
        </v-col>
      </v-row>

      <groups-auto
        v-if="isAuto"
        :model-value="(launchOptions.Data as GroupsStrategyAuto)"
        @update:model-value="(v) => (launchOptions.Data = v)"
      >
      </groups-auto>
      <groups-manual
        v-else
        :model-value="(launchOptions.Data as GroupsStrategyManual)"
        @update:model-value="(v) => (launchOptions.Data = v)"
      >
      </groups-manual>

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
import {
  GroupsStrategyKind,
  type GroupsStrategy,
  type GroupsStrategyAuto,
  type GroupsStrategyManual,
} from "@/controller/api_gen";
import { computed } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import GroupsAuto from "./GroupsAuto.vue";
import GroupsManual from "./GroupsManual.vue";

// interface Props {}
// const props = defineProps<Props>();
//

const emit = defineEmits<{
  (e: "launch", groups: GroupsStrategy): void;
}>();

let launchOptions = $ref<GroupsStrategy>({
  Kind: GroupsStrategyKind.GroupsStrategyAuto,
  Data: { Groups: [] },
});

const isAuto = computed(
  () => launchOptions.Kind == GroupsStrategyKind.GroupsStrategyAuto
);

const hint = computed(() =>
  isAuto
    ? "Le nombre et la taille de chaque groupe est fixée au lancement, et chaque partie démarre automatiquement."
    : "La taille des groupes n'est pas spécifiée, et chaque partie doit être démarrée manuellement."
);
const strategyItems = [
  {
    value: GroupsStrategyKind.GroupsStrategyAuto,
    title: "Automatique",
  },
  { value: GroupsStrategyKind.GroupsStrategyManual, title: "Manuel" },
];

const isValid = computed(() => {
  switch (launchOptions.Kind) {
    case GroupsStrategyKind.GroupsStrategyAuto: {
      const groups = (launchOptions.Data as GroupsStrategyAuto).Groups;
      return groups?.length && groups.every((v) => v > 0);
    }
    case GroupsStrategyKind.GroupsStrategyManual:
      return (launchOptions.Data as GroupsStrategyManual).NbGroups > 0;
  }
});
</script>

<style scoped></style>
