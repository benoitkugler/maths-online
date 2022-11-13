<template>
  <v-card class="my-2 bg-grey-lighten-4">
    <v-card-text>
      <v-textarea
        color="primary"
        variant="outlined"
        label="Nouveau commentaire"
        placeholder="Partarger vos retours"
        hide-details
        v-model="comment"
        :disabled="props.disabled"
      ></v-textarea>
    </v-card-text>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn
        color="primary-darken-1"
        variant="outlined"
        @click="emit('send', comment)"
        :disabled="props.disabled"
      >
        Envoyer
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import { $ref } from "vue/macros";

interface Props {
  disabled: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "send", comment: string): void;
}>();

let comment = $ref("");
</script>
