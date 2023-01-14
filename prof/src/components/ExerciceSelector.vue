<template>
  <v-card class="mt-3 px-2" title="Choisir un exercice">
    <v-card-text>
      <resource-query-row
        v-model="props.query"
        :all-tags="props.tags"
        @update:model-value="updateQuery"
      >
      </resource-query-row>

      <div style="height: 47vh" class="overflow-y-auto">
        <v-expansion-panels class="pa-2">
          <v-expansion-panel v-for="(group, index) in exercices" :key="index">
            <v-expansion-panel-title>
              {{ group.Group.Title }}
            </v-expansion-panel-title>
            <v-expansion-panel-text>
              <v-list>
                <v-list-item
                  link
                  v-for="(exercice, index) in group.Variants"
                  :key="index"
                  @click="emit('selected', exercice)"
                >
                  <v-row>
                    <v-col>
                      <small>({{ exercice.Id }})</small> {{ exercice.Subtitle }}
                    </v-col>
                    <v-spacer></v-spacer>
                    <v-col cols="auto">
                      <tag-chip
                        :tag="{
                          Tag: exercice.Difficulty || 'Aucune difficulté',
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
        {{ exercices.length }} / {{ serverNbExercices }} variantes d'exercices
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
  ExercicegroupExt,
  ExerciceHeader,
  Query,
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
  (e: "selected", question: ExerciceHeader): void;
  (e: "update:query", query: Query): void;
}>();

let exercices = $ref<ExercicegroupExt[]>([]);
let serverNbExercices = $ref(0);

onMounted(() => {
  if (
    props.query.TitleQuery ||
    props.query.LevelTags?.length ||
    props.query.ChapterTags?.length
  ) {
    fetchExercices();
  }
});

async function updateQuery(qu: Query) {
  emit("update:query", qu);
  await fetchExercices();
}

async function fetchExercices() {
  const result = await controller.EditorSearchExercices(props.query);
  if (result == undefined) return;
  exercices = result.Groups || [];
  serverNbExercices = result.NbExercices;
}
</script>

<style></style>
