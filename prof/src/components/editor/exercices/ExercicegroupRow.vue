<template>
  <v-list-item link @click="emit('clicked')" class="px-1">
    <v-row no-gutters justify="space-between">
      <v-col cols="auto" align-self="center">
        <OriginButton
          :origin="props.group.Origin"
          @update-public="(b) => emit('updatePublic', props.group.Group.Id, b)"
        ></OriginButton>

        <v-btn
          class="mx-2 my-1"
          size="x-small"
          icon
          @click.stop="emit('duplicate')"
          title="Dupliquer cet exercice"
        >
          <v-icon icon="mdi-content-copy" color="secondary"></v-icon>
        </v-btn>
      </v-col>

      <v-col class="my-3 mx-1" style="text-align: left" align-self="center">
        {{ props.group.Group.Title }}

        <VariantList :variants="variants"></VariantList>
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
import { Visibility, type ExercicegroupExt } from "@/controller/api_gen";
import type { VariantG } from "@/controller/editor";
import { computed } from "vue";
import OriginButton from "../../OriginButton.vue";
import TagListField from "../TagListField.vue";
import VariantList from "../VariantList.vue";

interface Props {
  group: ExercicegroupExt;
  allTags: string[];
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "clicked"): void;
  (e: "duplicate"): void;
  (e: "updatePublic", exercicegroupID: number, isPublic: boolean): void;
  (e: "updateTags", tags: string[]): void;
}>();

const isEditable = computed(
  () => props.group.Origin.Visibility == Visibility.Personnal
);

const variants = computed<VariantG[]>(() => {
  return props.group.Variants || [];
});
</script>

<style scoped></style>
