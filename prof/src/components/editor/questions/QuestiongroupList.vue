<template>
  <v-card class="pt-2 pb-0">
    <v-dialog
      :model-value="groupToDelete != null"
      @update:model-value="groupToDelete = null"
      max-width="700"
    >
      <v-card title="Confirmer" v-if="groupToDelete != null">
        <v-card-text
          >Etes-vous certain de vouloir supprimer la question
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

      <v-col> <v-card-title>Liste des questions</v-card-title> </v-col>

      <v-col align-self="center" style="text-align: right" cols="4">
        <v-btn
          class="mx-2"
          @click="createQuestiongroup"
          title="Créer une nouvelle question"
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
          <v-list style="height: 60vh" class="overflow-y-auto">
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

            <div v-for="(questionGroup, index) in displayedGroups" :key="index">
              <resource-group-row
                :group="questionToResource(questionGroup)"
                :all-tags="props.tags"
                :is-question="true"
                @clicked="startEdit(questionGroup)"
                @duplicate="duplicate(questionGroup)"
                @delete="groupToDelete = questionGroup"
                @create-review="reviewToCreate = questionGroup.Group"
                @update-public="(b) => updatePublic(questionGroup.Group.Id, b)"
              ></resource-group-row>
            </div>
          </v-list>
          <v-row no-gutters>
            <v-col align-self="center" cols="4">
              {{ groups.length || 0 }} questions ({{ serverNbQuestions }}
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
  type Query,
  type Question,
  type Questiongroup,
  type QuestiongroupExt,
  type TagsDB,
  PublicStatus,
  type TaskUses,
  type Sheet,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { useRouter } from "vue-router";
import ResourceQueryRow from "../ResourceQueryRow.vue";
import ResourceGroupRow from "../ResourceGroupRow.vue";
import { questionToResource } from "@/controller/editor";
import { ref, computed, onActivated, onMounted } from "vue";
import UsesCard from "../UsesCard.vue";
import ConfirmPublish from "@/components/ConfirmPublish.vue";

interface Props {
  tags: TagsDB; // queried once for all
  initialQuery: Query | null;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "edit", group: QuestiongroupExt, questions: Question[]): void;
  (e: "back"): void;
}>();

defineExpose({ createQuestiongroup });

const router = useRouter();

// groups are cut in slice of `pagination` length,
// and currentPage is the index of the page
const pagination = 6;
const currentPage = ref(1);
const pageLength = computed(() => Math.ceil(groups.value.length / pagination));
const displayedGroups = computed(() =>
  groups.value.slice(
    (currentPage.value - 1) * pagination,
    currentPage.value * pagination
  )
);

const groups = ref<QuestiongroupExt[]>([]);
const serverNbQuestions = ref(0);

const query = ref<Query>(
  props.initialQuery
    ? props.initialQuery
    : {
        TitleQuery: "",
        LevelTags: [],
        ChapterTags: [],
        SubLevelTags: [],
        Origin: OriginKind.All,
        Matiere: controller.settings.FavoriteMatiere,
      }
);

onMounted(fetchQuestions);
onActivated(fetchQuestions);

async function updateQuery(qu: Query) {
  query.value = qu;
  await fetchQuestions();
}

async function fetchQuestions() {
  const result = await controller.EditorSearchQuestions(query.value);
  if (result == undefined) {
    return;
  }
  groups.value = result.Groups || [];
  serverNbQuestions.value = result.NbQuestions;
}

async function createQuestiongroup() {
  const out = await controller.EditorCreateQuestiongroup();
  if (out == undefined) {
    return;
  }
  await startEdit(out);
}

async function startEdit(group: QuestiongroupExt) {
  // load the variants
  const out = await controller.EditorGetQuestions({ id: group.Group.Id });
  if (out == undefined) {
    return;
  }
  emit("edit", group, out);
}

async function duplicate(group: QuestiongroupExt) {
  const ok = await controller.EditorDuplicateQuestiongroup({
    id: group.Group.Id,
  });
  if (!ok) return;
  await fetchQuestions();
}

const groupToDelete = ref<QuestiongroupExt | null>(null);
async function deleteGroup() {
  if (groupToDelete.value == null) return;
  const res = await controller.EditorDeleteQuestiongroup({
    id: groupToDelete.value.Group.Id,
  });
  groupToDelete.value = null;
  if (res === undefined) return;

  if (!res.Deleted) {
    deletedBlocked.value = res.BlockedBy;
    return;
  }
  await fetchQuestions();
}

const deletedBlocked = ref<TaskUses>(null);
function goToSheet(sh: Sheet) {
  deletedBlocked.value = null;

  router.push({ name: "homework", query: { idSheet: sh.Id } });
}

async function updatePublic(id: number, isPublic: boolean) {
  const res = await controller.EditorUpdateQuestiongroupVis({
    ID: id,
    Public: isPublic,
  });
  if (res === undefined) {
    return;
  }

  const index = groups.value.findIndex((gr) => gr.Group.Id == id);
  groups.value[index].Origin.PublicStatus = isPublic
    ? PublicStatus.AdminPublic
    : PublicStatus.AdminNotPublic;
}

const reviewToCreate = ref<Questiongroup | null>(null);
async function createReview() {
  if (reviewToCreate.value == null) return;
  const res = await controller.ReviewCreate({
    Kind: ReviewKind.KQuestion,
    Id: reviewToCreate.value.Id,
  });
  reviewToCreate.value = null;
  if (res == undefined) return;
  router.push({ name: "reviews", query: { id: res.Id } });
}
</script>
