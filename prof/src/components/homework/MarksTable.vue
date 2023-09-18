<template>
  <v-table>
    <tr>
      <th class="py-2 text-left">Elève</th>
      <th>
        <div class="bg-blue-lighten-1 rounded py-1 px-1">Moyenne</div>
      </th>
      <th v-for="tr in props.travaux" :key="tr.Id">
        <div class="bg-blue-lighten-4 rounded mx-2 my-1 py-1 text-subtitle-1">
          {{ props.sheets.get(tr.IdSheet)!.Sheet.Title }}
        </div>
      </th>
    </tr>
    <tr
      v-for="(student, index) in data?.Students"
      :key="index"
      :class="{ 'bg-grey-lighten-2': index % 2 == 0 }"
    >
      <td class="px-2 text-left">{{ student.Label }}</td>
      <td class="pa-1 text-center font-weight-bold">
        {{ getMoyenne(student) }}
      </td>
      <td class="text-center" v-for="tr in props.travaux" :key="tr.Id">
        <MarksTableCell :data="getMark(tr, student)"></MarksTableCell>
      </td>
    </tr>
  </v-table>
</template>

<script setup lang="ts">
import type {
  HomeworkMarksOut,
  SheetExt,
  StudentHeader,
  Travail,
  StudentTravailMark
} from "@/controller/api_gen";
import MarksTableCell from "./MarksTableCell.vue";

interface Props {
  data: HomeworkMarksOut;
  travaux: Travail[];
  sheets: Map<number, SheetExt>;
}

const props = defineProps<Props>();

function getMark(tr: Travail, student: StudentHeader) {
  const sheetMarks = (props.data?.Marks || {})[tr.Id];
  const mark: StudentTravailMark = (sheetMarks.Marks || {})[student.Id] || {
    Mark: 0,
    Dispensed: false,
    NbTries: 0
  };
  return mark;
}

function getMoyenne(student: StudentHeader) {
  let total = 0;
  let nbTravaux = 0;
  props.travaux.forEach(tr => {
    const m = getMark(tr, student);
    if (m.Dispensed) {
      return;
    }
    total += m.Mark;
    nbTravaux += 1;
  });
  if (nbTravaux == 0) return "-";
  return (total / nbTravaux).toFixed(2);
}
</script>
