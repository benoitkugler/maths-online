<template>
  <v-col cols="12" md="6" class="my-1">
    <v-list-item
      :class="{ [colorClass]: true, 'px-0': true, 'mx-1': true }"
      rounded
    >
      <v-row no-gutters>
        <v-col align-self="center">
          <v-menu width="600px">
            <template v-slot:activator="{ isActive, props: slotProps }">
              <v-card
                density="compact"
                v-on="{ isActive }"
                v-bind="slotProps"
                variant="tonal"
                class="ml-2 py-2"
              >
                <h4 class="ml-4">
                  {{ config.Config.Name || "Sans titre" }}
                </h4>
                <v-card-subtitle>
                  {{ formatDifficulties(config.Config) }}
                </v-card-subtitle>
              </v-card>
            </template>

            <QuestionsRecap :config="props.config"></QuestionsRecap>
          </v-menu>
        </v-col>

        <v-col cols="auto" align-self="center" class="my-3">
          <v-btn
            class="ml-2"
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

          <v-btn
            v-if="config.Origin.Visibility == Visibility.Admin"
            class="ml-3 mr-2 my-1"
            size="x-small"
            icon
            @click="emit('duplicate')"
            title="Dupliquer et importer"
          >
            <v-icon icon="mdi-content-copy" color="secondary"></v-icon>
          </v-btn>
          <v-menu offset-y close-on-content-click v-else>
            <template v-slot:activator="{ isActive, props }">
              <v-btn
                icon
                title="Plus d'options"
                v-on="{ isActive }"
                v-bind="props"
                size="x-small"
                class="ml-3 mr-2 my-1"
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
              <v-list-item>
                <OriginButton
                  :origin="config.Origin"
                  @update-public="(b) => emit('update-public', b)"
                  @create-review="emit('create-review')"
                ></OriginButton>
              </v-list-item>
            </v-list>
          </v-menu>
        </v-col>
      </v-row>
    </v-list-item>
  </v-col>
</template>

<script setup lang="ts">
import {
  Visibility,
  type Trivial,
  type TrivialExt,
} from "@/controller/api_gen";
import { computed } from "vue";
import OriginButton from "../OriginButton.vue";
import QuestionsRecap from "./QuestionsRecap.vue";
import { colorForOrigin } from "@/controller/utils";

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

const colorClass = computed(() => "bg-" + colorForOrigin(props.config.Origin));

function isPersonnal(config: TrivialExt) {
  return config.Origin.Visibility == Visibility.Personnal;
}

function formatDifficulties(config: Trivial) {
  const l = config.Questions.Difficulties || [];
  if (l.length) {
    return l.join(", ");
  }
  return "Toutes difficultés";
}
</script>

<style></style>
