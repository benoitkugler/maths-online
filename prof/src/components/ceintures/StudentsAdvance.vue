<template>
  <v-card title="Progression des élèves">
    <template v-slot:append>
      <v-row>
        <v-col>
          <v-btn-toggle variant="outlined" v-model="mode">
            <v-btn value="couleur">Couleur</v-btn>
            <v-btn value="note">Note</v-btn>
          </v-btn-toggle>
        </v-col>
        <v-divider vertical thickness="1"></v-divider>
        <v-col>
          <v-select
            variant="outlined"
            density="comfortable"
            hide-details
            label="Classe"
            :items="
              props.classrooms.map((c) => ({ title: c.name, value: c.id }))
            "
            v-model="selected"
            no-data-text="Vous n'avez aucune classe."
            @update:model-value="fetchAdvance"
          ></v-select>
        </v-col>
      </v-row>
    </template>
    <v-card-text>
      <v-table>
        <tr>
          <th width="150px">Élève</th>
          <th v-for="(k, v) in DomainLabels" :key="v" class="pa-2">
            {{ k }}
          </th>
        </tr>
        <tr
          v-for="(student, i) in advances"
          :key="i"
          :class="{ 'bg-grey-lighten-4': i % 2 == 0 }"
        >
          <td class="px-2">
            {{ student.Student.Name }} {{ student.Student.Surname }}
          </td>
          <td v-for="(k, v) in DomainLabels" :key="v" class="pa-2 text-center">
            <RankIcon
              v-if="mode == 'couleur'"
              :rank="student.Advance.Advance[v]"
            ></RankIcon>
            <span v-else>{{ (student.Advance.Advance[v] || 0) * 2 }}</span>
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

const mode = ref<"couleur" | "note">("couleur");
</script>
