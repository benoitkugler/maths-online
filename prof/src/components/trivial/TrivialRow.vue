<template>
  <v-list-item
    :class="{ 'my-3': true, 'mx-4': true, [colorClass]: true }"
    rounded
  >
    <v-row>
      <v-col cols="6" sm="4" md="3" align-self="center" class="pr-0 pl-1 my-3">
        <v-btn
          v-if="config.Origin.Visibility == Visibility.Admin"
          class="ml-3 mr-1 my-1"
          size="x-small"
          icon
          @click="emit('duplicate')"
          title="Dupliquer cette session"
        >
          <v-icon icon="mdi-content-copy" color="secondary"></v-icon>
        </v-btn>
        <template v-else>
          <v-menu offset-y close-on-content-click>
            <template v-slot:activator="{ isActive, props }">
              <v-btn
                icon
                title="Plus d'options"
                v-on="{ isActive }"
                v-bind="props"
                size="x-small"
                class="ml-3"
              >
                <v-icon icon="mdi-dots-vertical"></v-icon>
              </v-btn>
            </template>
            <v-list>
              <v-list-item>
                <v-btn
                  flat
                  size="small"
                  title="Editer"
                  @click="emit('edit')"
                  v-if="isPersonnal(config)"
                >
                  <v-icon icon="mdi-pencil" class="mr-2"></v-icon>
                  Modifier...
                </v-btn>
              </v-list-item>
              <v-list-item>
                <v-btn
                  v-if="isPersonnal(config)"
                  flat
                  size="small"
                  @click="emit('delete')"
                  title="Supprimer cette session"
                >
                  <v-icon icon="mdi-delete" color="red" class="mr-2"></v-icon>
                  Supprimer
                </v-btn></v-list-item
              >
            </v-list>
          </v-menu>

          <OriginButton
            :origin="config.Origin"
            @update-public="(b) => emit('update-public', b)"
            @create-review="emit('create-review')"
          ></OriginButton>
        </template>
        <v-btn
          title="Lancer"
          size="small"
          @click="emit('launch')"
          :disabled="props.disableLaunch"
        >
          <v-icon icon="mdi-play" color="green"></v-icon>
          Lancer
        </v-btn>

        <v-tooltip text="Paramétrer l'accès libre...">
          <template v-slot:activator="{ props }">
            <v-btn
              v-bind="props"
              icon
              size="x-small"
              class="ml-2"
              @click="emit('show-selfaccess')"
            >
              <v-icon icon="mdi-account-multiple"></v-icon>
            </v-btn>
          </template>
        </v-tooltip>
      </v-col>

      <v-col
        cols="6"
        sm="5"
        md="3"
        class="px-1"
        align-self="center"
        style="text-align: center"
      >
        {{ config.Config.Name }}
        <small class="text-grey">
          {{ formatCategories(config.Config) }}
        </small>
      </v-col>

      <v-col
        class="d-none d-sm-block"
        cols="3"
        md="2"
        align-self="center"
        style="text-align: center"
      >
        <small class="text-primary">
          {{ formatDifficulties(config.Config) }}
        </small>
      </v-col>

      <v-col align-self="center" md="4">
        <v-card class="bg-grey-lighten-4">
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
  </v-list-item>
</template>

<script setup lang="ts">
import {
  Visibility,
  type Trivial,
  type TrivialExt,
} from "@/controller/api_gen";
import { visiblityColors } from "@/controller/editor";
import { colorsPerCategorie } from "@/controller/trivial";
import { computed } from "vue";
import OriginButton from "../OriginButton.vue";

interface Props {
  config: TrivialExt;
  disableLaunch: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update-public", isPublic: boolean): void;
  (e: "create-review"): void;
  (e: "duplicate"): void;
  (e: "edit"): void;
  (e: "launch"): void;
  (e: "delete"): void;
  (e: "show-selfaccess"): void;
}>();

const colors = colorsPerCategorie;

const colorClass = computed(() =>
  props.config.Origin.Visibility == Visibility.Admin
    ? "bg-" + visiblityColors[Visibility.Admin]
    : ""
);

function isPersonnal(config: TrivialExt) {
  return config.Origin.Visibility == Visibility.Personnal;
}

/** return the list of tags shared by all the list */
function commonTags(tags: string[][]) {
  const crible: { [key: string]: number } = {};
  tags.forEach((l) =>
    l.forEach((tag) => (crible[tag] = (crible[tag] || 0) + 1))
  );
  return Object.entries(crible)
    .filter((entry) => entry[1] == tags.length)
    .map((entry) => entry[0]);
}

function formatCategories(config: Trivial) {
  const allUnions: string[][] = [];
  config.Questions.Tags.forEach((cat) =>
    allUnions.push(...(cat || []).map((s) => (s || []).map((ts) => ts.Tag)))
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
