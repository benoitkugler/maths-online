<template>
  <v-table>
    <tr>
      <th class="py-2">Elève</th>
      <th v-for="tr in props.travaux" :key="tr.Id">
        <div class="bg-blue-lighten-4 rounded mx-2 my-1 py-1 text-subtitle-1">
          {{ props.sheets.get(tr.IdSheet)!.Sheet.Title }}
        </div>
      </th>
      <th>
        <div class="bg-blue-lighten-1 rounded py-1 px-1">Moyenne</div>
      </th>
    </tr>
    <tr
      v-for="(student, index) in data?.Students"
      :key="index"
      :class="{ 'bg-grey-lighten-2': index % 2 == 0 }"
    >
      <td style="text-align: left" class="px-2">{{ student.Label }}</td>
      <td style="text-align: center" v-for="tr in props.travaux" :key="tr.Id">
        {{ getMark(tr, student) }}
      </td>
      <td style="text-align: right" class="pa-1">
        {{ getMoyenne(student) }}
      </td>
    </tr>
  </v-table>
</template>

<script setup lang="ts">
import type {
  HomeworkMarksOut,
  SheetExt,
  StudentHeader,
  Travail
} from "@/controller/api_gen";

interface Props {
  data: HomeworkMarksOut;
  travaux: Travail[];
  sheets: Map<number, SheetExt>;
}

const props = defineProps<Props>();

function _getMark(tr: Travail, student: StudentHeader) {
  const sheetMarks = (props.data?.Marks || {})[tr.Id];
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
