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

          <v-btn icon title="Plus d'options" size="x-small" class="mx-2">
            <v-icon>mdi-dots-vertical</v-icon>
            <v-menu activator="parent">
              <v-list density="compact">
                <!-- self access -->
                <v-list-item
                  title="Paramétrer l'accès libre..."
                  prepend-icon="mdi-account-multiple"
                  @click="emit('show-selfaccess')"
                ></v-list-item>
                <!-- duplicate -->
                <v-list-item
                  v-if="config.Origin.Visibility == Visibility.Admin"
                  title="Dupliquer et importer"
                  prepend-icon="mdi-content-copy"
                  @click="emit('duplicate')"
                ></v-list-item>
                <template v-else>
                  <v-divider thickness="1"></v-divider>
                  <OriginButton
                    as-list-item
                    :origin="config.Origin"
                    @update-public="(b) => emit('update-public', b)"
                    @create-review="emit('create-review')"
                  ></OriginButton>
                  <v-divider thickness="1"></v-divider>
                  <v-list-item
                    v-if="isPersonnal(config)"
                    title=" Modifier..."
                    prepend-icon="mdi-pencil"
                    @click="emit('edit')"
                  >
                  </v-list-item>
                  <v-list-item
                    v-if="isPersonnal(config)"
                    title="Supprimer"
                    prepend-icon="mdi-delete"
                    @click="emit('delete')"
                  >
                  </v-list-item>
                </template>
              </v-list>
            </v-menu>
          </v-btn>
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
