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

  <v-dialog
    :model-value="reviewToCreate != null"
    @update:model-value="reviewToCreate = null"
    max-width="700"
  >
    <confirm-publish @create-review="createReview"></confirm-publish>
  </v-dialog>

  <v-card
    class="my-5 mx-auto"
    width="90%"
    title="Isy'Triv"
    subtitle="Configurer et lancer une partie Isy'Triv"
  >
    <template v-slot:append>
      <v-btn
        size="small"
        @click="createConfig"
        title="Créer une nouvelle partie"
      >
        <v-icon icon="mdi-plus" color="success"></v-icon>
        Créer
      </v-btn>
      <matiere-select
        v-model:matiere="matiere"
        @update:matiere="_init"
      ></matiere-select>
    </template>

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
              isLaunching || !config.NbQuestionsByCategories.every((v) => v > 0)
            "
            @update-public="(b:boolean) => updatePublic(config.Config, b)"
            @create-review="reviewToCreate = config.Config"
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
  PublicStatus,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { ref, computed, onMounted, onActivated, watchEffect } from "vue";
import TrivialRow from "../components/trivial/TrivialRow.vue";
import EditConfig from "../components/trivial/EditConfig.vue";
import LaunchOptions from "../components/trivial/LaunchOptions.vue";
import SessionMonitor from "../components/trivial/SessionMonitor.vue";
import { useRouter } from "vue-router";
import { emptyTagsDB } from "@/controller/editor";
import ConfirmPublish from "@/components/ConfirmPublish.vue";
import MatiereSelect from "@/components/MatiereSelect.vue";
import SelfaccessConfig from "@/components/trivial/SelfaccessConfig.vue";

const router = useRouter();

const allKnownTags = ref<TagsDB>(emptyTagsDB());

const editedConfig = ref<Trivial | null>(null);

const matiere = ref(controller.settings.FavoriteMatiere);

const _configs = ref<TrivialExt[]>([]);

const configsByLevels = computed(() => {
  const byLevel = new Map<string, TrivialExt[]>();
  _configs.value.forEach((cf) => {
    if (!cf.Levels?.length) {
      // add to unclassified
      const l = byLevel.get("") || [];
      l.push(cf);
      byLevel.set("", l);
    } else {
      cf.Levels.forEach((level) => {
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

  // show unclassified first
  const unclassified = byLevel.get("");
  const others = Array.from(byLevel.entries()).filter((v) => v[0] != "");
  others.sort((a, b) => -a[0].localeCompare(b[0]));
  const head: typeof others = [["", unclassified || []]];
  return unclassified?.length ? head.concat(others) : others;
});

const isLaunching = ref(false);

onMounted(_init);
onActivated(_init);

async function _init() {
  fetchSessionMeta();

  await controller.ensureSettings();
  matiere.value = controller.settings.FavoriteMatiere;

  const res = await controller.GetTrivialPoursuit({ matiere: matiere.value });

  if (res === undefined) return;
  _configs.value = Object.values(res || {});

  const tags = await controller.EditorGetTags();
  if (tags) {
    allKnownTags.value = tags;
  }
}

async function createConfig() {
  const res = await controller.CreateTrivialPoursuit({
    matiere: matiere.value,
  });
  if (res === undefined) {
    return;
  }
  _configs.value.push(res);
  // launch the edition
  editedConfig.value = res.Config;
}

async function updateConfig(config: Trivial) {
  // remove empty categories
  config.Questions.Tags.forEach(
    (q, i) =>
      (config.Questions.Tags[i] = (q || []).filter((v) => v && v.length != 0))
  );
  const res = await controller.UpdateTrivialPoursuit(config);
  if (res === undefined) {
    return;
  }
  const index = _configs.value.findIndex((v) => v.Config.Id == config.Id);
  _configs.value[index] = res;
  editedConfig.value = null;
}

async function updatePublic(config: Trivial, isPublic: boolean) {
  const res = await controller.UpdateTrivialVisiblity({
    ConfigID: config.Id,
    Public: isPublic,
  });
  if (res === undefined) {
    return;
  }
  const index = _configs.value.findIndex((v) => v.Config.Id == config.Id);
  _configs.value[index].Origin.PublicStatus = isPublic
    ? PublicStatus.AdminPublic
    : PublicStatus.AdminNotPublic;
}

const reviewToCreate = ref<Trivial | null>(null);
async function createReview() {
  if (reviewToCreate.value == null) return;
  const res = await controller.ReviewCreate({
    Kind: ReviewKind.KTrivial,
    Id: reviewToCreate.value.Id,
  });
  reviewToCreate.value = null;
  if (res == undefined) return;

  router.push({ name: "reviews", query: { id: res.Id } });
}

async function duplicateConfig(config: Trivial) {
  const res = await controller.DuplicateTrivialPoursuit({ id: config.Id });
  if (res === undefined) {
    return;
  }
  console.log(config, res);

  _configs.value.push(res);
}

const trivialToDelete = ref<Trivial | null>(null);

async function deleteConfig() {
  if (trivialToDelete.value == null) return;
  const id = trivialToDelete.value.Id;
  const res = await controller.DeleteTrivialPoursuit({ id: id });
  trivialToDelete.value = null;
  if (res === undefined) return;
  _configs.value = _configs.value.filter((c) => c.Config.Id != id);
}

const selfaccessConfig = ref<Trivial | null>(null);

const launchingConfig = ref<Trivial | null>(null);
// workaround for https://github.com/vuetifyjs/vuetify/issues/16770
watchEffect(() => {
  document.documentElement.style.overflow =
    launchingConfig.value != null ? "hidden" : "";
});
async function launchSession(groups: GroupsStrategy) {
  if (launchingConfig.value == null) {
    return;
  }
  const configID = launchingConfig.value.Id;
  isLaunching.value = true;
  const res = await controller.LaunchSessionTrivialPoursuit({
    IdConfig: configID,
    Groups: groups,
  });
  launchingConfig.value = null;
  isLaunching.value = false;
  if (res === undefined) {
    return;
  }
  fetchSessionMeta();

  // automatically jump to monitor screen
  showMonitor.value = true;
}

const sessionMeta = ref<RunningSessionMetaOut>({ NbGames: 0 });
async function fetchSessionMeta() {
  const res = await controller.GetTrivialRunningSessions();
  if (res == undefined) {
    return;
  }
  sessionMeta.value = res;
}

const showMonitor = ref(false);
function showMonitorChanged(show: boolean) {
  showMonitor.value = show;
  if (!show) {
    closeMonitor();
  }
}

function closeMonitor() {
  showMonitor.value = false;
  fetchSessionMeta();
}
</script>

<style>
:deep(.v-dialog .v-overlay__content) {
  max-width: 900px;
}
</style>
