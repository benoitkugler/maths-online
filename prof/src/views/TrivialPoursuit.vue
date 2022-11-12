<template>
  <v-dialog
    :model-value="editedConfig != null"
    @update:model-value="editedConfig = null"
    :retain-focus="false"
    max-width="1200"
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
    :model-value="trivialToDelete != null"
    @update:model-value="trivialToDelete = null"
    max-width="800px"
  >
    <v-card title="Confirmer">
      <v-card-text
        >Etes-vous certain de vouloir supprimer la partie
        <i>{{ trivialToDelete?.Name }}</i> ? <br />
        Cette opération est irréversible.
      </v-card-text>
      <v-card-actions>
        <v-btn @click="trivialToDelete = null">Retour</v-btn>
        <v-spacer></v-spacer>
        <v-btn color="red" @click="deleteConfig" variant="elevated">
          Supprimer
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog
    fullscreen
    :model-value="showMonitor"
    @update:model-value="showMonitorChanged"
    :retain-focus="false"
  >
    <session-monitor @closed="closeMonitor"></session-monitor>
  </v-dialog>

  <v-card class="my-5 mx-auto" width="90%">
    <v-row class="mx-0">
      <v-col cols="9">
        <v-card-title>Triv'Maths</v-card-title>
        <v-card-subtitle
          >Configurer et lancer une partie de Triv'Maths</v-card-subtitle
        >
      </v-col>

      <v-col align-self="center" style="text-align: right" cols="3">
        <v-btn
          size="small"
          @click="createConfig"
          title="Créer une nouvelle partie"
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
        <trivial-row
          :config="config"
          :disable-launch="
            isLaunching || !config.NbQuestionsByCategories.every((v) => v > 0)
          "
          @update-public="(b) => updatePublic(config.Config, b)"
          @create-review="createReview(config.Config)"
          @duplicate="duplicateConfig(config.Config)"
          @edit="editedConfig = config.Config"
          @launch="launchingConfig = config.Config"
          @delete="trivialToDelete = config.Config"
        ></trivial-row>
      </v-list-item>
    </v-list>
  </v-card>
</template>

<script setup lang="ts">
import {
  ReviewKind,
  type RunningSessionMetaOut,
  type Trivial,
  type TrivialExt,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { computed, onMounted } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import TrivialRow from "../components/trivial/TrivialRow.vue";
import EditConfig from "../components/trivial/EditConfig.vue";
import LaunchOptions from "../components/trivial/LaunchOptions.vue";
import SessionMonitor from "../components/trivial/SessionMonitor.vue";
import { useRouter } from "vue-router";

const router = useRouter();

let allKnownTags = $ref<string[]>([]);

let editedConfig = $ref<Trivial | null>(null);

let _configs = $ref<TrivialExt[]>([]);

const configs = computed(() => {
  const a = _configs.map((v) => v);
  a.sort((u, v) => u.Config.Id - v.Config.Id);
  return a;
});

let isLaunching = $ref(false);

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

async function createReview(config: Trivial) {
  const res = await controller.ReviewCreate({
    Kind: ReviewKind.KTrivial,
    Id: config.Id,
  });
  if (res == undefined) return;

  router.push({ name: "reviews", query: { id: res.Id } });
}

async function duplicateConfig(config: Trivial) {
  const res = await controller.DuplicateTrivialPoursuit({ id: config.Id });
  if (res === undefined) {
    return;
  }
  console.log(config, res);

  _configs.push(res);
}

let trivialToDelete = $ref<Trivial | null>(null);

async function deleteConfig() {
  if (trivialToDelete == null) return;
  const id = trivialToDelete.Id;
  await controller.DeleteTrivialPoursuit({ id: id });
  trivialToDelete = null;
  _configs = _configs.filter((c) => c.Config.Id != id);
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
