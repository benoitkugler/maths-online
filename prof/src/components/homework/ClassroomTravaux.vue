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

  <div class="ma-2">
    <v-row no-gutters>
      <v-col>
        <!-- <v-card-title>{{ classroom.Classroom.name }}</v-card-title> -->
      </v-col>
      <v-col align-self="center" cols="auto">
        <v-btn
          density="compact"
          @click="onShowNotes"
          title="Afficher les notes pour les feuilles sélectionnées"
          :disabled="
            !classroom.Travaux?.length ||
            (inSelect && Array.from(selectedTravaux).length == 0)
          "
        >
          {{ inSelect ? "Afficher" : "Voir les notes..." }}
        </v-btn>
      </v-col>
    </v-row>
    <v-row v-if="!classroom.Travaux?.length" class="mt-6">
      <v-col align-self="center" style="text-align: center">
        <i
          >Déplacer et déposer une feuille ici pour définir un nouveau
          travail...</i
        >
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
        ></travail-card>
        <v-card
          v-else
          :color="selectedTravaux.has(travail.Id) ? 'blue' : 'grey-lighten-3'"
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
  SheetExt,
  Travail,
} from "@/controller/api_gen";
import { $computed, $ref } from "vue/macros";
import ClassroomNotes from "./ClassroomNotes.vue";
import TravailCard from "./TravailCard.vue";
import { computed } from "vue";

interface Props {
  classroom: ClassroomTravaux;
  sheets: Map<number, SheetExt>;
  classrooms: Classroom[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update", travail: Travail): void;
  (e: "delete", travail: Travail): void;
  (e: "copy", travail: Travail, idClassroom: number): void;
}>();

let inSelect = $ref(false);
let selectedTravaux = $ref<Set<number>>(new Set());

let showNotes = $ref(false);

function onToggle(tr: Travail) {
  if (selectedTravaux.has(tr.Id)) {
    selectedTravaux.delete(tr.Id);
  } else {
    selectedTravaux.add(tr.Id);
  }
}
function onShowNotes() {
  if (!inSelect) {
    // start with all selected
    selectedTravaux = new Set(props.classroom.Travaux?.map((tr) => tr.Id));
    inSelect = true;
  } else {
    inSelect = false;
    showNotes = true;
  }
}
</script>
