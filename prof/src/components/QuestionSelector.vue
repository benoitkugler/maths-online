<template>
  <v-card class="mt-3 px-2" title="Choisir une question">
    <v-card-text>
      <resource-query-row
        v-model="props.query"
        :all-tags="props.tags"
        @update:model-value="updateQuery"
      >
      </resource-query-row>

      <div style="height: 47vh" class="overflow-y-auto">
        <v-expansion-panels class="pa-2">
          <v-expansion-panel v-for="(group, index) in questions" :key="index">
            <v-expansion-panel-title>
              {{ group.Group.Title }}
            </v-expansion-panel-title>
            <v-expansion-panel-text>
              <v-list>
                <v-list-item
                  link
                  v-for="(question, index) in group.Variants"
                  :key="index"
                  @click="emit('selected', question)"
                >
                  <v-row>
                    <v-col>
                      <small>({{ question.Id }})</small> {{ question.Subtitle }}
                    </v-col>
                    <v-spacer></v-spacer>
                    <v-col cols="auto">
                      <tag-chip
                        :tag="{
                          Tag: question.Difficulty || 'Aucune difficulté',
                          Section: 0,
                        }"
                      ></tag-chip>
                    </v-col>
                  </v-row>
                </v-list-item>
              </v-list>
            </v-expansion-panel-text>
          </v-expansion-panel>
        </v-expansion-panels>
      </div>
      <div class="my-2">
        {{ questions.length }} / {{ serverNbQuestions }} variantes de questions
        affichées
      </div>
    </v-card-text>
    <v-card-actions>
      <v-btn @click="emit('closed')" color="warning">Retour</v-btn>
      <v-spacer></v-spacer>
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import { onMounted } from "vue";
import { $ref } from "vue/macros";
import type {
  Query,
  QuestiongroupExt,
  QuestionHeader,
  TagsDB,
} from "../controller/api_gen";
import { controller } from "../controller/controller";
import TagChip from "./editor/utils/TagChip.vue";
import ResourceQueryRow from "./editor/ResourceQueryRow.vue";

interface Props {
  tags: TagsDB;
  query: Query;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  (e: "closed"): void;
  (e: "selected", question: QuestionHeader): void;
  (e: "update:query", query: Query): void;
}>();

let questions = $ref<QuestiongroupExt[]>([]);
let serverNbQuestions = $ref(0);

onMounted(() => {
  if (
    props.query.TitleQuery ||
    props.query.LevelTags?.length ||
    props.query.ChapterTags?.length
  ) {
    fetchQuestions();
  }
});

async function updateQuery(qu: Query) {
  emit("update:query", qu);
  await fetchQuestions();
}

async function fetchQuestions() {
  const result = await controller.EditorSearchQuestions(props.query);
  if (result == undefined) return;
  questions = result.Groups || [];
  serverNbQuestions = result.NbQuestions;
}
</script>

<style></style>
