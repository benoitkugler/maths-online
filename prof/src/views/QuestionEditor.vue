<template>
  <div class="ma-2">
    <v-row>
      <v-col>
        <QuestiongroupPannel
          v-if="viewKind == 'editor'"
          :group="currentGroup!"
          :variants="currentVariants"
          :all-tags="allKnownTags"
          @back="backToList"
          @preview="(qu) => preview?.showQuestion(qu)"
        ></QuestiongroupPannel>
        <keep-alive>
          <QuestiongroupList
            v-if="viewKind == 'questions'"
            :tags="allKnownTags"
            @edit="editQuestion"
          ></QuestiongroupList>
        </keep-alive>
      </v-col>
      <v-col cols="auto">
        <keep-alive>
          <ClientPreview ref="preview"></ClientPreview>
        </keep-alive>
      </v-col>
    </v-row>
  </div>
</template>

<script setup lang="ts">
import type {
  LoopbackShowQuestion,
  Question,
  QuestiongroupExt,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { LoopbackServerEventKind } from "@/controller/loopback_gen";
import { onMounted } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import ClientPreview from "../components/editor/ClientPreview.vue";
import QuestiongroupList from "../components/editor/questions/QuestiongroupList.vue";
import QuestiongroupPannel from "../components/editor/questions/QuestiongroupPannel.vue";

let allKnownTags = $ref<string[]>([]);
let preview = $ref<InstanceType<typeof ClientPreview> | null>(null);

let viewKind: "questions" | "editor" = $ref("questions");
let currentGroup: QuestiongroupExt | null = $ref(null);
let currentVariants: Question[] = $ref([]);

onMounted(async () => {
  fetchTags();
});

async function fetchTags() {
  const tags = await controller.EditorGetTags();
  allKnownTags = tags || [];
}

function backToList() {
  currentGroup = null;
  currentVariants = [];

  fetchTags(); // required since the tags may have changed

  viewKind = "questions";
  preview?.pause();
}

function editQuestion(group: QuestiongroupExt, variants: Question[]) {
  currentGroup = group;
  currentVariants = variants;
  viewKind = "editor";
}
</script>

<style></style>
