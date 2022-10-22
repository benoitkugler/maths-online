<template>
  <v-dialog
    :model-value="editedConfig != null"
    @update:model-value="editedConfig = null"
    :retain-focus="false"
  >
    <edit-config
      v-if="editedConfig != null"
      :edited="editedConfig"
      :all-known-tags="allKnownTags"
      @closed="editedConfig = null"
      @update="updateConfig"
    >
    </edit-config>
  </v-dialog>

  <v-dialog
    :model-value="launchingConfig != null"
    @update:model-value="launchingConfig = null"
    max-width="870px"
  >
    <launch-options @launch="launchSession"></launch-options>
  </v-dialog>

  <v-dialog
    fullscreen
    :model-value="showMonitor"
    @update:model-value="showMonitorChanged"
    :retain-focus="false"
  >
    <session-monitor @closed="closeMonitor"></session-monitor>
  </v-dialog>

  <v-card class="my-5 mx-auto" width="80%">
    <v-row>
      <v-col>
        <v-card-title>Triv'Maths</v-card-title>
        <v-card-subtitle
          >Configurer et lancer une partie de Triv'Maths</v-card-subtitle
        >
      </v-col>

      <v-col align-self="center" style="text-align: right" cols="4">
        <v-btn
          class="mx-2"
          @click="createConfig"
          title="Créer une nouvelle session"
        >
          <v-icon icon="mdi-plus" color="success"></v-icon>
          Créer
        </v-btn>
      </v-col>
    </v-row>

    <v-alert color="secondary" v-if="sessionMeta.NbGames > 0" class="my-2 mx-4">
      <v-row>
        <v-col>
          Parties en cours : <v-chip>{{ sessionMeta.NbGames }}</v-chip>
        </v-col>
        <v-spacer></v-spacer>
        <v-col>
          <v-btn @click="showMonitor = true"> Suivre les parties </v-btn>
        </v-col>
      </v-row>
    </v-alert>

    <v-list>
      <v-list-item
        v-for="config in configs"
        :key="config.Config.Id"
        class="my-3"
      >
        <v-row>
          <v-col cols="3" align-self="center">
            <origin-button
              :origin="config.Origin"
              @update-public="(b) => updatePublic(config.Config, b)"
            ></origin-button>
            <v-btn
              class="mx-2 my-1"
              size="x-small"
              icon
              @click="duplicateConfig(config.Config)"
              title="Dupliquer cette session"
            >
              <v-icon icon="mdi-content-copy" color="secondary"></v-icon>
            </v-btn>

            <v-btn
              icon
              size="x-small"
              title="Editer"
              class="mx-2"
              @click="editedConfig = config.Config"
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
              @click="launchingConfig = config.Config"
              :disabled="
                isLaunching ||
                !config.NbQuestionsByCategories.every((v) => v > 0)
              "
            >
              <v-icon icon="mdi-play" color="green"></v-icon>
            </v-btn>

            <v-btn
              v-if="isPersonnal(config)"
              class="mx-2"
              size="x-small"
              icon
              @click="deleteConfig(config.Config)"
              title="Supprimer cette session"
            >
              <v-icon icon="mdi-delete" color="red"></v-icon>
            </v-btn>
          </v-col>
          <v-col cols="3" align-self="center" style="text-align: center">
            {{ config.Config.Name }}
            <small class="text-grey">
              {{ formatCategories(config.Config) }}
            </small>
          </v-col>
          <v-col cols="2" align-self="center" style="text-align: center">
            <small class="text-primary">
              {{ formatDifficulties(config.Config) }}
            </small>
          </v-col>
          <v-col align-self="center" cols="4">
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
                    class="my-1"
                    cols="2"
                    align-self="center"
                    v-for="(categorie, index) in config.Config.Questions.Tags ||
                    []"
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
    </v-list>
  </v-card>
</template>

<script setup lang="ts">
import {
  Visibility,
  type RunningSessionMetaOut,
  type Trivial,
  type TrivialExt,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { commonTags } from "@/controller/editor";
import { colorsPerCategorie } from "@/controller/trivial";
import { computed, onMounted } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import OriginButton from "../components/OriginButton.vue";
import EditConfig from "../components/trivial/EditConfig.vue";
import LaunchOptions from "../components/trivial/LaunchOptions.vue";
import SessionMonitor from "../components/trivial/SessionMonitor.vue";

let allKnownTags = $ref<string[]>([]);

let editedConfig = $ref<Trivial | null>(null);

let _configs = $ref<TrivialExt[]>([]);

const configs = computed(() => {
  const a = _configs.map((v) => v);
  a.sort((u, v) => u.Config.Id - v.Config.Id);
  return a;
});

let isLaunching = $ref(false);

const colors = colorsPerCategorie;

onMounted(async () => {
  fetchSessionMeta();

  const res = await controller.GetTrivialPoursuit();

  if (res === undefined) {
    return;
  }
  _configs = Object.values(res || {});

  const tags = await controller.EditorGetTags();
  allKnownTags = tags || [];
});

function isPersonnal(config: TrivialExt) {
  return config.Origin.Visibility == Visibility.Personnal;
}

async function createConfig() {
  const res = await controller.CreateTrivialPoursuit();
  if (res === undefined) {
    return;
  }
  _configs.push(res);
}

async function updateConfig(config: Trivial) {
  // remove empty categories
  config.Questions.Tags = config.Questions.Tags.map((q) =>
    (q || []).filter((v) => v && v.length != 0)
  );
  const res = await controller.UpdateTrivialPoursuit(config);
  if (res === undefined) {
    return;
  }
  const index = _configs.findIndex((v) => v.Config.Id == config.Id);
  _configs[index] = res;
  editedConfig = null;
}

async function updatePublic(config: Trivial, isPublic: boolean) {
  const res = await controller.UpdateTrivialVisiblity({
    ConfigID: config.Id,
    Public: isPublic,
  });
  if (res === undefined) {
    return;
  }
  const index = _configs.findIndex((v) => v.Config.Id == config.Id);
  _configs[index].Origin.IsPublic = isPublic;
}

async function duplicateConfig(config: Trivial) {
  const res = await controller.DuplicateTrivialPoursuit({ id: config.Id });
  if (res === undefined) {
    return;
  }
  console.log(config, res);

  _configs.push(res);
}

async function deleteConfig(config: Trivial) {
  await controller.DeleteTrivialPoursuit({ id: config.Id });
  _configs = _configs.filter((c) => c.Config.Id != config.Id);
}

let launchingConfig = $ref<Trivial | null>(null);
async function launchSession(groups: number[]) {
  if (launchingConfig == null) {
    return;
  }
  const configID = launchingConfig.Id;
  isLaunching = true;
  const res = await controller.LaunchSessionTrivialPoursuit({
    IdConfig: configID,
    Groups: groups,
  });
  launchingConfig = null;
  isLaunching = false;
  if (res === undefined) {
    return;
  }
  fetchSessionMeta();

  // automatically jump to monitor screen
  showMonitor = true;
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

let sessionMeta = $ref<RunningSessionMetaOut>({ NbGames: 0 });
async function fetchSessionMeta() {
  const res = await controller.GetTrivialRunningSessions();
  if (res == undefined) {
    return;
  }
  sessionMeta = res;
}

let showMonitor = $ref(false);
function showMonitorChanged(show: boolean) {
  showMonitor = show;
  if (!show) {
    closeMonitor();
  }
}

function closeMonitor() {
  showMonitor = false;
  fetchSessionMeta();
}
</script>

<style>
:deep(.v-dialog .v-overlay__content) {
  max-width: 900px;
}
</style>
