<template>
  <v-dialog
    :model-value="editedConfig != null"
    @update:model-value="editedConfig = null"
  >
    <v-card
      v-if="editedConfig != null"
      max-height="80vh"
      width="800"
      class="overflow-y-auto py-2"
    >
      <v-row>
        <v-col>
          <v-card-title>Modifier la session de TrivialPoursuit</v-card-title>
        </v-col>
        <v-col style="text-align: right">
          <v-btn icon flat class="mx-2" @click="editedConfig = null">
            <v-icon icon="mdi-close"></v-icon>
          </v-btn>
        </v-col>
      </v-row>
      <v-card-text>
        <v-list>
          <v-list-subheader>
            <h3>Choix des questions</h3>
            <small
              >Chaque catégorie est définie par une <i>union</i> d'<i
                >intersections</i
              >
              d'étiquettes.</small
            >
          </v-list-subheader>
          <v-list-item
            v-for="(categorie, index) in editedConfig?.Questions"
            rounded
            :style="{
              'border-color': colors[index],
              borderWidth: '2px',
              borderStyle: 'solid'
            }"
            class="my-2"
          >
            <v-list-item-subtitle
              >Catégorie {{ index + 1 }}</v-list-item-subtitle
            >
            <tags-selector
              :all-tags="allKnownTags"
              :model-value="categorie || []"
              @update:model-value="v => editedConfig!.Questions[index] = v"
            ></tags-selector>
          </v-list-item>
        </v-list>

        <v-row class="mt-2">
          <v-col cols="6">
            <v-text-field
              density="compact"
              variant="outlined"
              label="Durée limite pour une question"
              type="number"
              min="1"
              suffix="sec"
              v-model.number="editedConfig!.QuestionTimeout"
            ></v-text-field>
          </v-col>
        </v-row>
      </v-card-text>

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="success" @click="updateConfig">
          Enregistrer les modifications
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog
    :model-value="launchingConfig != null"
    @update:model-value="launchingConfig = null"
  >
    <launch-options @launch="launchSession"></launch-options>
  </v-dialog>

  <v-dialog
    fullscreen
    :model-value="sessionToMonitor != null"
    @update:model-value="sessionToMonitor = null"
  >
    <session-monitor
      v-if="sessionToMonitor != null"
      :running-session="sessionToMonitor!.Running "
      @close="sessionToMonitor = null"
    ></session-monitor>
  </v-dialog>

  <v-card class="my-3 mx-auto" width="60%">
    <v-row>
      <v-col>
        <v-card-title>Trivial Poursuit</v-card-title>
        <v-card-subtitle
          >Configurer et lancer une partie de Trivial Poursuit</v-card-subtitle
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

    <v-list>
      <v-list-item v-for="config in configs" class="my-2">
        <v-row>
          <v-col cols="auto" align-self="center">
            <v-btn
              icon
              size="x-small"
              title="Editer"
              class="mx-2"
              @click="editedConfig = config.Config"
              :disabled="config.Running.SessionID != ''"
            >
              <v-icon icon="mdi-pencil"></v-icon>
            </v-btn>
            <v-btn
              v-if="config.Running.SessionID != ''"
              @click="monitor(config)"
            >
              Suivre
            </v-btn>
            <v-btn
              v-else
              icon
              size="x-small"
              title="Lancer"
              class="mx-2"
              @click="launchingConfig = config.Config"
              :disabled="
                isLaunching || !config.NbQuestionsByCategories.every(v => v > 0)
              "
            >
              <v-icon icon="mdi-play" color="green"></v-icon>
            </v-btn>

            <v-btn
              v-if="config.Running.SessionID == ''"
              class="mx-2"
              size="x-small"
              icon
              @click="deleteConfig(config.Config)"
              title="Supprimer cette session"
            >
              <v-icon icon="mdi-delete" color="red"></v-icon>
            </v-btn>
            <v-btn
              v-else
              class="mx-2"
              size="x-small"
              icon
              @click="stopSession(config.Config)"
              title="Stopper la session"
            >
              <v-icon icon="mdi-close" color="red"></v-icon>
            </v-btn>
          </v-col>
          <v-col align-self="center" md="8">
            <v-row justify="center" class="bg-grey-lighten-1 rounded">
              <v-col
                cols="6"
                align-self="center"
                style="text-align: center"
                v-if="config.Config.Questions!.every(v=>!v)"
              >
                <i>Aucune question configurée.</i>
              </v-col>
              <v-col
                cols="2"
                align-self="center"
                v-for="(categorie, index) in config.Config.Questions || []"
                v-show="categorie && categorie.length != 0"
              >
                <v-chip :color="colors[index]" variant="outlined">
                  {{ config.NbQuestionsByCategories[index] }}
                </v-chip>
              </v-col>
            </v-row>
          </v-col>
        </v-row>
      </v-list-item>
    </v-list>
  </v-card>
