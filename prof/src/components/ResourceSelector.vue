<template>
  <v-card class="mt-3 px-2">
    <v-row class="mt-0">
      <v-col>
        <v-card-title>Choisir une ressource</v-card-title>
      </v-col>
      <v-col cols="auto">
        <v-btn icon="mdi-close" flat @click="emit('closed')"> </v-btn>
      </v-col>
    </v-row>
    <v-card-text class="pt-1">
      <resource-query-row
        :model-value="props.query"
        :all-tags="props.tags"
        @update:model-value="updateQuery"
      >
      </resource-query-row>

      <div style="height: 62vh" class="overflow-y-auto mt-1">
        <v-card
          v-for="(group, index) in displayedGroups"
          :key="index"
          class="ma-1 pb-1"
        >
          <v-row class="py-3 px-4">
            <v-col align-self="center">
              {{ group.Title }}
            </v-col>
            <v-col
              cols="auto"
              align-self="center"
              v-if="props.mode == 'questions'"
            >
              <v-btn
                variant="outlined"
                size="small"
                @click="emit('selected-group', group)"
                >Inclure aléatoirement</v-btn
              >
            </v-col>
          </v-row>
          <v-card-text class="py-0">
            <v-list class="py-0" style="column-count: 2">
              <v-list-item
                rounded
                link
                density="compact"
                v-for="(variant, index) in group.Variants"
                :key="index"
                @click="emit('selected-variant', variant)"
                style="break-inside: avoid-column"
              >
                <v-row>
                  <v-col>
                    <small>({{ variant.Id }})</small> {{ variant.Subtitle }}
                  </v-col>
                  <v-spacer></v-spacer>
                  <v-col cols="auto" align-self="center">
                    <tag-chip
                      :tag="{
                        Tag: variant.Difficulty || 'Aucune difficulté',
                      }"
                    ></tag-chip>
                  </v-col>
                </v-row>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </div>
      <v-row no-gutters>
        <v-col align-self="center" cols="4">
          {{ groups.length || 0 }} ressource{{
            groups.length > 1 ? "s" : ""
          }}
          ({{ serverNbVariants }}
          variantes)
        </v-col>
        <v-col align-self="center" cols="8">
          <v-pagination
            density="comfortable"
            rounded="circle"
            v-model="currentPage"
            :length="Math.ceil(groups.length / pagination)"
          ></v-pagination>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { onMounted } from "vue";
import type { Query, TagsDB } from "../controller/api_gen";
import { controller } from "../controller/controller";
import TagChip from "./editor/utils/TagChip.vue";
import ResourceQueryRow from "./editor/ResourceQueryRow.vue";
import { computed } from "vue";
import {
  exerciceToResource,
  questionToResource,
  type ResourceGroup,
  type VariantG,
} from "@/controller/editor";
import { ref } from "vue";

interface Props {
  tags: TagsDB;
  query: Query;
  mode: "questions" | "exercices";
}

const props = defineProps<Props>();
const emit = defineEmits<{
  (e: "closed"): void;
  (e: "selected-variant", question: VariantG): void;
  (e: "selected-group", group: ResourceGroup): void;
  (e: "update:query", query: Query): void;
}>();

onMounted(() => {
  if (
    props.query.TitleQuery ||
    props.query.LevelTags?.length ||
    props.query.ChapterTags?.length
  ) {
    fetchResources();
  }
});

const groups = ref<ResourceGroup[]>([]);
const serverNbVariants = ref(0);
// groups are cut in slice of `pagination` length,
// and currentPage is the index of the page
const pagination = 6;
const currentPage = ref(1);
const displayedGroups = computed(() =>
  groups.value.slice(
    (currentPage.value - 1) * pagination,
    currentPage.value * pagination
  )
);

async function updateQuery(qu: Query) {
  emit("update:query", qu);
  await fetchResources();
}

async function fetchResources() {
  switch (props.mode) {
    case "exercices":
      {
        const result = await controller.EditorSearchExercices(props.query);
        if (result == undefined) return;
        groups.value = (result.Groups || []).map(exerciceToResource);
        serverNbVariants.value = result.NbExercices;
      }
      break;
    case "questions":
      {
        const result = await controller.EditorSearchQuestions(props.query);
        if (result == undefined) return;
        groups.value = (result.Groups || []).map(questionToResource);
        serverNbVariants.value = result.NbQuestions;
      }
      break;
  }
  currentPage.value = 1;
}
</script>

<style></style>
