<template>
  <v-row>
    <v-col cols="5" md="3" align-self="center">
      <OriginButton
        :origin="config.Origin"
        @update-public="(b) => emit('update-public', b)"
      ></OriginButton>
      <v-btn
        class="mx-2 my-1"
        size="x-small"
        icon
        @click="emit('duplicate')"
        title="Dupliquer cette session"
      >
        <v-icon icon="mdi-content-copy" color="secondary"></v-icon>
      </v-btn>

      <v-btn
        icon
        size="x-small"
        title="Editer"
        class="mx-2"
        @click="emit('edit')"
        v-if="isPersonnal(config)"
      >
        <v-icon icon="mdi-pencil"></v-icon>
      </v-btn>

      <v-btn
        v-if="isPersonnal(config)"
        icon
        size="x-small"
        title="Lancer"
        class="mx-2"
        @click="emit('launch')"
        :disabled="props.disableLaunch"
      >
        <v-icon icon="mdi-play" color="green"></v-icon>
      </v-btn>

      <v-btn
        v-if="isPersonnal(config)"
        class="mx-2"
        size="x-small"
        icon
        @click="emit('delete')"
        title="Supprimer cette session"
      >
        <v-icon icon="mdi-delete" color="red"></v-icon>
      </v-btn>
    </v-col>
    <v-col cols="4" md="3" align-self="center" style="text-align: center">
      {{ config.Config.Name }}
      <small class="text-grey">
        {{ formatCategories(config.Config) }}
      </small>
    </v-col>
    <v-col cols="3" md="2" align-self="center" style="text-align: center">
      <small class="text-primary">
        {{ formatDifficulties(config.Config) }}
      </small>
    </v-col>
    <v-col align-self="center" md="4">
      <v-card class="bg-grey-lighten-2">
        <v-card-text class="px-0 py-1">
          <v-row justify="center">
            <v-col
              cols="6"
              align-self="center"
              style="text-align: center"
              v-if="config.Config.Questions.Tags.every((v) => !v)"
            >
              <i>Aucune question configurée.</i>
            </v-col>
            <v-col
              class="my-1 px-0"
              cols="2"
              align-self="center"
              style="text-align: center"
              v-for="(categorie, index) in config.Config.Questions.Tags || []"
              :key="index"
              v-show="categorie && categorie.length != 0"
            >
              <v-chip :color="colors[index]" variant="outlined">
                {{ config.NbQuestionsByCategories[index] }}
              </v-chip>
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import {
  Visibility,
  type Trivial,
  type TrivialExt,
} from "@/controller/api_gen";
import { commonTags } from "@/controller/editor";
import { colorsPerCategorie } from "@/controller/trivial";
import OriginButton from "../OriginButton.vue";

interface Props {
  config: TrivialExt;
  disableLaunch: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update-public", isPublic: boolean): void;
  (e: "duplicate"): void;
  (e: "edit"): void;
  (e: "launch"): void;
  (e: "delete"): void;
}>();

const colors = colorsPerCategorie;

function isPersonnal(config: TrivialExt) {
  return config.Origin.Visibility == Visibility.Personnal;
}

function formatCategories(config: Trivial) {
  const allUnions: string[][] = [];
  config.Questions.Tags.forEach((cat) =>
    allUnions.push(...(cat || []).map((s) => s || []))
  );
  const common = commonTags(allUnions);
  if (common.length != 0) {
    return "(" + common.join(", ") + ")";
  }
  return "";
}

function formatDifficulties(config: Trivial) {
  const l = config.Questions.Difficulties || [];
  if (l.length) {
    return l.join(", ");
  }
  return "Toutes difficultés";
}
</script>

<style scoped></style>
