<template>
  <v-list-item @click="emit('clicked')" class="px-1">
    <v-row no-gutters justify="space-between">
      <v-col cols="auto" align-self="center">
        <OriginButton
          :origin="props.group.Origin"
          @update-public="(b) => emit('updatePublic', props.group.Group.Id, b)"
        ></OriginButton>
      </v-col>

      <v-col class="my-3 mx-1" style="text-align: left" align-self="center">
        {{ props.group.Group.Title }}
        <v-badge
          color="info"
          inline
          :content="props.group.Questions?.length || 0"
        ></v-badge>
      </v-col>
      <v-col cols="5" align-self="center">
        <TagListField
          :readonly="!isEditable"
          :model-value="props.group.Tags || []"
          @update:model-value="(l) => emit('updateTags', l)"
          :all-tags="props.allTags"
          y-padding
        >
        </TagListField>
      </v-col>
    </v-row>
  </v-list-item>
</template>

<script setup lang="ts">
import { Visibility, type QuestiongroupExt } from "@/controller/api_gen";
import { computed } from "vue";
import OriginButton from "../../OriginButton.vue";
import TagListField from "../TagListField.vue";

interface Props {
  group: QuestiongroupExt;
  allTags: string[];
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "clicked"): void;
  (e: "updatePublic", questiongroupID: number, isPublic: boolean): void;
  (e: "updateTags", tags: string[]): void;
}>();

const isEditable = computed(
  () => props.group.Origin.Visibility == Visibility.Personnal
);
</script>

<style scoped></style>
