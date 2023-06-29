<template>
  <v-chip v-if="task.IdWork.Kind == WorkKind.WorkExercice">
    / {{ taskBareme(task) }}
  </v-chip>

  <v-menu
    v-else-if="task.IdWork.Kind == WorkKind.WorkMonoquestion"
    offset-y
    :close-on-content-click="false"
    :model-value="monoquestionToEdit != null"
    @update:model-value="monoquestionToEdit = null"
  >
    <template v-slot:activator="{ isActive, props }">
      <v-chip
        elevation="2"
        v-on="{ isActive }"
        v-bind="props"
        @click="startEditMonoquestion()"
      >
        / {{ taskBareme(task) }}
      </v-chip>
    </template>
    <MonoquestionDetails
      v-if="monoquestionToEdit != null"
      :monoquestion="monoquestionToEdit!"
      @update="
        (qu) => {
          monoquestionToEdit = null;
          emit('updateMonoquestion', qu);
        }
      "
    ></MonoquestionDetails>
  </v-menu>

  <v-menu
    v-else-if="props.task.IdWork.Kind == WorkKind.WorkRandomMonoquestion"
    offset-y
    :close-on-content-click="false"
    :model-value="randomMonoquestionToEdit != null"
    @update:model-value="randomMonoquestionToEdit = null"
  >
    <template v-slot:activator="{ isActive, props }">
      <v-chip
        elevation="2"
        v-on="{ isActive }"
        v-bind="props"
        @click="startEditRandomMonoquestion()"
      >
        / {{ taskBareme(task) }}
      </v-chip>
    </template>
    <RandomMonoquestionDetails
      v-if="randomMonoquestionToEdit != null"
      :random-monoquestion="randomMonoquestionToEdit!"
      @update="
        (qu) => {
          emit('updateRandomMonoquestion', qu);
          randomMonoquestionToEdit = null;
        }
      "
    ></RandomMonoquestionDetails>
  </v-menu>
</template>

<script setup lang="ts">
import {
  WorkKind,
  type TaskExt,
  type Monoquestion,
  type RandomMonoquestion,
} from "@/controller/api_gen";
import { taskBareme } from "@/controller/utils";
import MonoquestionDetails from "./MonoquestionDetails.vue";
import { $ref } from "vue/macros";
import { controller } from "@/controller/controller";
import RandomMonoquestionDetails from "./RandomMonoquestionDetails.vue";

interface Props {
  task: TaskExt;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "updateMonoquestion", v: Monoquestion): void;
  (e: "updateRandomMonoquestion", v: RandomMonoquestion): void;
}>();

let monoquestionToEdit = $ref<Monoquestion | null>(null);
let randomMonoquestionToEdit = $ref<RandomMonoquestion | null>(null);

// load the monoquestion details and show an editor
async function startEditMonoquestion() {
  const res = await controller.HomeworkGetMonoquestion({
    "id-monoquestion": props.task.IdWork.ID,
  });
  if (res == undefined) return;

  monoquestionToEdit = res;
}
// load the monoquestion details and show an editor
async function startEditRandomMonoquestion() {
  const res = await controller.HomeworkGetRandomMonoquestion({
    "id-randommonoquestion": props.task.IdWork.ID,
  });
  if (res == undefined) return;

  randomMonoquestionToEdit = res;
}
</script>

<style scoped></style>
