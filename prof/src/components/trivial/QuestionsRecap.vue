<template>
  <v-card title="Configuration des questions">
    <v-card-text>
      <v-row>
        <v-col>Difficultés</v-col>
        <v-col cols="auto">
          {{ formatDifficulties(props.config.Config) }}
        </v-col>
      </v-row>

      <v-row v-if="commonT.length" class="mb-2" justify="center">
        <v-col>Etiquettes partagées</v-col>
        <v-col cols="auto">
          <TagChip
            v-for="(tag, index) in commonT"
            :key="index"
            :tag="tag"
          ></TagChip>
        </v-col>
      </v-row>

      <i v-if="props.config.Config.Questions.Tags.every(v => !v?.length)">
        Aucun question configurée
      </i>
      <div v-else>
        <v-list-item
          v-for="(nbQuestions, index) in props.config.NbQuestionsByCategories"
          :key="index"
          rounded
          :style="{
            'border-color': colorsPerCategorie[index],
            borderWidth: '2px',
            borderStyle: 'solid'
          }"
          class="my-2"
        >
          <v-row>
            <v-col align-self="center">
              <div
                v-for="(intersec, j) in ownTagsFor(index)"
                :key="j"
                class="my-1"
              >
                <TagChip
                  v-for="(tag, k) in intersec"
                  :key="k"
                  :tag="tag"
                ></TagChip>
                <i v-if="!intersec.length">Aucune étiquette supplémentaire</i>
              </div>
            </v-col>
            <v-col cols="2" align-self="center" style="text-align: center">
              <v-chip :color="colorsPerCategorie[index]" variant="outlined">
                {{ nbQuestions }}
              </v-chip>
            </v-col>
          </v-row>
        </v-list-item>
      </div>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { TagSection, Trivial, TrivialExt } from "@/controller/api_gen";
import { colorsPerCategorie } from "@/controller/trivial";
import { computed } from "vue";
import TagChip from "../editor/utils/TagChip.vue";

interface Props {
  config: TrivialExt;
}

const props = defineProps<Props>();

function formatDifficulties(config: Trivial) {
  const l = config.Questions.Difficulties || [];
  if (l.length) {
    return l.join(", ");
  }
  return "Toutes difficultés";
}

const commonT = computed(() => {
  const allUnions: TagSection[][] = [];
  props.config.Config.Questions.Tags.forEach(cat =>
    allUnions.push(...(cat || []).map(s => s || []))
  );
  return commonTags(allUnions);
});

/** return the list of tags shared by all the list */
function commonTags(tags: TagSection[][]) {
  const crible = new Map<string, number>();
  tags.forEach(l =>
    l.forEach(tag =>
      crible.set(
        JSON.stringify(tag),
        (crible.get(JSON.stringify(tag)) || 0) + 1
      )
    )
  );
  return Array.from(crible.entries())
    .filter(entry => entry[1] == tags.length)
    .map(entry => JSON.parse(entry[0]) as TagSection);
}

// do not return shared tags
function ownTagsFor(categorie: number) {
  const out = props.config.Config.Questions.Tags[categorie] || [];
  return out.map(l =>
    (l || []).filter(
      tag => commonT.value.findIndex(other => other.Tag == tag.Tag) == -1
    )
  );
}
</script>

<style scoped></style>
