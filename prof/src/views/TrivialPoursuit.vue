<template>
  <v-dialog
    fullscreen
    :model-value="editedConfig != null"
    @update:model-value="editedConfig = null"
  >
    <v-card v-if="editedConfig != null">
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
          <v-list-subheader>Choix des questions</v-list-subheader>
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

  <v-card class="ma-3">
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
      <v-list-item v-for="config in configs">
        <v-list-item-media>
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
        </v-list-item-media>
        {{ config }}
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

let configs = $ref<TrivialConfig[]>([]);

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
  configs = Object.values(res || {});

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
  configs.push(res);
}
</script>

<style>
:deep(.v-dialog .v-overlay__content) {
  max-width: 900px;
}
</style>
