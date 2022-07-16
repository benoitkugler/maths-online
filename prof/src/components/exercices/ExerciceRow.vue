<template>
  <v-list-item
    dense
    rounded
    :class="'py-1 my-2 ' + colorClass"
    @click="emit('clicked', props.exercice)"
  >
    <v-row no-gutters>
      <v-col cols="auto" align-self="center">
        <v-btn
          v-if="isPersonnal"
          size="x-small"
          icon
          @click.stop="emit('delete', props.exercice)"
          title="Supprimer"
        >
          <v-icon icon="mdi-delete" color="red" size="small"></v-icon>
        </v-btn>

        <OriginButton
          :origin="props.exercice.Origin"
          @update-public="
            (b) => emit('updatePublic', props.exercice.Exercice.Id, b)
          "
        ></OriginButton>
      </v-col>
      <v-col align-self="center">
        <div class="ml-2">
          <i>
            <small> ({{ props.exercice.Exercice.Id }}) </small>
          </i>
          {{
            props.exercice.Exercice.Title
              ? props.exercice.Exercice.Title
              : "..."
          }}
        </div>
      </v-col>
      <v-col cols="4">
        <v-chip>
          {{ props.exercice.Questions?.length || 0 }} question(s)</v-chip
        >
      </v-col>
    </v-row>
  </v-list-item>
</template>

<script setup lang="ts">
import { Visibility, type ExerciceHeader } from "@/controller/api_gen";
import { visiblityColors } from "@/controller/editor";
import { computed } from "vue";
import OriginButton from "../OriginButton.vue";

interface Props {
  exercice: ExerciceHeader;
}
const props = defineProps<Props>();
const emit = defineEmits<{
  (e: "delete", question: ExerciceHeader): void;
  (e: "clicked", question: ExerciceHeader): void;
  (e: "duplicate", question: ExerciceHeader): void;
  (e: "updatePublic", exerciceID: number, isPublic: boolean): void;
}>();

const colorClass = computed(
  () => "bg-" + visiblityColors[props.exercice.Origin.Visibility]
);

const isPersonnal = props.exercice.Origin.Visibility == Visibility.Personnal;
</script>
