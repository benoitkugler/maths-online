<template>
  <v-card
    :title="'Résultats de la classe : ' + props.classroom.name"
    :subtitle="
      viewKind == 'marks'
        ? 'Les notes affichées sont /20.'
        : 'Nombre de tentatives réussies et échouées, pour la classe.'
    "
  >
    <template v-slot:append>
      <div style="min-width: 300px">
        <v-select
          v-model="viewKind"
          label="Statistique"
          :items="viewItems"
          density="compact"
          variant="outlined"
          hide-details
        ></v-select>
      </div>
    </template>
    <v-card-text>
      <v-fade-transition :hide-on-leave="true">
        <template v-if="isLoading"> Chargement des notes... </template>
        <MarksTable
          :data="data!"
          :travaux="props.travaux"
          :sheets="props.sheets"
          v-else-if="viewKind == 'marks'"
        ></MarksTable>
        <MarksStats
          :data="data!"
          :travaux="props.travaux"
          :sheets="props.sheets"
          v-else-if="viewKind == 'tasks'"
        ></MarksStats>
      </v-fade-transition>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type {
  Classroom,
  HomeworkMarksOut,
  SheetExt,
  Travail,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { ref, computed } from "vue";
import { onMounted } from "vue";
import MarksStats from "./MarksStats.vue";
import MarksTable from "./MarksTable.vue";

interface Props {
  classroom: Classroom;
  travaux: Travail[];
  sheets: Map<number, SheetExt>;
}

const props = defineProps<Props>();

// const emit = defineEmits<{}>();

const viewKind = ref<"marks" | "tasks">("marks");

const viewItems = [
  { title: "Notes par élève", value: "marks" },
  { title: "Réussite par exercice", value: "tasks" },
];

onMounted(fetchNotes);

const isLoading = computed(() => data.value == null);
const data = ref<HomeworkMarksOut | null>(null);

async function fetchNotes() {
  data.value = null;
  const res = await controller.HomeworkGetMarks({
    IdClassroom: props.classroom.id,
    IdTravaux: props.travaux.map((tr) => tr.Id),
  });
  if (res == undefined) return;
  data.value = res;
}
</script>
