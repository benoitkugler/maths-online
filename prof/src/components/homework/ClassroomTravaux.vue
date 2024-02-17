<template>
  <v-dialog v-model="showNotes">
    <classroom-notes
      v-if="showNotes"
      :classroom="props.classroom.Classroom"
      :travaux="
        props.classroom.Travaux?.filter((tr) => selectedTravaux.has(tr.Id)) ||
        []
      "
      :sheets="props.sheets"
    ></classroom-notes>
  </v-dialog>

  <v-dialog
    max-width="800px"
    :model-value="showDispensesFor != null"
    @update:model-value="showDispensesFor = null"
  >
    <travail-dispenses
      v-if="showDispensesFor != null"
      :travail="showDispensesFor"
      :sheet="props.sheets.get(showDispensesFor.IdSheet)!.Sheet"
    ></travail-dispenses>
  </v-dialog>

  <div class="ma-2">
    <v-row no-gutters>
      <v-col>
        <!-- <v-card-title>{{ classroom.Classroom.name }}</v-card-title> -->
      </v-col>
      <v-col align-self="center" cols="auto">
        <v-btn
          density="comfortable"
          title="Afficher les notes pour les feuilles sélectionnées"
          class="mx-1"
          @click="onShowNotes"
          :disabled="
            !classroom.Travaux?.length ||
            (inSelect && Array.from(selectedTravaux).length == 0)
          "
        >
          {{ inSelect ? "Afficher" : "Voir les notes..." }}
        </v-btn>
        <v-btn density="comfortable" class="mx-1" @click="emit('create')">
          <v-icon color="green">mdi-plus</v-icon>
          Ajouter
        </v-btn>
      </v-col>
    </v-row>
    <v-row v-if="!classroom.Travaux?.length" class="mt-6 mb-3">
      <v-col align-self="center" style="text-align: center">
        <i>
          Aucun travail. <br />
          Vous pouvez ajouter un nouveau travail depuis les favoris.
        </i>
      </v-col>
    </v-row>
    <v-row class="mt-1">
      <v-col
        lg="6"
        cols="12"
        v-for="(travail, index) in props.classroom.Travaux"
        :key="index"
      >
        <travail-card
          v-if="!inSelect"
          :travail="travail"
          :sheet="props.sheets.get(travail.IdSheet)!"
          :classrooms="props.classrooms"
          @update="(tr) => emit('update', tr)"
          @delete="emit('delete', travail)"
          @copy="(idClassroom) => emit('copy', travail, idClassroom)"
          @set-favorite="(s) => emit('setFavorite', s)"
          @edit-sheet="(s) => emit('editSheet', s)"
          @show-dispenses="showDispensesFor = travail"
        ></travail-card>
        <v-card
          v-else
          :color="
            selectedTravaux.has(travail.Id)
              ? 'blue-lighten-2'
              : 'grey-lighten-3'
          "
          @click="onToggle(travail)"
        >
          <v-card-text style="text-align: center">
            <v-row>
              <v-col cols="2">
                <v-icon
                  :color="selectedTravaux.has(travail.Id) ? '' : 'transparent'"
                >
                  mdi-check</v-icon
                >
              </v-col>
              <v-col>
                {{ props.sheets.get(travail.IdSheet)!.Sheet.Title }}
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script setup lang="ts">
import type {
  Classroom,
  ClassroomTravaux,
  Int,
  Sheet,
  SheetExt,
  Travail,
} from "@/controller/api_gen";
import ClassroomNotes from "./ClassroomNotes.vue";
import TravailCard from "./TravailCard.vue";
import { ref } from "vue";
import TravailDispenses from "./TravailDispenses.vue";

interface Props {
  classroom: ClassroomTravaux;
  sheets: Map<number, SheetExt>;
  classrooms: Classroom[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "create"): void;
  (e: "update", travail: Travail): void;
  (e: "delete", travail: Travail): void;
  (e: "copy", travail: Travail, idClassroom: Int): void;
  (e: "setFavorite", sheet: Sheet): void;
  (e: "editSheet", sheet: SheetExt): void;
}>();

const inSelect = ref(false);
const selectedTravaux = ref<Set<number>>(new Set());

const showNotes = ref(false);

function onToggle(tr: Travail) {
  if (selectedTravaux.value.has(tr.Id)) {
    selectedTravaux.value.delete(tr.Id);
  } else {
    selectedTravaux.value.add(tr.Id);
  }
}
function onShowNotes() {
  if (!inSelect.value) {
    // start with all noted selected
    selectedTravaux.value = new Set(
      props.classroom.Travaux?.filter((tr) => tr.Noted).map((tr) => tr.Id)
    );
    inSelect.value = true;
  } else {
    inSelect.value = false;
    showNotes.value = true;
  }
}

const showDispensesFor = ref<Travail | null>(null);
</script>
