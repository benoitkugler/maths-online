<template>
  <v-dialog
    :model-value="editedConfig != null"
    @update:model-value="editedConfig = null"
  >
    <v-card
      v-if="editedConfig != null"
      max-height="80vh"
      width="800"
      class="overflow-y-auto"
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
        <v-row justify="center">
          <v-col cols="6">
            <v-text-field
              density="comfortable"
              variant="outlined"
              label="Durée limite pour une question"
              type="number"
              min="1"
              suffix="sec"
              v-model.number="editedConfig!.QuestionTimeout"
            ></v-text-field>
          </v-col>
        </v-row>

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
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="success" @click="updateConfig">
          Enregistrer les modifications
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog title="Démarrer la session">
    <v-form class="mx-4 my-2">
      <v-text-field
        label="Nombre de joueurs"
        type="number"
        min="1"
        v-model.number="launchOptions.GroupStrategy.Data.MaxPlayersPerGroup"
      ></v-text-field>

      <v-row>
        <v-col cols="4" sm="6" md="9"></v-col>
        <v-col cols="8" sm="6" md="3" class="text-right">
          <v-btn block @click="launchGame" :disabled="!isValid" color="success">
            Lancer la partie
          </v-btn>
        </v-col>
      </v-row>

      <v-alert class="px-4 my-2" :model-value="gameCode != ''" color="info">
        <v-row no-gutters>
          <v-col align-self="center">
            Code de
            <a href="/test-eleve" target="_blank">la partie à rejoindre</a> :
          </v-col>
          <v-col>
            <v-chip size="big" class="pa-3"
              ><b>{{ gameCode }}</b></v-chip
            >
          </v-col>
        </v-row>
      </v-alert>
    </v-form>
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
              @click="editedConfig = config"
            >
              <v-icon icon="mdi-pencil"></v-icon>
            </v-btn>
            <v-btn icon size="x-small" title="Lancer" class="mx-2">
              <v-icon icon="mdi-play" color="green"></v-icon>
            </v-btn>
            <v-btn
              class="mx-2"
              size="x-small"
              icon
              @click="deleteConfig(config)"
              title="Supprimer cette session"
              :disabled="config.IsLaunched"
            >
              <v-icon icon="mdi-delete" color="red"></v-icon>
            </v-btn>
          </v-col>
          <v-col align-self="center" md="8">
            <v-row justify="center">
              <v-col
                cols="6"
                align-self="center"
                style="text-align: center"
                v-if="config.Questions!.every(v=>!v)"
              >
                <i>Aucune question</i>
              </v-col>
              <v-col
                cols="2"
                align-self="center"
                v-for="(categorie, index) in config.Questions || []"
                v-show="categorie && categorie.length != 0"
              >
                <v-chip :color="colors[index]" variant="outlined">AAAA</v-chip>
              </v-col>
            </v-row>
          </v-col>
        </v-row>
      </v-list-item>
    </v-list>
  </v-card>
</template>

<script setup lang="ts">
import type { LaunchSessionIn, TrivialConfig } from "@/controller/api_gen";
import { GroupStrategyKind } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { computed, onMounted, reactive } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import TagsSelector from "../components/trivial/TagsSelector.vue";

let launchOptions: LaunchSessionIn = reactive({
  IdConfig: 1,
  GroupStrategy: {
    Kind: GroupStrategyKind.RandomGroupStrategy,
    Data: { MaxPlayersPerGroup: 2, TotalPlayersNumber: 5 }
  }
});

let allKnownTags = $ref<string[]>([]);

let editedConfig = $ref<TrivialConfig | null>(null);

let _configs = $ref<TrivialConfig[]>([]);

const configs = computed(() => {
  const a = _configs.map(v => v);
  a.sort((u, v) => u.Id - v.Id);
  return a;
});

const isValid = computed(
  () => !isSpinning && launchOptions.GroupStrategy.Data.MaxPlayersPerGroup >= 1
);
let isSpinning = $ref(false);

let gameCode = $ref("");

const colors = ["purple", "green", "orange", "yellow", "blue"];

onMounted(async () => {
  const res = await controller.GetTrivialPoursuit();

  if (res === undefined) {
    return;
  }
  _configs = Object.values(res || {});

  const tags = await controller.EditorGetTags();
  allKnownTags = tags || [];
});

async function launchGame(): Promise<void> {
  isSpinning = true;
  const res = await controller.LaunchSession(launchOptions);
  isSpinning = false;

  if (res === undefined) {
    return;
  }

  gameCode = res.SessionID;
}

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
  await controller.UpdateTrivialPoursuit(editedConfig!);
  editedConfig = null;
}

async function deleteConfig(config: TrivialConfig) {
  await controller.DeleteTrivialPoursuit({ id: config.Id });
  _configs = _configs.filter(c => c.Id != config.Id);
}
</script>

<style>
:deep(.v-dialog .v-overlay__content) {
  max-width: 900px;
}
</style>
