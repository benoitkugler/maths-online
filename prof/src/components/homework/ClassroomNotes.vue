<template>
  <v-card :title="'Notes de la classe ' + props.classroom.name">
    <v-card-text>
      <v-table>
        <tr>
          <th class="py-2">Elève</th>
          <th v-for="tr in props.travaux" :key="tr.Id">
            <div class="bg-blue-lighten-4 rounded mx-2 py-1">
              {{ props.sheets.get(tr.IdSheet)!.Sheet.Title }} ( /20)
            </div>
          </th>
          <th>
            <div class="bg-blue-lighten-1 rounded py-1">Moyenne ( /20)</div>
          </th>
        </tr>
        <template v-if="isLoading"> Chargement des notes... </template>
        <template v-else>
          <tr
            v-for="(student, index) in data?.Students"
            :key="index"
            :class="{ 'bg-grey-lighten-2': index % 2 == 0 }"
          >
            <th style="text-align: left" class="px-2">{{ student.Label }}</th>
            <td
              style="text-align: center"
              v-for="tr in props.travaux"
              :key="tr.Id"
            >
              {{ getMark(tr, student) }}
            </td>
            <td style="text-align: right" class="pa-1">
              {{ getMoyenne(student) }}
            </td>
          </tr>
        </template>
      </v-table>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type {
  Classroom,
  HomeworkMarksOut,
  SheetExt,
  StudentHeader,
  Travail
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { computed } from "@vue/reactivity";
import { onMounted } from "vue";
import { $ref } from "vue/macros";

interface Props {
  classroom: Classroom;
  travaux: Travail[];
  sheets: Map<number, SheetExt>;
}

const props = defineProps<Props>();

// const emit = defineEmits<{}>();

onMounted(fetchNotes);

let isLoading = computed(() => data == null);
let data = $ref<HomeworkMarksOut | null>(null);

async function fetchNotes() {
  data = null;
  const res = await controller.HomeworkGetMarks({
    IdClassroom: props.classroom.id,
    IdTravaux: props.travaux.map(tr => tr.Id)
  });
  if (res == undefined) return;
  data = res;
}

function _getMark(tr: Travail, student: StudentHeader) {
  const sheetMarks = (data?.Marks || {})[tr.Id];
  const mark = (sheetMarks.Marks || {})[student.Id] || 0;
  const ignored = (sheetMarks.Ignored || []).includes(student.Id);
  return { mark, ignored };
}

function getMark(tr: Travail, student: StudentHeader) {
  const m = _getMark(tr, student);
  if (m.ignored) {
    return `${m.mark.toFixed(1)} (*)`;
  }
  return m.mark.toFixed(1);
}

function getMoyenne(student: StudentHeader) {
  let total = 0;
  let nbTravaux = 0;
  props.travaux.forEach(tr => {
    const m = _getMark(tr, student);
    if (m.ignored) {
      return;
    }
    total += m.mark;
    nbTravaux += 1;
  });
  if (nbTravaux == 0) return "-";
  return (total / nbTravaux).toFixed(2);
}
</script>
