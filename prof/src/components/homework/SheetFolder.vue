<template>
  <v-card>
    <v-row>
      <v-col><v-card-title>Feuilles d'exercices</v-card-title></v-col>
      <v-col cols="auto" align-self="center">
        <v-btn class="mx-2" @click="emit('create')">
          <v-icon color="green">mdi-plus</v-icon>
          Créer une feuille</v-btn
        >
      </v-col>
    </v-row>
    <v-card-text>
      <v-expansion-panels>
        <v-expansion-panel v-for="level in byLevels" :key="level[0]">
          <v-expansion-panel-title class="bg-pink-lighten-3">
            {{ level[0] || "Non classé" }}
          </v-expansion-panel-title>
          <v-expansion-panel-text>
            <v-list density="compact" class="py-0">
              <v-list-item
                v-for="sheet in level[1]"
                :key="sheet.Sheet.Id"
                style="cursor: grab"
                @dragstart="(e: DragEvent) => onDrapStart(e, sheet)"
                draggable="true"
                color="grey"
              >
                <v-list-item-title>{{ sheet.Sheet.Title }}</v-list-item-title>

                <template v-slot:append="{}">
                  <v-list-item-action>
                    <v-btn
                      class="mx-1"
                      icon="mdi-pencil"
                      size="x-small"
                      @click="emit('edit', sheet)"
                    ></v-btn>
                    <v-btn
                      class="mx-1"
                      size="x-small"
                      icon
                      @click="emit('duplicate', sheet)"
                    >
                      <v-icon color="secondary">mdi-content-copy</v-icon>
                    </v-btn>
                    <v-btn
                      class="mx-1"
                      size="x-small"
                      icon
                      @click="emit('delete', sheet)"
                    >
                      <v-icon color="red">mdi-delete</v-icon>
                    </v-btn>
                  </v-list-item-action>
                </template>
              </v-list-item>
            </v-list>
          </v-expansion-panel-text>
        </v-expansion-panel>
      </v-expansion-panels>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { SheetExt } from "@/controller/api_gen";
import { computed } from "vue";

interface Props {
  sheets: Map<number, SheetExt>;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "create"): void;
  (e: "edit", sheet: SheetExt): void;
  (e: "duplicate", sheet: SheetExt): void;
  (e: "delete", sheet: SheetExt): void;
}>();

const byLevels = computed(() => {
  const tmp = new Map<string, SheetExt[]>();
  for (const sh of props.sheets.values()) {
    const l = tmp.get(sh.Sheet.Level) || [];
    tmp.set(sh.Sheet.Level, l.concat(sh));
  }
  const out = Array.from(tmp.entries());
  out.sort((a, b) => a[0].localeCompare(b[0]));
  out.forEach((v) =>
    v[1].sort((a, b) => a.Sheet.Title.localeCompare(b.Sheet.Title))
  );
  return out;
});

function onDrapStart(payload: DragEvent, sheet: SheetExt) {
  if (payload.dataTransfer == null) return;
  payload.dataTransfer.setData("text/json", JSON.stringify({ sheet }));
  payload.dataTransfer.dropEffect = "copy";
}
</script>
