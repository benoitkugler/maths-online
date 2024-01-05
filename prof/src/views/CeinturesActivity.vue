<template>
  <v-container class="pb-1 fill-height">
    <v-skeleton-loader
      v-if="scheme == null"
      width="600"
      class="mx-auto"
      type="card"
    ></v-skeleton-loader>
    <v-row v-else justify="space-evenly">
      <v-col cols="auto">
        <v-card>
          <v-card-text>
            <table>
              <tr>
                <th></th>
                <th v-for="(k, v) in DomainLabels" :key="v" class="pa-2">
                  {{ k }}
                </th>
              </tr>
              <tr v-for="rank in nbRanks - 1 as Rank" :key="rank">
                <th class="px-2">{{ RankLabels[rank as Rank] }}</th>
                <td
                  v-for="(domainName, domain) in DomainLabels"
                  :key="domain"
                  class="pa-2"
                >
                  <v-card
                    v-if="simulateProgression == null"
                    :elevation="
                      sameStage(currentStage, { Domain: domain, Rank: rank })
                        ? 12
                        : 1
                    "
                    :color="rankColors[rank as Rank]"
                    @click="currentStage = { Domain: domain, Rank: rank }"
                  >
                    <v-card-text>
                      {{ stageText({ Domain: domain, Rank: rank }) }}
                    </v-card-text>
                  </v-card>
                  <v-card
                    v-else
                    :color="
                      simulateProgression[domain] >= rank
                        ? 'green-lighten-2'
                        : isPending({ Domain: domain, Rank: rank })
                        ? 'blue-lighten-3'
                        : 'grey'
                    "
                    @click="
                      simulateProgression[domain] =
                        simulateProgression[domain] >= rank ? rank - 1 : rank;
                      updatePending();
                    "
                    :disabled="
                      !(
                        simulateProgression[domain] >= rank ||
                        isPending({ Domain: domain, Rank: rank })
                      )
                    "
                  >
                    <v-card-text class="text-center">
                      <v-icon
                        :icon="
                          simulateProgression[domain] >= rank
                            ? 'mdi-check'
                            : isPending({ Domain: domain, Rank: rank })
                            ? 'mdi-play'
                            : 'mdi-lock'
                        "
                      >
                      </v-icon>
                    </v-card-text>
                  </v-card>
                </td>
              </tr>
            </table>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col
        cols="5"
        align-self="center"
        v-if="currentStage != null && requiert != null"
      >
        <v-card title="Progression">
          <template v-slot:append>
            <v-btn
              v-if="simulateProgression == null"
              @click="
                simulateProgression = initAdvance();
                updatePending();
              "
              >Simuler...</v-btn
            >
            <v-btn v-else @click="simulateProgression = null">Retour</v-btn>
          </template>
          <v-card-text>
            <v-fade-transition hide-on-leave>
              <div v-if="simulateProgression == null">
                <template v-if="requiert.needs.length">
                  <v-row no-gutters justify="center"
                    ><v-col cols="auto">
                      <stage-chip
                        v-for="(need, i) in requiert.needs"
                        :key="i"
                        :stage="need"
                        link
                        @click="currentStage = need"
                      ></stage-chip> </v-col
                  ></v-row>
                  <v-row no-gutters justify="center">
                    <v-col cols="auto">
                      <v-icon>mdi-arrow-down</v-icon>
                    </v-col>
                  </v-row>
                </template>
                <v-row
                  no-gutters
                  justify="center"
                  class="bg-grey-lighten-3 rounded-lg"
                >
                  <v-col cols="auto">
                    <stage-chip :stage="currentStage"></stage-chip> </v-col
                ></v-row>
                <template v-if="requiert.isNeededIn.length">
                  <v-row no-gutters justify="center">
                    <v-col cols="auto">
                      <v-icon>mdi-arrow-down</v-icon>
                    </v-col>
                  </v-row>
                  <v-row no-gutters justify="center"
                    ><v-col cols="auto">
                      <stage-chip
                        v-for="(need, i) in requiert.isNeededIn"
                        :key="i"
                        :stage="need"
                        link
                        @click="currentStage = need"
                      ></stage-chip> </v-col
                  ></v-row>
                </template>
              </div>
              <div v-else>
                <v-row class="mt-2">
                  <v-col align-self="center" cols="auto">
                    <v-icon>mdi-chevron-left</v-icon>
                  </v-col>
                  <v-col align-self="center">
                    Sélectionner l'avancement pour afficher les niveaux
                    atteignables...
                  </v-col>
                </v-row>
              </div>
            </v-fade-transition>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import StageChip from "@/components/ceintures/StageChip.vue";
import { Advance, Beltquestion, Rank, Stage } from "@/controller/api_gen";
import { GetSchemeOut, DomainLabels, RankLabels } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { rankColors, sameStage } from "@/controller/utils";
import { computed, ref } from "vue";
import { onMounted } from "vue";

onMounted(fetchScheme);

const scheme = ref<GetSchemeOut | null>(null);

async function fetchScheme() {
  const res = await controller.CeinturesGetScheme();
  if (res === undefined) return;
  scheme.value = res;
}

function initAdvance(): Advance {
  return Array.from({ length: nbDomains }).map(() => Rank.StartRank);
}

const simulateProgression = ref<Advance | null>(null);

const pending = ref<Stage[]>([]);
async function updatePending() {
  if (simulateProgression.value == null) return;
  const res = await controller.CeinturesGetPending(simulateProgression.value);
  if (res === undefined) return;
  pending.value = res || [];
}
function isPending(stage: Stage) {
  return pending.value.find((s) => sameStage(s, stage)) !== undefined;
}

const nbDomains = Object.values(DomainLabels).length;
const nbRanks = Object.values(RankLabels).length;

const currentStage = ref<Stage>({ Domain: 0, Rank: 1 });

const questionsByStage = computed(() => {
  const out: Beltquestion[][][] = Array.from({ length: nbDomains }).map(() =>
    Array.from({ length: nbRanks }).map(() => [])
  );
  for (const question of scheme.value?.Questions || []) {
    out[question.Domain][question.Rank].push(question);
  }
  return out;
});

function stageText(stage: Stage) {
  return `${
    questionsByStage.value[stage.Domain][stage.Rank].length
  } question(s)`;
}

interface deps {
  needs: Stage[];
  isNeededIn: Stage[];
}

const requiert = computed(() => {
  const stage = currentStage.value;
  if (stage == null || scheme.value == null) return null;
  const out: deps = { needs: [], isNeededIn: [] };
  for (const link of scheme.value.Scheme || []) {
    if (sameStage(link.For, stage)) {
      out.needs.push(link.Need);
    } else if (sameStage(link.Need, stage)) {
      out.isNeededIn.push(link.For);
    }
  }
  // always include same domain side ranks
  if (stage.Rank > Rank.Blanche) {
    out.needs.push({ Domain: stage.Domain, Rank: stage.Rank - 1 });
  }
  if (stage.Rank < nbRanks - 1) {
    out.isNeededIn.push({ Domain: stage.Domain, Rank: stage.Rank + 1 });
  }

  return out;
});
</script>

<style></style>
