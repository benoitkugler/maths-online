<template>
  <v-dialog v-model="showErrVarsDetails" max-width="800">
    <v-card subtitle="Valeurs des paramètres aléatoires">
      <v-card-text>
        L'erreur est rencontrée pour les valeurs suivantes :
        <v-list>
          <v-list-item v-for="(entry, index) in errVars" :key="index">
            <v-row no-gutters>
              <v-col>
                {{ entry[0] }}
              </v-col>
              <v-col class="text-grey">
                {{ entry[1] }}
              </v-col>
            </v-row>
          </v-list-item>
        </v-list>
      </v-card-text>
    </v-card>
  </v-dialog>

  <v-snackbar
    :model-value="showError"
    @update:model-value="emit('close')"
    :timeout="timeout"
    color="warning"
  >
    <v-row v-if="props.error != null">
      <v-col>
        <v-row no-gutters>
          <v-col>
            <b
              >Erreur dans le contenu de la
              {{ props.isCorrection ? "correction" : "question" }}</b
            >
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <div>
              <i v-html="props.error.Error"></i>
            </div>
          </v-col>
        </v-row>
      </v-col>
      <v-col
        v-if="errVars.length > 0"
        cols="3"
        align-self="center"
        class="px-1"
      >
        <v-btn variant="outlined" @click="onShowDetails"> Détails </v-btn>
      </v-col>
      <v-col
        cols="2"
        align-self="center"
        style="text-align: right"
        class="px-1"
      >
        <v-btn icon size="x-small" @click="emit('close')">
          <v-icon icon="mdi-close" color="warning"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
  </v-snackbar>
</template>

<script setup lang="ts">
import type { errEnonce } from "@/controller/api_gen";
import { watch } from "vue";
import { ref, computed } from "vue";

interface Props {
  error: errEnonce | null;
  isCorrection: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

const showError = computed(() => props.error != null);

const errVars = computed(() => {
  const out = Object.entries(props.error?.Vars || {});
  out.sort((a, b) => a[0].localeCompare(b[0]));
  return out;
});

const showErrVarsDetails = ref(false);
let timeout = 5000;

watch(props, () => (timeout = 5000));

function onShowDetails() {
  showErrVarsDetails.value = true;
  // remove the snackbar timeout so that the error stays on screen
  timeout = -1;
}
</script>

<style scoped></style>
