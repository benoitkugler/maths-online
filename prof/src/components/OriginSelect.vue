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
          @update:model-value="(o) => emit('update:origin', o!)"
          hide-details
        >
          <v-radio :value="OriginKind.All">
            <template v-slot:label>
              <v-chip
                class="w-100"
                label
                variant="elevated"
                color="grey-lighten-3"
                >Tous</v-chip
              >
            </template>
          </v-radio>
          <v-radio :value="OriginKind.OnlyPersonnal">
            <template v-slot:label>
              <v-chip
                class="w-100"
                label
                variant="elevated"
                :color="ColorPersonnal"
                >Entrées personnelles</v-chip
              >
            </template></v-radio
          >
          <v-radio :value="OriginKind.OnlyAdmin">
            <template v-slot:label>
              <v-chip class="w-100" label variant="elevated" :color="ColorAdmin"
                >Entrées officielles</v-chip
              >
            </template></v-radio
          >
        </v-radio-group>
      </v-card-text>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import { OriginKind } from "@/controller/api_gen";
import { ColorAdmin, ColorPersonnal } from "@/controller/utils";

interface Props {
  origin: OriginKind;
}
const props = defineProps<Props>();
const emit = defineEmits<{
  (e: "update:origin", origin: OriginKind): void;
}>();
</script>
