<template>
  <v-card title="Démarrer la session">
    <v-card-text class="mt-2">
      <v-row>
        <v-col>
          <v-select
            label="Type de lancement"
            density="compact"
            variant="outlined"
            :items="strategyItems"
            v-model="launchOptions.Kind"
            :hint="hint"
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
import { ref, computed } from "vue";
import GroupsAuto from "./GroupsAuto.vue";
import GroupsManual from "./GroupsManual.vue";

// interface Props {}
// const props = defineProps<Props>();
//

const emit = defineEmits<{
  (e: "launch", groups: GroupsStrategy): void;
}>();

const launchOptions = ref<GroupsStrategy>({
  Kind: GroupsStrategyKind.GroupsStrategyAuto,
  Data: { Groups: [] },
});

const isAuto = computed(
  () => launchOptions.value.Kind == GroupsStrategyKind.GroupsStrategyAuto
);

const hint = computed(() =>
  isAuto.value
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
  switch (launchOptions.value.Kind) {
    case GroupsStrategyKind.GroupsStrategyAuto: {
      const groups = (launchOptions.value.Data as GroupsStrategyAuto).Groups;
      return groups?.length && groups.every((v) => v > 0);
    }
    case GroupsStrategyKind.GroupsStrategyManual:
      return (launchOptions.value.Data as GroupsStrategyManual).NbGroups > 0;
    default:
      throw "";
  }
});
</script>

<style scoped></style>
