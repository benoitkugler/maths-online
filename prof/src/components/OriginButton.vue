<template>
  <v-menu offset-y close-on-content-click>
    <template v-slot:activator="{ isActive, props }">
      <v-btn
        v-on="{ isActive }"
        v-bind="props"
        class="mx-1"
        size="x-small"
        icon
        title="Options de partage"
        @click.stop
        :color="isPersonnalAndShared ? 'blue' : undefined"
      >
        <v-icon icon="mdi-share-variant" size="small"></v-icon>
      </v-btn>
    </template>
    <OriginCard
      :origin="props.origin"
      @update="(b) => emit('updatePublic', b)"
    ></OriginCard>
  </v-menu>
</template>

<script setup lang="ts">
import { Visibility, type Origin } from "@/controller/api_gen";
import { computed } from "vue";
import OriginCard from "./OriginCard.vue";

interface Props {
  origin: Origin;
}
const props = defineProps<Props>();
const emit = defineEmits<{
  (e: "updatePublic", isPublic: boolean): void;
}>();

const isPersonnalAndShared = computed(
  () => props.origin.Visibility == Visibility.Personnal && props.origin.IsPublic
);
</script>
