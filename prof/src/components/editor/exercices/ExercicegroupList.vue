<template>
  <v-card class="pt-2 pb-0">
    <v-row>
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
      <v-row>
        <v-col>
          <v-text-field
            label="Rechercher"
            hint="Rechercher un exercice par nom."
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

            <div v-for="exerciceGroup in groups" :key="exerciceGroup.Group.Id">
              <exercicegroup-row
                :group="exerciceGroup"
                :all-tags="props.tags"
                @clicked="startEdit(exerciceGroup)"
                @duplicate="duplicate(exerciceGroup)"
                @update-public="updatePublic"
                @create-review="createReview(exerciceGroup.Group)"
                @update-tags="
                  (tags) => updateGroupTags(exerciceGroup.Group, tags)
                "
              ></exercicegroup-row>
            </div>
          </v-list>
          <div class="my-2">
            {{ groups.length }} / {{ serverNbGroups }} exercices affichés. ({{
              displayedNbExercices
            }}
            / {{ serverNbExercices }}
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
  ReviewKind,
  type Exercicegroup,
  type ExercicegroupExt,
} from "@/controller/api_gen";
import { controller, IsDev } from "@/controller/controller";
import { computed, onActivated, onMounted } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import OriginSelect from "../../OriginSelect.vue";
import ExercicegroupRow from "./ExercicegroupRow.vue";

interface Props {
  tags: string[]; // queried once for all
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "edit", group: ExercicegroupExt): void;
}>();

let groups = $ref<ExercicegroupExt[]>([]);
let serverNbGroups = $ref(0);
let serverNbExercices = $ref(0);
const displayedNbExercices = computed(() => {
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

onMounted(fetchExercices);
onActivated(fetchExercices);

function updateQuerySearch() {
  const debounceDelay = 200;
  // cancel pending call
  clearTimeout(timerId);

  // delay new call 500ms
  timerId = setTimeout(() => {
    fetchExercices();
  }, debounceDelay);
}

async function updateQueryTags() {
  await fetchExercices();
}

async function updateQueryOrigin(o: OriginKind) {
  queryOrigin = o;
  await fetchExercices();
}

async function duplicate(group: ExercicegroupExt) {
  const ok = await controller.EditorDuplicateExercicegroup({
    id: group.Group.Id,
  });
  if (!ok) return;
  await fetchExercices();
}

async function fetchExercices() {
  const result = await controller.EditorSearchExercices({
    TitleQuery: querySearch,
    Tags: queryTags,
    Origin: queryOrigin,
  });
  if (result == undefined) {
    return;
  }
  groups = result.Groups || [];
  serverNbGroups = result.NbGroups;
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
    Public: isPublic,
  });
  if (res === undefined) {
    return;
  }

  const index = groups.findIndex((gr) => gr.Group.Id == id);
  groups[index].Origin.IsPublic = isPublic;
}

async function createReview(ex: Exercicegroup) {
  const res = await controller.ReviewCreate({
    Kind: ReviewKind.KExercice,
    Id: ex.Id,
  });
  if (res == undefined) return;
  // TODO; maybe go to review ?
}

async function updateGroupTags(group: Exercicegroup, newTags: string[]) {
  const rep = await controller.EditorUpdateExerciceTags({
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
