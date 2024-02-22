<template>
  <v-container class="pb-1 fill-height" fluid>
    <v-fade-transition hide-on-leave>
      <v-skeleton-loader
        v-if="scheme == null"
        width="600"
        class="mx-auto"
        type="card"
      ></v-skeleton-loader>
      <v-row v-else-if="stageToEdit == null" justify="space-evenly">
        <v-col cols="8" align-self="center">
          <v-alert v-if="missingQuestions.length" icon="mdi-alert">
            <v-row>
              <v-col align-self="center" cols="4"
                >Certaines ceintures sont vides !</v-col
              >
              <v-col align-self="center" class="text-center">
                <stage-chip
                  v-for="(stage, index) in missingQuestions"
                  :key="index"
                  :stage="stage"
                  :small="true"
                ></stage-chip>
              </v-col>
            </v-row>
          </v-alert>
          <v-card class="overflow-x-auto">
            <v-card-text class="pa-1">
              <table style="table-layout: fixed; width: 200%">
                <tr>
                  <th></th>
                  <th v-for="(k, v) in DomainLabels" :key="v" class="pa-2">
                    <v-card
                      link
                      @click="currentDomain = v"
                      :color="currentDomain == v ? 'grey-lighten-3' : undefined"
                      height="50px"
                      ><v-card-text class="pa-1 font-weight-bold">
                        {{ k }}
                      </v-card-text></v-card
                    >
                  </th>
                </tr>
                <tr v-for="rank in nbRanks - 1" :key="rank">
                  <th class="px-2">{{ RankLabels[rank as Rank] }}</th>

                  <td
                    v-for="(stage, index) in stagesFor(rank as Rank)"
                    :key="index"
                    class="pa-1"
                  >
                    <v-card
                      v-if="simulateProgression == null"
                      :color="rankColors[stage.Rank]"
                      link
                      @click="stageToEdit = stage"
                    >
                      <v-card-text class="text-center pa-1">
                        {{ stageText(stage) }}
                      </v-card-text>
                    </v-card>
                    <v-card
                      v-else
                      :color="
                        simulateProgression[stage.Domain] >= rank
                          ? 'green-lighten-2'
                          : isPending(stage)
                          ? 'blue-lighten-3'
                          : 'grey'
                      "
                      @click="
                        simulateProgression[stage.Domain] = (
                          simulateProgression[stage.Domain] >= rank
                            ? rank - 1
                            : rank
                        ) as Rank;
                        updatePending();
                      "
                      :disabled="
                        !(
                          simulateProgression[stage.Domain] >= rank ||
                          isPending(stage)
                        )
                      "
                    >
                      <v-card-text class="text-center pa-1">
                        <v-icon
                          :icon="
                            simulateProgression[stage.Domain] >= rank
                              ? 'mdi-check'
                              : isPending(stage)
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
        <v-col cols="4" align-self="center" v-if="currentDomain != null">
          <v-card
            title="Progression et prérequis"
            :subtitle="DomainLabels[currentDomain]"
          >
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

            <v-card-text class="overflow-y-auto" style="max-height: 72dvh">
              <v-fade-transition hide-on-leave>
                <domain-line
                  :domain="currentDomain"
                  :scheme="scheme.Scheme"
                  @go-to="(d) => (currentDomain = d)"
                ></domain-line>
              </v-fade-transition>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
      <stage-questions-editor
        v-else
        :stage="stageToEdit"
        :readonly="false"
        @back="
          stageToEdit = null;
          fetchScheme();
        "
        @go-to="(r) => (stageToEdit = { Domain: stageToEdit!.Domain, Rank: r })"
      ></stage-questions-editor>
    </v-fade-transition>
  </v-container>
</template>

<script setup lang="ts">
import DomainLine from "@/components/ceintures/DomainLine.vue";
import StageChip from "@/components/ceintures/StageChip.vue";
import StageQuestionsEditor from "@/components/ceintures/StageQuestionsEditor.vue";
import { Advance, Domain, Rank, Stage } from "@/controller/api_gen";
import { GetSchemeOut, DomainLabels, RankLabels } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { rankColors, sameStage } from "@/controller/utils";
import { computed } from "vue";
import { ref } from "vue";
import { onMounted } from "vue";

onMounted(fetchScheme);

const scheme = ref<GetSchemeOut | null>(null);

async function fetchScheme() {
  const res = await controller.CeinturesGetScheme();
  if (res === undefined) return;
  scheme.value = res;
}

const missingQuestions = computed(() => {
  const out: Stage[] = [];
  scheme.value?.NbQuestions.forEach((ar, domain) =>
    ar.forEach((nb, rank) => {
      if (rank == Rank.StartRank) return;
      if (nb == 0) out.push({ Domain: domain as Domain, Rank: rank as Rank });
    })
  );
  return out;
});

function initAdvance(): Advance {
  return Array.from({ length: nbDomains }).map(() => Rank.StartRank) as Advance;
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

const nbDomains = Domain.Matrices + 1;
const nbRanks = Rank.Noire + 1;

const currentDomain = ref<Domain | null>(0);

function stagesFor(rank: Rank): Stage[] {
  return Object.values(Domain).map((d) => ({ Domain: d, Rank: rank }));
}

function stageText(stage: Stage) {
  const s = scheme.value;
  if (s == null) return;
  return `${s.NbQuestions[stage.Domain][stage.Rank]} qu.`;
}

// show/hide question editor
const stageToEdit = ref<Stage | null>(null);
</script>

<style></style>
