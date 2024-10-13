<template>
  <v-card title="Progression des élèves">
    <template v-slot:append>
      <v-select
        variant="outlined"
        density="comfortable"
        label="Classe"
        :items="props.classrooms.map((c) => ({ title: c.name, value: c.id }))"
        v-model="selected"
        no-data-text="Vous n'avez aucune classe."
        @update:model-value="fetchAdvance"
      ></v-select>
    </template>
    <v-card-text>
      <v-table>
        <tr>
          <th width="150px">Élève</th>
          <th v-for="(k, v) in DomainLabels" :key="v" class="pa-2">
            {{ k }}
          </th>
        </tr>
        <tr v-for="(student, i) in advances" :key="i">
          <td>{{ student.Student.Name }} {{ student.Student.Surname }}</td>
          <td v-for="(k, v) in DomainLabels" :key="v" class="pa-2 text-center">
            <RankIcon :rank="student.Advance.Advance[v]"></RankIcon>
          </td>
        </tr>
      </v-table>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import {
  IdClassroom,
  Classroom,
  studentAdvance,
  DomainLabels,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { onMounted } from "vue";
import { ref } from "vue";
import RankIcon from "./RankIcon.vue";

interface Props {
  classrooms: Classroom[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  //   (e: "goTo", d: Domain): void;
}>();

onMounted(fetchAdvance);

const selected = ref<IdClassroom | null>(
  props.classrooms.length ? props.classrooms[0].id : null
);

const advances = ref<studentAdvance[]>([]);

async function fetchAdvance() {
  if (selected.value == null) return;
  const res = await controller.CeinturesGetStudentsAdvance({
    "classroom-id": selected.value,
  });
  if (res === undefined) return;
  advances.value = res || [];
}
</script>
