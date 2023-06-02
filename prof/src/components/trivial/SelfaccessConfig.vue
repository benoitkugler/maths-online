<template>
  <v-card
    title="Accès en autonomie"
    subtitle="Les élèves des classes sélectionnées peuvent lancer une partie de ce TrivMaths en autonomie."
  >
    <v-card-text>
      <v-list
        v-if="data != null"
        select-strategy="classic"
        v-model:selected="selected"
        class="my-1 overflow-y-auto"
        style="max-height: 55vh"
      >
        <v-list-item
          v-for="(classroom, index) in data.Classrooms"
          :key="index"
          :title="classroom.name"
          :value="classroom.id"
          rounded
          class="my-1"
        >
          <template v-slot:append="{ isActive }">
            <v-list-item-action end>
              <v-checkbox-btn :model-value="isActive"></v-checkbox-btn>
            </v-list-item-action>
          </template>
        </v-list-item>
      </v-list>
    </v-card-text>
    <v-card-actions>
      <v-spacer> </v-spacer>
      <v-btn color="success" @click="save">Enregistrer</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import type { Trivial, TrivialSelfaccess } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { onMounted } from "vue";
import { onActivated } from "vue";
import { $ref } from "vue/macros";

interface Props {
  config: Trivial;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

onMounted(() => fetch());
onActivated(() => fetch());

let selected = $ref<number[]>([]);
let data = $ref<TrivialSelfaccess | null>(null);

async function fetch() {
  const res = await controller.TrivialGetSelfaccess({
    "id-trivial": props.config.Id,
  });
  if (res == undefined) return;
  data = res;
  selected = res.Actives || [];
}

async function save() {
  await controller.TrivialUpdateSelfaccess({
    IdTrivial: props.config.Id,
    IdClassrooms: selected,
  });
  emit("close");
}
</script>

<style scoped></style>
