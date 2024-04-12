<template>
  <v-row v-if="!props.homeworks.Travaux.length" justify="center">
    <v-col style="text-align: center" class="my-3">
      <i>Aucune classe n'est enregistrée.</i>
    </v-col>
  </v-row>

  <v-tabs v-model="tab" align-tabs="center">
    <v-tab v-for="(classroom, index) in props.homeworks.Travaux" :key="index">
      {{ classroom.Classroom.name }}
    </v-tab>
  </v-tabs>
  <v-window v-model="tab" class="mt-1" style="height: 95%">
    <v-window-item
      v-for="(classroom, index) in props.homeworks.Travaux"
      :key="index"
      :value="index"
      class="fill-height rounded pt-0 mx-1"
    >
      <classroom-travaux
        :classroom="classroom"
        :sheets="props.homeworks.Sheets"
        :classrooms="props.homeworks.Travaux.map((i) => i.Classroom)"
        @create="emit('create', classroom.Classroom.id)"
        @set-favorite="(s) => emit('setFavorite', s)"
        @edit-sheet="(s) => emit('editSheet', s)"
        @update="(tr) => emit('update', tr)"
        @copy="(tr, cl) => emit('copy', tr, cl)"
        @delete="(tr) => emit('delete', tr)"
      ></classroom-travaux>
    </v-window-item>
  </v-window>
</template>

<script setup lang="ts">
import type { HomeworksT } from "@/controller/utils";
import { ref } from "vue";
import ClassroomTravaux from "./ClassroomTravaux.vue";
import type { Int, Sheet, SheetExt, Travail } from "@/controller/api_gen";

interface Props {
  homeworks: HomeworksT;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "create", idClassroom: Int): void;
  (e: "createWith", idSheet: Int, idClassroom: Int): void;
  (e: "setFavorite", sheet: Sheet): void;
  (e: "delete", travail: Travail): void;
  (e: "copy", travail: Travail, target: Int): void;
  (e: "update", travail: Travail): void;
  (e: "editSheet", sheet: SheetExt): void;
}>();

const tab = ref(0);

defineExpose({ showClassroom });

function showClassroom(id: number) {
  const index = props.homeworks.Travaux.findIndex(
    (cl) => cl.Classroom.id == id
  );
  if (index != -1) tab.value = index;
}
</script>
