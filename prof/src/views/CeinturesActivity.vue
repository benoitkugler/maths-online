<template>
  <v-container class="pa-2 fill-height" fluid>
    <v-dialog v-model="showStudentsAdvance" v-if="scheme != null">
      <StudentsAdvance :classrooms="scheme.Classrooms || []"></StudentsAdvance>
    </v-dialog>
    <v-fade-transition hide-on-leave>
      <v-skeleton-loader
        v-if="scheme == null"
        width="600"
        class="mx-auto"
        type="card"
      ></v-skeleton-loader>
      <template v-else-if="scheme.IsAdmin">
        <v-row v-if="stageToEdit == null" justify="space-evenly">
          <v-col cols="7" align-self="center">
            <v-alert v-if="missingQuestions.length" icon="mdi-alert" closable>
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
              <template v-slot:append>
                <v-btn @click="showStudentsAdvance = true">
                  <v-icon>mdi-view-list</v-icon>
                  Progression élèves</v-btn
                >
              </template>
              <v-card-text class="pa-1">
                <table
                  style="table-layout: fixed; width: 200%"
                  v-if="scheme != null"
                >
                  <tbody>
                    <tr>
                      <th></th>
                      <th v-for="(k, v) in DomainLabels" :key="v" class="pa-2">
                        <v-card
                          link
                          @click="currentDomain = v"
                          :color="
                            currentDomain == v ? 'grey-lighten-3' : undefined
                          "
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
                        <StageHeaderCard
                          v-if="simulateProgression == null"
                          :stage="stage"
                          :header="scheme.Stages[stage.Domain][stage.Rank]"
                          @click="stageToEdit = stage"
                        >
                        </StageHeaderCard>
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
                  </tbody>
                </table>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="5" align-self="center" v-if="currentDomain != null">
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
                    :stages="scheme.Stages[currentDomain]"
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
      </template>
      <v-row v-else>
        <v-col class="text-center">
          <v-btn @click="showStudentsAdvance = true">
            <v-icon>mdi-view-list</v-icon>
            Afficher la progression élève</v-btn
          >
        </v-col>
      </v-row>
    </v-fade-transition>
  </v-container>
</template>

<script setup lang="ts">
import DomainLine from "@/components/ceintures/DomainLine.vue";
import StudentsAdvance from "@/components/ceintures/StudentsAdvance.vue";
import StageChip from "@/components/ceintures/StageChip.vue";
import StageHeaderCard from "@/components/ceintures/StageHeaderCard.vue";
import StageQuestionsEditor from "@/components/ceintures/StageQuestionsEditor.vue";
import { Advance, Domain, Rank, Stage } from "@/controller/api_gen";
import { GetSchemeOut, DomainLabels, RankLabels } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { sameStage } from "@/controller/utils";
import { computed } from "vue";
import { ref } from "vue";
import { onMounted } from "vue";
import { onActivated } from "vue";

onMounted(fetchScheme);
onActivated(fetchScheme);

const scheme = ref<GetSchemeOut | null>(null);

async function fetchScheme() {
  const res = await controller.CeinturesGetScheme();
  if (res === undefined) return;
  scheme.value = res;
}

const missingQuestions = computed(() => {
  const out: Stage[] = [];
  scheme.value?.Stages.forEach((ar, domain) =>
    ar.forEach((stage, rank) => {
      if (rank == Rank.StartRank) return;
      if (!stage.Questions?.length)
        out.push({ Domain: domain as Domain, Rank: rank as Rank });
    })
  );
  return out;
});

function initAdvance(): Advance {
  return Array.from({ length: nbDomains }).map(() => Rank.StartRank) as Advance;
}

const showStudentsAdvance = ref(false);

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

// show/hide question editor
const stageToEdit = ref<Stage | null>(null);
</script>

<style></style>
