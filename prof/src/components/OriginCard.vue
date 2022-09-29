<template>
  <v-card :color="color">
    <v-card-text>
      <span v-html="visibilityLabel"></span>
      <div v-if="isPersonnal && props.origin.AllowPublish">
        <v-switch
          label="Partager à la communauté"
          hide-details
          :model-value="props.origin.IsPublic"
          :color="props.origin.IsPublic ? 'blue' : undefined"
          @update:model-value="(b:boolean) => emit('update', b)"
        ></v-switch>
      </div>
      <div v-else-if="isPersonnal">
        <v-btn>Demander à rendre public : à venir...</v-btn>
      </div>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import {
  Visibility,
  VisibilityLabels,
  type Origin,
} from "@/controller/api_gen";
import { visiblityColors } from "@/controller/editor";
import { computed } from "vue";
interface Props {
  origin: Origin;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update", isPublic: boolean): void;
}>();

const color = computed(() => visiblityColors[props.origin.Visibility]);
const isPersonnal = computed(
  () => props.origin.Visibility == Visibility.Personnal
);
const visibilityLabel = computed(() => {
  switch (props.origin.Visibility) {
    case Visibility.Admin:
      return "<b>" + VisibilityLabels[props.origin.Visibility] + "</b>";
    case Visibility.Personnal:
      return VisibilityLabels[props.origin.Visibility];
    default:
      throw new Error("");
  }
});
</script>
