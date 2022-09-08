<template>
  <v-list-item @click="emit('clicked')" class="px-1">
    <v-row no-gutters justify="space-between">
      <v-col cols="auto" align-self="center">
        <OriginButton
          :origin="props.group.Origin"
          @update-public="b => emit('updatePublic', props.group.Group.Id, b)"
        ></OriginButton>
      </v-col>

      <v-col class="my-3 mx-1" style="text-align: left" align-self="center">
        {{ props.group.Group.Title }}

        <v-tooltip bottom>
          <template v-slot:activator="{ props }">
            <v-badge
              @click.stop="() => {}"
              v-bind="props"
              color="info"
              inline
              :content="group.Variants?.length || 0"
            ></v-badge>
          </template>
          <v-card max-width="500" subtitle="Variantes">
            <v-card-text class="pa-0">
              <v-list style="max-height: 70vh" class="py-0">
                <v-list-item
                  density="compact"
                  rounded
                  class="my-1"
                  v-for="question in group.Variants"
                  :key="question.Id"
                >
                  <v-row no-gutters>
                    <v-col align-self="center"> ({{ question.Id }}) </v-col>
                    <v-col align-self="center" class="my-4 px-3" cols="auto">
                      {{ question.Subtitle || "..." }}
                    </v-col>
                    <v-col align-self="center" style="text-align: right">
                      <TagChip
                        v-if="question.Difficulty"
                        :tag="question.Difficulty"
                      ></TagChip>
                      <v-chip v-else size="small" label title="Difficulté"
                        >Aucune</v-chip
                      >
                    </v-col>
                  </v-row>
                </v-list-item>
              </v-list>
            </v-card-text>
          </v-card>
        </v-tooltip>
      </v-col>
      <v-col cols="5" align-self="center">
        <TagListField
          :readonly="!isEditable"
          :model-value="props.group.Tags || []"
          @update:model-value="l => emit('updateTags', l)"
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
import TagChip from "../utils/TagChip.vue";

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
