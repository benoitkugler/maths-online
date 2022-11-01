<template>
  <v-card class="pt-2 pb-0">
    <v-row>
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
      <v-row>
        <v-col>
          <v-text-field
            label="Rechercher"
            hint="Rechercher une question par nom."
            variant="outlined"
            density="compact"
            v-model="querySearch"
            @update:model-value="updateQuerySearch"
            persistent-hint
            clearable
          ></v-text-field>
        </v-col>
        <v-col>
          <v-autocomplete
            variant="outlined"
            density="compact"
            multiple
            chips
            closable-chips
            :items="props.tags"
            color="primary"
            label="Catégories"
            no-data-text="Aucune catégorie n'est encore utilisée."
            v-model="queryTags"
            @update:model-value="updateQueryTags"
            @blur="updateQueryTags"
            hint="Restreint la recherche à l'intersection des catégories sélectionnées."
            persistent-hint
          ></v-autocomplete>
        </v-col>
        <v-col cols="auto" align-self="center">
          <origin-select
            :origin="queryOrigin"
            @update:origin="updateQueryOrigin"
          ></origin-select>
        </v-col>
      </v-row>
      <v-row no-gutters>
        <v-col>
          <v-list style="height: 60vh" class="overflow-y-auto">
            <div
              v-if="groups.length == 0"
              style="width: 100%; text-align: center"
            >
              <i>
                {{
                  querySearch == "" && queryTags.length == 0
                    ? "Entrer une recherche..."
                    : "Aucun résultat"
                }}
              </i>
            </div>

            <div v-for="questionGroup in groups" :key="questionGroup.Group.Id">
              <questiongroup-row
                :group="questionGroup"
                :all-tags="props.tags"
                @clicked="startEdit(questionGroup)"
                @duplicate="duplicate(questionGroup)"
                @update-public="updatePublic"
                @update-tags="
                  (tags) => updateGroupTags(questionGroup.Group, tags)
                "
              ></questiongroup-row>
            </div>
          </v-list>
          <div class="my-2">
            {{ groups.length }} / {{ serverNbGroups }} questions affichées. ({{
              displayedNbQuestions
            }}
            / {{ serverNbQuestions }}
            variantes)
          </div>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import {
  OriginKind,
  type Question,
  type Questiongroup,
  type QuestiongroupExt,
} from "@/controller/api_gen";
import { controller, IsDev } from "@/controller/controller";
import { computed, onActivated, onMounted } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import OriginSelect from "../../OriginSelect.vue";
import QuestiongroupRow from "./QuestiongroupRow.vue";

interface Props {
  tags: string[]; // queried once for all
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "edit", group: QuestiongroupExt, questions: Question[]): void;
}>();

let groups = $ref<QuestiongroupExt[]>([]);
let serverNbGroups = $ref(0);
let serverNbQuestions = $ref(0);
const displayedNbQuestions = computed(() => {
  let nb = 0;
  groups.forEach((group) => {
    nb += group.Variants?.length || 0;
  });
  return nb;
});

let querySearch = $ref("");

let queryTags = $ref<string[]>(IsDev ? ["DEV"] : []);

let queryOrigin = $ref(OriginKind.All);

let timerId = 0;

onMounted(fetchQuestions);
onActivated(fetchQuestions);

function updateQuerySearch() {
  const debounceDelay = 200;
  // cancel pending call
  clearTimeout(timerId);

  // delay new call 500ms
  timerId = setTimeout(() => {
    fetchQuestions();
  }, debounceDelay);
}

async function updateQueryTags() {
  await fetchQuestions();
}

async function updateQueryOrigin(o: OriginKind) {
  queryOrigin = o;
  await fetchQuestions();
}

async function fetchQuestions() {
  const result = await controller.EditorSearchQuestions({
    TitleQuery: querySearch,
    Tags: queryTags,
    Origin: queryOrigin,
  });
  if (result == undefined) {
    return;
  }
  groups = result.Groups || [];
  serverNbGroups = result.NbGroups;
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
  // load the questions
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

async function updateGroupTags(group: Questiongroup, newTags: string[]) {
  const rep = await controller.EditorUpdateQuestionTags({
    Id: group.Id,
    Tags: newTags,
  });
  if (rep == undefined) {
    return;
  }
  const index = groups.findIndex((gr) => gr.Group.Id == group.Id);
  groups[index].Tags = newTags;
}
</script>
