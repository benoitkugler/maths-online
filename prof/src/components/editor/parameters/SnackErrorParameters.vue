<template>
  <v-snackbar
    :model-value="showErrorParameters"
    @update:model-value="emit('close')"
    color="warning"
  >
    <v-row v-if="props.error != null" style="width: 600px">
      <v-col>
        <v-row no-gutters>
          <v-col>
            Erreur dans la définition <b>{{ props.error.Origin }}</b> :
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <div>{{ props.error.Details }}</div>
          </v-col>
        </v-row>
      </v-col>
      <v-col cols="2" align-self="center" style="text-align: right">
        <v-btn icon size="x-small" flat @click="emit('close')">
          <v-icon icon="mdi-close" color="warning"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
  </v-snackbar>
</template>

<script setup lang="ts">
import type { ErrParameters } from "@/controller/api_gen";
import { computed } from "vue";

interface Props {
  error: ErrParameters | null;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

const showErrorParameters = computed(() => props.error != null);
</script>

<style scoped></style>
