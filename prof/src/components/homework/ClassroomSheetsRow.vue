<template>
  <v-card class="mx-2 bg-grey-lighten-4">
    <v-row no-gutters>
      <v-col>
        <v-card-title>{{ classroom.Classroom.name }}</v-card-title>
      </v-col>
      <v-col align-self="center" cols="auto">
        <v-btn
          class="mx-2"
          @click="emit('add')"
          title="Créer une nouvelle fiche"
        >
          <v-icon icon="mdi-plus" color="success"></v-icon>
          Ajouter
        </v-btn>
      </v-col>
    </v-row>
    <v-card-text>
      <v-row>
        <v-col
          v-if="!classroom.Sheets?.length"
          align-self="center"
          style="text-align: center"
        >
          <i>Aucune fiche n'est encore définie.</i>
        </v-col>

        <v-col
          lg="3"
          md="4"
          sm="6"
          xs="12"
          v-for="(sheet, index) in classroom.Sheets"
          :key="index"
        >
          <sheet-card
            :sheet="sheet"
            :classrooms="props.classrooms"
            @delete="emit('delete', sheet)"
            @copy="(idClassroom) => emit('copy', sheet.Sheet.Id, idClassroom)"
            @update="emit('update', sheet)"
          ></sheet-card>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type {
  Classroom,
  ClassroomSheets,
  SheetExt,
} from "@/controller/api_gen";
import SheetCard from "./SheetCard.vue";

interface Props {
  classroom: ClassroomSheets;
  classrooms: Classroom[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "add"): void;
  (e: "update", sheet: SheetExt): void;
  (e: "delete", sheet: SheetExt): void;
  (e: "copy", idSheet: number, idClassroom: number): void;
}>();
</script>
