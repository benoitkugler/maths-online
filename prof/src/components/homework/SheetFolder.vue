<template>
  <v-card>
    <v-row class="pa-2">
      <v-col>
        <v-card-title>Feuilles d'exercices</v-card-title>
        <v-card-subtitle>
          Les feuilles ci-dessous peuvent être partagées entre plusieurs classes
          et sont conservées si une classe est supprimée.
        </v-card-subtitle>
      </v-col>
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
            <v-list density="comfortable" class="py-0" style="column-count: 2">
              <v-list-item
                class="bg-grey-lighten-3 py-2 mb-1 px-2"
                v-for="sheet in level[1]"
                :key="sheet.Sheet.Id"
                style="break-inside: avoid-column"
                color="grey"
                rounded
              >
                <v-menu>
                  <template v-slot:activator="{ isActive, props }">
                    <v-card
                      v-on="{ isActive }"
                      v-bind="props"
                      class="pa-2 mr-2"
                    >
                      <v-list-item-title>
                        {{ sheet.Sheet.Title }}
                      </v-list-item-title>
                      <v-list-item-subtitle>
                        {{ subtitle(sheet) }}
                      </v-list-item-subtitle>
                    </v-card>
                  </template>
                  <PreviewSheet :sheet="sheet"></PreviewSheet>
                </v-menu>
                <template v-slot:append="{}">
                  <v-list-item-action>
                    <v-menu offset-y close-on-content-click>
                      <template v-slot:activator="{ isActive, props: props2 }">
                        <v-btn
                          v-on="{ isActive }"
                          v-bind="props2"
                          variant="outlined"
                          color="primary-darken-1"
                          title="Ajouter à une classe"
                        >
                          Assigner
                        </v-btn>
                      </template>
                      <v-list>
                        <v-list-subheader v-if="!props.classrooms.length">
                          Aucune classe.
                        </v-list-subheader>
                        <v-list-subheader v-else
                          >Assigner cette feuille à...</v-list-subheader
                        >
                        <v-list-item
                          v-for="(classroom, index) in props.classrooms"
                          :key="index"
                          link
                          @click="emit('assign', sheet.Sheet.Id, classroom.id)"
                        >
                          {{ classroom.name }}
                        </v-list-item>
                      </v-list>
                    </v-menu>
                    <!--  -->
                    <v-btn
                      v-if="sheet.Origin.Visibility == Visibility.Admin"
                      @click="emit('duplicate', sheet)"
                      title="Dupliquer"
                      icon
                      class="ml-4"
                      size="small"
                    >
                      <v-icon color="secondary"> mdi-content-copy </v-icon>
                    </v-btn>

                    <v-menu v-else offset-y close-on-content-click>
                      <template v-slot:activator="{ isActive, props }">
                        <v-btn
                          v-on="{ isActive }"
                          v-bind="props"
                          icon
                          class="ml-4"
                          size="small"
                        >
                          <v-icon>mdi-dots-vertical</v-icon>
                        </v-btn>
                      </template>
                      <v-list density="compact">
                        <v-list-item>
                          <v-btn
                            variant="flat"
                            class="mr-2"
                            size="small"
                            @click="emit('edit', sheet)"
                          >
                            <template v-slot:prepend>
                              <v-icon icon="mdi-pencil" class="mr-4"></v-icon>
                            </template>
                            Editer
                          </v-btn>
                        </v-list-item>
                        <v-list-item>
                          <v-btn
                            variant="flat"
                            class="mr-2"
                            size="small"
                            @click="emit('duplicate', sheet)"
                          >
                            <template v-slot:prepend>
                              <v-icon color="secondary" class="mr-4">
                                mdi-content-copy
                              </v-icon>
                            </template>
                            Dupliquer
                          </v-btn>
                        </v-list-item>
                        <v-list-item>
                          <v-btn
                            variant="flat"
                            class="mr-2"
                            size="small"
                            @click="emit('delete', sheet)"
                          >
                            <template v-slot:prepend>
                              <v-icon color="red" class="mr-4"
                                >mdi-delete</v-icon
                              >
                            </template>
                            Supprimer
                          </v-btn>
                        </v-list-item>

                        <v-list-item>
                          <OriginButton
                            variant="text"
                            :origin="sheet.Origin"
                            @update-public="
                              b => emit('updatePublic', sheet.Sheet, b)
                            "
                            @create-review="emit('createReview', sheet.Sheet)"
                          ></OriginButton>
                        </v-list-item>
                      </v-list>
                    </v-menu>
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
import {
  Visibility,
  type Classroom,
  type SheetExt,
  type Sheet
} from "@/controller/api_gen";
import { computed } from "vue";
import OriginButton from "../OriginButton.vue";
import PreviewSheet from "./PreviewSheet.vue";

interface Props {
  sheets: Map<number, SheetExt>;
  classrooms: Classroom[];
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

function subtitle(sheet: SheetExt) {
  if (sheet.NbTravaux == 0) {
    return "Inactive";
  } else if (sheet.NbTravaux == 1) {
    return "Active dans 1 travail";
  } else {
    return `Active dans ${sheet.NbTravaux} travaux`;
  }
}
</script>
