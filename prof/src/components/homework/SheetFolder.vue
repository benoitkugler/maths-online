<template>
  <v-card
    title="Feuilles d'exercices"
    subtitle="Les feuilles ci-dessous peuvent être partagées entre plusieurs classes
          et sont conservées si une classe est supprimée."
  >
    <template v-slot:append>
      <v-btn class="mx-2" @click="emit('create')">
        <v-icon color="green">mdi-plus</v-icon>
        Créer une feuille
      </v-btn>
      <MatiereSelect
        :matiere="props.matiere"
        @update:matiere="v => emit('update:matiere', v)"
      ></MatiereSelect>
    </template>

    <v-card-text>
      <v-expansion-panels>
        <v-expansion-panel v-for="level in byLevels" :key="level[0]">
          <v-expansion-panel-title class="bg-pink-lighten-3">
            {{ level[0] || "Non classé" }}
          </v-expansion-panel-title>
          <v-expansion-panel-text>
            <v-list density="comfortable" class="py-0" style="column-count: 2">
              <SheetCard
                v-for="sheet in level[1]"
                :key="sheet.Sheet.Id"
                :sheet="sheet"
                :classrooms="props.classrooms"
                @assign="target => emit('assign', sheet.Sheet.Id, target)"
                @edit="emit('edit', sheet)"
                @duplicate="emit('duplicate', sheet)"
                @delete="emit('delete', sheet)"
                @updatePublic="b => emit('updatePublic', sheet.Sheet, b)"
                @createReview="emit('createReview', sheet.Sheet)"
              ></SheetCard>
            </v-list>
          </v-expansion-panel-text>
        </v-expansion-panel>
      </v-expansion-panels>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type {
  Classroom,
  SheetExt,
  Sheet,
  MatiereTag
} from "@/controller/api_gen";
import { computed } from "vue";
import SheetCard from "./SheetCard.vue";
import MatiereSelect from "../MatiereSelect.vue";

interface Props {
  sheets: Map<number, SheetExt>;
  classrooms: Classroom[];
  matiere: MatiereTag;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "create"): void;
  (e: "assign", idSheet: number, idClassroom: number): void;
  (e: "edit", sheet: SheetExt): void;
  (e: "duplicate", sheet: SheetExt): void;
  (e: "delete", sheet: SheetExt): void;
  (e: "updatePublic", sheet: Sheet, pub: boolean): void;
  (e: "createReview", sheet: Sheet): void;
  (e: "update:matiere", matiere: MatiereTag): void;
}>();

const byLevels = computed(() => {
  const tmp = new Map<string, SheetExt[]>();
  for (const sh of props.sheets.values()) {
    if (sh.Sheet.Anonymous.Valid) continue; // do not show anonymous sheets
    const l = tmp.get(sh.Sheet.Level) || [];
    tmp.set(sh.Sheet.Level, l.concat(sh));
  }
  const out = Array.from(tmp.entries());
  out.sort((a, b) => a[0].localeCompare(b[0]));
  out.forEach(v =>
    v[1].sort((a, b) => a.Sheet.Title.localeCompare(b.Sheet.Title))
  );
  return out;
});
</script>
