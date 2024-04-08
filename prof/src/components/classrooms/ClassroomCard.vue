<template>
  <v-card
    :title="props.classroom.Classroom.name"
    class="ma-2"
    @click="emit('showStudents')"
  >
    <v-card-text>
      <v-row justify="center" class="mt-1">
        <v-col cols="auto">
          <v-chip link variant="elevated" color="primary">
            {{ eleveText }}
          </v-chip>
        </v-col>
      </v-row>
    </v-card-text>

    <v-card-actions>
      <v-btn icon color="red" title="Supprimer" @click.stop="emit('delete')">
        <v-icon icon="mdi-delete"></v-icon>
      </v-btn>
      <v-spacer></v-spacer>
      <v-btn @click.stop="emit('update')">
        <template v-slot:append>
          <v-icon>mdi-cog</v-icon>
        </template>
        Modifier
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import type { ClassroomExt } from "@/controller/api_gen";
import { computed } from "vue";

interface Props {
  classroom: ClassroomExt;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "showStudents"): void;
  (e: "delete"): void;
  (e: "update"): void;
}>();

const eleveText = computed(() => {
  switch (props.classroom.NbStudents) {
    case 0:
      return "Ajouter des élèves...";
    case 1:
      return "Un élève";
    default:
      return `${props.classroom.NbStudents} élèves`;
  }
});
</script>
