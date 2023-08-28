<template>
  <v-card class="pt-2 pb-0">
    <v-dialog
      :model-value="groupToDelete != null"
      @update:model-value="groupToDelete = null"
      max-width="700"
    >
      <v-card title="Confirmer" v-if="groupToDelete != null">
        <v-card-text
          >Etes-vous certain de vouloir supprimer l'exercice
          <i>{{ groupToDelete.Group.Title }}</i>
          ?
          <br />
          <br />
          Cette opération est irréversible.
        </v-card-text>
        <v-card-actions>
          <v-btn @click="groupToDelete = null">Retour</v-btn>
          <v-spacer></v-spacer>
          <v-btn color="red" @click="deleteGroup" variant="outlined">
            Supprimer
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog
      :model-value="deletedBlocked != null"
      @update:model-value="deletedBlocked = null"
      max-width="600"
    >
      <uses-card
        :uses="deletedBlocked || []"
        @back="deletedBlocked = null"
        @go-to="goToSheet"
      ></uses-card>
    </v-dialog>

    <v-dialog
      :model-value="reviewToCreate != null"
      @update:model-value="reviewToCreate = null"
      max-width="700"
    >
      <confirm-publish @create-review="createReview"></confirm-publish>
    </v-dialog>

    <v-row>
      <v-col cols="auto" align-self="center">
        <v-btn
          class="ml-2"
          size="small"
          icon
          title="Retour au sommaire"
          @click="emit('back')"
        >
          <v-icon icon="mdi-arrow-left"></v-icon>
        </v-btn>
      </v-col>

      <v-col> <v-card-title>Liste des exercices</v-card-title> </v-col>

      <v-col align-self="center" style="text-align: right" cols="4">
        <v-btn
          class="mx-2"
          @click="createExercicegroup"
          title="Créer un nouvel exercice"
        >
          <v-icon icon="mdi-plus" color="success"></v-icon>
          Créer
        </v-btn>
      </v-col>
    </v-row>

    <v-card-text class="py-1">
      <resource-query-row
        :all-tags="props.tags"
        :model-value="query"
        @update:model-value="updateQuery"
      ></resource-query-row>
      <v-row no-gutters>
        <v-col>
          <v-list style="height: 62vh" class="overflow-y-auto">
            <div
              v-if="groups.length == 0"
              style="width: 100%; text-align: center"
            >
              <i>
                {{
                  query.TitleQuery == "" &&
                  !query.LevelTags?.length &&
                  !query.ChapterTags?.length
                    ? "Entrer une recherche..."
                    : "Aucun résultat"
                }}
              </i>
            </div>

            <div v-for="(exerciceGroup, index) in displayedGroups" :key="index">
              <resource-group-row
                :group="exerciceToResource(exerciceGroup)"
                :all-tags="props.tags"
                :is-question="false"
                @clicked="startEdit(exerciceGroup)"
                @duplicate="duplicate(exerciceGroup)"
                @delete="groupToDelete = exerciceGroup"
                @create-review="reviewToCreate = exerciceGroup.Group"
                @update-public="b => updatePublic(exerciceGroup.Group.Id, b)"
              ></resource-group-row>
            </div>
          </v-list>
          <v-row no-gutters>
            <v-col align-self="center" cols="4">
              {{ groups.length || 0 }} exercices ({{ serverNbExercices }}
              variantes)
            </v-col>
            <v-col align-self="center" cols="8">
              <v-pagination
                density="comfortable"
                rounded="circle"
                v-model="currentPage"
                :length="pageLength"
              ></v-pagination>
            </v-col>
          </v-row>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import {
  OriginKind,
  ReviewKind,
  type Exercicegroup,
  type ExercicegroupExt,
  type Query,
  type TagsDB,
  PublicStatus,
  type TaskUses,
  type Sheet
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { computed, onActivated, onMounted } from "@vue/runtime-core";
import { useRouter } from "vue-router";
import { $ref } from "vue/macros";
import ResourceQueryRow from "../ResourceQueryRow.vue";
import { exerciceToResource } from "@/controller/editor";
import { ref } from "vue";
import ResourceGroupRow from "../ResourceGroupRow.vue";
import UsesCard from "../UsesCard.vue";
import ConfirmPublish from "@/components/ConfirmPublish.vue";

interface Props {
  tags: TagsDB; // queried once for all
  initialQuery: Query | null;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "edit", group: ExercicegroupExt): void;
  (e: "back"): void;
}>();

defineExpose({ createExercicegroup });

const router = useRouter();

// groups are cut in slice of `pagination` length,
// and currentPage is the index of the page
const pagination = 6;
let currentPage = $ref(1);
const pageLength = computed(() => Math.ceil(groups.length / pagination));
const displayedGroups = computed(() =>
  groups.slice((currentPage - 1) * pagination, currentPage * pagination)
);

let groups = $ref<ExercicegroupExt[]>([]);
let serverNbExercices = $ref(0);

let query = $ref<Query>(
  props.initialQuery
    ? props.initialQuery
    : {
        TitleQuery: "",
        LevelTags: [],
        ChapterTags: [],
        SubLevelTags: [],
        Origin: OriginKind.All
      }
);
onMounted(fetchExercices);
onActivated(fetchExercices);

async function updateQuery(qu: Query) {
  query = qu;
  await fetchExercices();
}

async function duplicate(group: ExercicegroupExt) {
  const ok = await controller.EditorDuplicateExercicegroup({
    id: group.Group.Id
  });
  if (!ok) return;
  await fetchExercices();
}

const groupToDelete = ref<ExercicegroupExt | null>(null);
async function deleteGroup() {
  if (groupToDelete.value == null) return;
  const res = await controller.EditorDeleteExercicegroup({
    id: groupToDelete.value.Group.Id
  });
  if (res === undefined) return;
  if (!res.Deleted) {
    deletedBlocked.value = res.BlockedBy;
    return;
  }
  await fetchExercices();
}

const deletedBlocked = ref<TaskUses>(null);
function goToSheet(sh: Sheet) {
  deletedBlocked.value = null;

  router.push({ name: "homework", query: { idSheet: sh.Id } });
}

async function fetchExercices() {
  const result = await controller.EditorSearchExercices(query);
  if (result == undefined) {
    return;
  }
  groups = result.Groups || [];
  serverNbExercices = result.NbExercices;
}

async function createExercicegroup() {
  const out = await controller.EditorCreateExercice();
  if (out == undefined) {
    return;
  }
  await startEdit(out);
}

async function startEdit(group: ExercicegroupExt) {
  // we defer loading of the exercice in the pannel
  emit("edit", group);
}

async function updatePublic(id: number, isPublic: boolean) {
  const res = await controller.EditorUpdateExercicegroupVis({
    ID: id,
    Public: isPublic
  });
  if (res === undefined) {
    return;
  }

  const index = groups.findIndex(gr => gr.Group.Id == id);
  groups[index].Origin.PublicStatus = isPublic
    ? PublicStatus.AdminPublic
    : PublicStatus.AdminNotPublic;
}

const reviewToCreate = ref<Exercicegroup | null>(null);
async function createReview() {
  if (reviewToCreate.value == null) return;
  const res = await controller.ReviewCreate({
    Kind: ReviewKind.KExercice,
    Id: reviewToCreate.value.Id
  });
  reviewToCreate.value = null;
  if (res == undefined) return;
  router.push({ name: "reviews", query: { id: res.Id } });
}
</script>