</template>

<script setup lang="ts">
import type {
  GroupStrategy,
  TrivialConfig,
  TrivialConfigExt
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { colorsPerCategorie } from "@/controller/trivial";
import { computed, onMounted } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import LaunchOptions from "../components/trivial/LaunchOptions.vue";
import SessionMonitor from "../components/trivial/SessionMonitor.vue";
import TagsSelector from "../components/trivial/TagsSelector.vue";

let allKnownTags = $ref<string[]>([]);

let editedConfig = $ref<TrivialConfig | null>(null);

let _configs = $ref<TrivialConfigExt[]>([]);

const configs = computed(() => {
  const a = _configs.map(v => v);
  a.sort((u, v) => u.Config.Id - v.Config.Id);
  return a;
});

let isLaunching = $ref(false);

let gameCode = $ref("");

const colors = colorsPerCategorie;

onMounted(async () => {
  const res = await controller.GetTrivialPoursuit();

  if (res === undefined) {
    return;
  }
  _configs = Object.values(res || {});

  const tags = await controller.EditorGetTags();
  allKnownTags = tags || [];
});

async function createConfig() {
  const res = await controller.CreateTrivialPoursuit(null);
  if (res === undefined) {
    return;
  }
  _configs.push(res);
}

async function updateConfig() {
  // remove empty categories
  editedConfig!.Questions = editedConfig!.Questions.map(q =>
    (q || []).filter(v => v && v.length != 0)
  );
  const res = await controller.UpdateTrivialPoursuit(editedConfig!);
  if (res === undefined) {
    return;
  }
  const index = _configs.findIndex(v => v.Config.Id == editedConfig!.Id);
  _configs[index] = res;
  editedConfig = null;
}

async function deleteConfig(config: TrivialConfig) {
  await controller.DeleteTrivialPoursuit({ id: config.Id });
  _configs = _configs.filter(c => c.Config.Id != config.Id);
}

let launchingConfig = $ref<TrivialConfig | null>(null);
async function launchSession(options: GroupStrategy) {
  if (launchingConfig == null) {
    return;
  }
  const configID = launchingConfig.Id;
  isLaunching = true;
  const res = await controller.LaunchSessionTrivialPoursuit({
    IdConfig: configID,
    GroupStrategy: options
  });
  launchingConfig = null;
  isLaunching = false;
  if (res === undefined) {
    return;
  }

  const index = _configs.findIndex(v => v.Config.Id == configID);
  _configs[index].Running = res;
}

async function stopSession(config: TrivialConfig) {
  const configID = config.Id;
  const res = await controller.StopSessionTrivialPoursuit({ id: config.Id });
  if (res === undefined) {
    return;
  }

  const index = _configs.findIndex(v => v.Config.Id == configID);
  _configs[index].Running = {
    SessionID: "",
    GroupStrategyKind: 0,
    GroupsID: []
  };
}

let sessionToMonitor = $ref<TrivialConfigExt | null>(null);
function monitor(config: TrivialConfigExt) {
  sessionToMonitor = config;
}
</script>

<style>
:deep(.v-dialog .v-overlay__content) {
  max-width: 900px;
}
</style>
