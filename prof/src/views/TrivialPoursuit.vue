<template>
  <v-dialog
    :model-value="editedConfig != null"
    @update:model-value="editedConfig = null"
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
        <v-btn @click="trivialToDelete = null" color="warning">Retour</v-btn>
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
  >
    <session-monitor @closed="closeMonitor"></session-monitor>
  </v-dialog>

  <v-dialog
    :model-value="selfaccessConfig != null"
    @update:model-value="selfaccessConfig = null"
    max-width="870px"
  >
    <selfaccess-config
      v-if="selfaccessConfig != null"
      :config="selfaccessConfig"
      @close="selfaccessConfig = null"
    ></selfaccess-config>
  </v-dialog>

  <v-card class="my-5 mx-auto" width="90%">
    <v-row class="mx-0">
      <v-col cols="9">
        <v-card-title>Isy'Triv</v-card-title>
        <v-card-subtitle
          >Configurer et lancer une partie Isy'Triv</v-card-subtitle
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
      <v-row justify="space-evenly">
        <v-col cols="auto" class="my-2">
          Parties en cours : <v-chip>{{ sessionMeta.NbGames }}</v-chip>
        </v-col>
        <v-col cols="auto" align-self="center">
          <v-btn @click="showMonitor = true"> Suivre les parties </v-btn>
        </v-col>
      </v-row>
    </v-alert>

    <v-card v-for="level in configsByLevels" :key="level[0]" class="ma-2">
      <v-card-title class="bg-pink-lighten-3">
        {{ level[0] || "Non classé" }}
      </v-card-title>
      <v-card-text class="mt-3">
        <v-row no-gutters>
          <trivial-row
            v-for="config in level[1]"
            :key="config.Config.Id"
            :config="config"
            :disable-launch="
              isLaunching || !config.NbQuestionsByCategories.every(v => v > 0)
            "
            @update-public="(b:boolean) => updatePublic(config.Config, b)"
            @create-review="createReview(config.Config)"
            @duplicate="duplicateConfig(config.Config)"
            @edit="editedConfig = config.Config"
            @launch="launchingConfig = config.Config"
            @delete="trivialToDelete = config.Config"
            @show-selfaccess="selfaccessConfig = config.Config"
          ></trivial-row>
        </v-row>
      </v-card-text>
    </v-card>
  </v-card>
</template>

<script setup lang="ts">
import {
  ReviewKind,
  type GroupsStrategy,
  type RunningSessionMetaOut,
  type TagsDB,
  type Trivial,
  type TrivialExt,
  PublicStatus
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import {
  computed,
  onMounted,
  onActivated,
  watchEffect
} from "@vue/runtime-core";
import { $ref } from "vue/macros";
import TrivialRow from "../components/trivial/TrivialRow.vue";
import EditConfig from "../components/trivial/EditConfig.vue";
import LaunchOptions from "../components/trivial/LaunchOptions.vue";
import SessionMonitor from "../components/trivial/SessionMonitor.vue";
import { useRouter } from "vue-router";
import { emptyTagsDB } from "@/controller/editor";
import SelfaccessConfig from "../components/trivial/SelfaccessConfig.vue";

const router = useRouter();

let allKnownTags = $ref<TagsDB>(emptyTagsDB());

let editedConfig = $ref<Trivial | null>(null);

let _configs = $ref<TrivialExt[]>([]);

const configsByLevels = computed(() => {
  const byLevel = new Map<string, TrivialExt[]>();
  _configs.forEach(cf => {
    if (!cf.Levels?.length) {
      // add to unclassified
      const l = byLevel.get("") || [];
      l.push(cf);
      byLevel.set("", l);
    } else {
      cf.Levels.forEach(level => {
        const l = byLevel.get(level) || [];
        l.push(cf);
        byLevel.set(level, l);
      });
    }
  });
  // inside each level, sort by id
  for (const list of byLevel.values()) {
    list.sort((u, v) => u.Config.Id - v.Config.Id);
  }

  const out = Array.from(byLevel.entries());
  out.sort((a, b) => -a[0].localeCompare(b[0]));

  return out;
});

let isLaunching = $ref(false);

onMounted(_init);
onActivated(_init);

async function _init() {
  fetchSessionMeta();

  const res = await controller.GetTrivialPoursuit();

  if (res === undefined) {
    return;
  }
  _configs = Object.values(res || {});

  const tags = await controller.EditorGetTags();
  if (tags) {
    allKnownTags = tags;
  }
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
  config.Questions.Tags = config.Questions.Tags.map(q =>
    (q || []).filter(v => v && v.length != 0)
  );
  const res = await controller.UpdateTrivialPoursuit(config);
  if (res === undefined) {
    return;
  }
  const index = _configs.findIndex(v => v.Config.Id == config.Id);
  _configs[index] = res;
  editedConfig = null;
}

async function updatePublic(config: Trivial, isPublic: boolean) {
  const res = await controller.UpdateTrivialVisiblity({
    ConfigID: config.Id,
    Public: isPublic
  });
  if (res === undefined) {
    return;
  }
  const index = _configs.findIndex(v => v.Config.Id == config.Id);
  _configs[index].Origin.PublicStatus = isPublic
    ? PublicStatus.AdminPublic
    : PublicStatus.AdminNotPublic;
}

async function createReview(config: Trivial) {
  const res = await controller.ReviewCreate({
    Kind: ReviewKind.KTrivial,
    Id: config.Id
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
  _configs = _configs.filter(c => c.Config.Id != id);
}

let selfaccessConfig = $ref<Trivial | null>(null);

let launchingConfig = $ref<Trivial | null>(null);
// workaround for https://github.com/vuetifyjs/vuetify/issues/16770
watchEffect(() => {
  document.documentElement.style.overflow =
    launchingConfig != null ? "hidden" : "";
});
async function launchSession(groups: GroupsStrategy) {
  if (launchingConfig == null) {
    return;
  }
  const configID = launchingConfig.Id;
  isLaunching = true;
  const res = await controller.LaunchSessionTrivialPoursuit({
    IdConfig: configID,
    Groups: groups
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
