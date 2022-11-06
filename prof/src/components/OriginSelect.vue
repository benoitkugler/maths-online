<template>
  <v-menu>
    <template v-slot:activator="{ isActive, props }">
      <v-btn
        v-on="{ isActive }"
        v-bind="props"
        class="mx-1"
        size="x-small"
        icon
        title="Selectionner l'auteur..."
        @click.stop
      >
        <v-icon icon="mdi-share-variant" size="small"></v-icon>
      </v-btn>
    </template>
    <v-card density="compact">
      <v-card-text>
        <v-radio-group
          :model-value="props.origin"
          hide-details
          @update:model-value="(o) => emit('update:origin', o)"
        >
          <v-radio label="Tout afficher" :value="OriginKind.All"></v-radio>
          <v-radio
            label="Entrées personnelles"
            :value="OriginKind.OnlyPersonnal"
          ></v-radio>
          <v-radio
            label="Entrées officielles"
            :value="OriginKind.OnlyAdmin"
          ></v-radio>
        </v-radio-group>
      </v-card-text>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import { OriginKind } from "@/controller/api_gen";

interface Props {
  origin: OriginKind;
}
const props = defineProps<Props>();
const emit = defineEmits<{
  (e: "update:origin", origin: OriginKind): void;
}>();
</script>
