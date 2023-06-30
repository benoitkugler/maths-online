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
      :class="
        (isDraggedOver ? 'bg-blue-lighten-3' : '') +
        ' fill-height rounded pt-0 mx-1'
      "
      @dragover="onDragOver"
      @dragleave="isDraggedOver = false"
      @drop="onDrop"
    >
      <classroom-travaux
        :classroom="classroom"
        :sheets="props.homeworks.Sheets"
        :classrooms="props.homeworks.Travaux.map((i) => i.Classroom)"
        @update="(tr) => emit('update', tr)"
        @copy="(tr, cl) => emit('copy', tr, cl)"
        @delete="(tr) => emit('delete', tr)"
      ></classroom-travaux>
    </v-window-item>
  </v-window>
</template>

<script setup lang="ts">
import type { HomeworksT } from "@/controller/utils";
import { $ref } from "vue/macros";
import ClassroomTravaux from "./ClassroomTravaux.vue";
import type { SheetExt, Travail } from "@/controller/api_gen";

interface Props {
  homeworks: HomeworksT;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "create", idSheet: number, idClassroom: number): void;
  (e: "delete", travail: Travail): void;
  (e: "copy", travail: Travail, target: number): void;
  (e: "update", travail: Travail): void;
}>();

let tab = $ref(0);

let isDraggedOver = $ref(false);
function onDrop(ev: DragEvent) {
  const eventData: { sheet?: SheetExt } = JSON.parse(
    ev.dataTransfer?.getData("text/json") || ""
  );
  isDraggedOver = false;
  if (!eventData.sheet) return;
  const classroom = props.homeworks.Travaux[tab].Classroom;
  emit("create", eventData.sheet.Sheet.Id, classroom.id);
}

function onDragOver(ev: DragEvent) {
  ev.preventDefault();
  isDraggedOver = true;
}
</script>
