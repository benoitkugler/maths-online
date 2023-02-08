<template>
  <v-card :title="'Notes de la classe ' + props.classroom.name">
    <v-card-text>
      <v-table>
        <tr>
          <th class="py-2">Elève</th>
          <th v-for="sheet in props.sheets" :key="sheet.Sheet.Id">
            <div class="bg-blue-lighten-4 rounded mx-2 py-1">
              {{ sheet.Sheet.Title }} ( /20)
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
              v-for="sheet in props.sheets"
              :key="sheet.Sheet.Id"
            >
              {{ getMark(sheet, student) }}
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
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { computed } from "@vue/reactivity";
import { onMounted } from "vue";
import { $ref } from "vue/macros";

interface Props {
  classroom: Classroom;
  sheets: SheetExt[];
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
    IdSheets: props.sheets.map((sh) => sh.Sheet.Id),
  });
  if (res == undefined) return;
  data = res;
}

function _getMark(sheet: SheetExt, student: StudentHeader) {
  const sheetMarks = (data?.Marks || {})[sheet.Sheet.Id] || {};
  return sheetMarks[student.Id] || 0;
}

function getMark(sheet: SheetExt, student: StudentHeader) {
  return _getMark(sheet, student).toFixed(1);
}

function getMoyenne(student: StudentHeader) {
  let total = 0;
  props.sheets.forEach((sh) => (total += _getMark(sh, student)));
  return (total / props.sheets.length).toFixed(2);
}
</script>
