<template>
  <v-card title="Description" subtitle="Note interne affichée aux professeurs.">
    <v-card-text>
      <v-textarea
        v-model="description"
        placeholder="Commentaire..."
        :readonly="props.readonly"
      ></v-textarea>
      <v-row no-gutters>
        <v-spacer></v-spacer>
        <v-col cols="auto">
          <v-btn @click="emit('save', description)" :disabled="props.readonly"
            >Enregistrer</v-btn
          >
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { ref } from "vue";

interface Props {
  description: string;
  readonly: boolean;
}

const props = defineProps<Props>();

const description = ref(props.description);

const emit = defineEmits<{
  (e: "save", description: string): void;
}>();
</script>

<style scoped></style>
