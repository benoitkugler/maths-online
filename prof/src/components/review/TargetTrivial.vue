<template>
  <div>
    <v-row justify="center">
      <v-col cols="auto" align-self="center">Niveau de difficulté :</v-col>
      <v-col cols="auto" align-self="center">
        <v-row no-gutters>
          <v-col
            v-for="(tag, index) in props.data.Config.Questions.Difficulties"
            :key="index"
          >
            <tag-chip :tag="{ Tag: tag, Section: 0 }"></tag-chip>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
    <v-list>
      <categorie-row
        v-for="(categorie, index) in categories"
        :key="index"
        :index="index"
      >
        <v-row>
          <v-col align-self="center" cols="2">
            <v-list-item-subtitle>
              Catégorie {{ index + 1 }}
            </v-list-item-subtitle>
          </v-col>
          <v-col align-self="center" class="my-1">
            <v-row v-for="(inter, j) in categorie.Tags" :key="j">
              <v-col>
                <tag-list-field
                  :model-value="inter || []"
                  readonly
                  :allTags="emptyTagsDB()"
                ></tag-list-field>
              </v-col>
            </v-row>
          </v-col>
          <v-col cols="2" align-self="center">
            <v-chip
              >{{ categorie.QuestionNumber }} question{{
                categorie.QuestionNumber >= 2 ? "s" : ""
              }}</v-chip
            >
          </v-col>
        </v-row>
      </categorie-row>
    </v-list>
  </div>
</template>

<script setup lang="ts">
import type { TargetTrivial } from "@/controller/api_gen";
import { computed } from "@vue/reactivity";
import TagListField from "../editor/TagListField.vue";
import CategorieRow from "../trivial/CategorieRow.vue";
import TagChip from "../editor/utils/TagChip.vue";
import { emptyTagsDB } from "@/controller/editor";

interface Props {
  data: TargetTrivial;
}

const props = defineProps<Props>();

const categories = computed(() => {
  return props.data.Config.Questions.Tags.map((cat, index) => {
    return {
      QuestionNumber: props.data.NbQuestionsByCategories[index],
      Tags: cat,
    };
  });
});
</script>
