<template>
  <v-list-item link @click="emit('clicked')" class="px-1">
    <v-row no-gutters justify="space-between">
      <v-col cols="auto" align-self="center">
        <OriginButton
          :origin="props.group.Origin"
          @update-public="(b) => emit('updatePublic', props.group.Group.Id, b)"
          @create-review="emit('createReview')"
        ></OriginButton>

        <v-btn
          class="mx-2 my-1"
          size="x-small"
          icon
          @click.stop="emit('duplicate')"
          title="Dupliquer cette question"
        >
          <v-icon icon="mdi-content-copy" color="secondary"></v-icon>
        </v-btn>
      </v-col>

      <v-col class="my-3 mx-1" style="text-align: left" align-self="center">
        {{ props.group.Group.Title }}

        <VariantList :variants="variants"></VariantList>
      </v-col>
      <v-col cols="auto" align-self="center">
        <TagIndex :tags="props.group.Tags || []"></TagIndex>
      </v-col>
    </v-row>
  </v-list-item>
</template>

<script setup lang="ts">
import type { QuestiongroupExt, TagsDB } from "@/controller/api_gen";
import type { VariantG } from "@/controller/editor";
import { computed } from "vue";
import OriginButton from "../../OriginButton.vue";
import VariantList from "../VariantList.vue";
import TagIndex from "../TagIndex.vue";

interface Props {
  group: QuestiongroupExt;
  allTags: TagsDB;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "clicked"): void;
  (e: "duplicate"): void;
  (e: "updatePublic", questiongroupID: number, isPublic: boolean): void;
  (e: "createReview"): void;
}>();

const variants = computed<VariantG[]>(() => {
  return props.group.Variants || [];
});
</script>

<style scoped></style>
