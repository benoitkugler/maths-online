<template>
  <v-card class="bg-grey-lighten-4">
    <v-card-text>
      <v-textarea
        color="primary"
        variant="outlined"
        label="Modifier votre commentaire"
        hide-details
        v-model="comment"
      ></v-textarea>
    </v-card-text>
    <v-card-actions>
      <v-btn color="warning" @click="emit('back')"> Annuler </v-btn>
      <v-spacer></v-spacer>
      <v-btn
        color="primary-darken-1"
        variant="outlined"
        @click="emit('send', comment)"
      >
        Modifier
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import { $ref } from "vue/macros";

interface Props {
  initialMessage: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "back"): void;
  (e: "send", comment: string): void;
}>();

let comment = $ref(props.initialMessage);
</script>
