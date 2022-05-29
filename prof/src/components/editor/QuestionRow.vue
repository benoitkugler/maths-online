<template>
  <v-list-item
    dense
    rounded
    :class="'py-1 my-2 ' + colorClass"
    @click="emit('clicked', props.question)"
  >
    <v-row no-gutters>
      <v-col cols="auto" align-self="center">
        <v-btn
          v-if="isPersonnal"
          size="x-small"
          icon
          @click.stop="emit('delete', props.question)"
          title="Supprimer"
        >
          <v-icon icon="mdi-delete" color="red" size="small"></v-icon>
        </v-btn>

        <OriginButton
          :origin="props.question.Origin"
          @update-public="(b) => emit('updatePublic', props.question.Id, b)"
        ></OriginButton>

        <v-btn
          v-if="isPersonnal && !props.question.IsInGroup"
          class="mx-1"
          size="x-small"
          icon
          @click.stop="emit('duplicate', props.question)"
          title="Dupliquer en différenciant la difficulté"
        >
          <v-icon
            icon="mdi-content-copy"
            color="info"
            size="small"
          ></v-icon> </v-btn
      ></v-col>
      <v-col align-self="center">
        <div class="ml-2">
          <i>
            <small> ({{ props.question.Id }}) </small>
          </i>
          {{ props.question.Title ? props.question.Title : "..." }}
        </div>
      </v-col>
      <v-col style="text-align: right" align-self="center">
        <TagChip :tag="tag" :key="tag" v-for="tag in question.Tags"></TagChip>
      </v-col>
    </v-row>
  </v-list-item>
</template>

<script setup lang="ts">
import { Visibility, type QuestionHeader } from "@/controller/api_gen";
import { visiblityColors } from "@/controller/editor";
import { computed } from "vue";
import OriginButton from "../OriginButton.vue";
import TagChip from "./utils/TagChip.vue";

interface Props {
  question: QuestionHeader;
}
const props = defineProps<Props>();
const emit = defineEmits<{
  (e: "delete", question: QuestionHeader): void;
  (e: "clicked", question: QuestionHeader): void;
  (e: "duplicate", question: QuestionHeader): void;
  (e: "updatePublic", questionID: number, isPublic: boolean): void;
}>();

const colorClass = computed(
  () => "bg-" + visiblityColors[props.question.Origin.Visibility]
);

const isPersonnal = props.question.Origin.Visibility == Visibility.Personnal;
</script>
