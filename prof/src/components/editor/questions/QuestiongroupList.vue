<template>
  <v-card class="pt-2 pb-0">
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
              <questiongroup-row
                :group="questionGroup"
                :all-tags="props.tags"
                @clicked="startEdit(questionGroup)"
                @duplicate="duplicate(questionGroup)"
                @update-public="updatePublic"
                @create-review="createReview(questionGroup.Group)"
              ></questiongroup-row>
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
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { computed, onActivated, onMounted } from "@vue/runtime-core";
import { useRouter } from "vue-router";
import { $ref } from "vue/macros";
import QuestiongroupRow from "./QuestiongroupRow.vue";
import ResourceQueryRow from "../ResourceQueryRow.vue";

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
let currentPage = $ref(1);
const pageLength = computed(() => Math.ceil(groups.length / pagination));
const displayedGroups = computed(() =>
  groups.slice((currentPage - 1) * pagination, currentPage * pagination)
);

let groups = $ref<QuestiongroupExt[]>([]);
let serverNbQuestions = $ref(0);

let query = $ref<Query>(
  props.initialQuery
    ? props.initialQuery
    : {
        TitleQuery: "",
        LevelTags: [],
        ChapterTags: [],
        Origin: OriginKind.All,
      }
);

onMounted(fetchQuestions);
onActivated(fetchQuestions);

async function updateQuery(qu: Query) {
  query = qu;
  await fetchQuestions();
}

async function fetchQuestions() {
  const result = await controller.EditorSearchQuestions(query);
  if (result == undefined) {
    return;
  }
  groups = result.Groups || [];
  serverNbQuestions = result.NbQuestions;
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

async function updatePublic(id: number, isPublic: boolean) {
  const res = await controller.EditorUpdateQuestiongroupVis({
    ID: id,
    Public: isPublic,
  });
  if (res === undefined) {
    return;
  }

  const index = groups.findIndex((gr) => gr.Group.Id == id);
  groups[index].Origin.IsPublic = isPublic;
}

async function createReview(ex: Questiongroup) {
  const res = await controller.ReviewCreate({
    Kind: ReviewKind.KQuestion,
    Id: ex.Id,
  });
  if (res == undefined) return;
  router.push({ name: "reviews", query: { id: res.Id } });
}
</script>
