<template>
  <v-card title="Confirmer la demande de publication">
    <v-card-text>
      Merci pour votre participation ! <br /><br />
      Le contenu officiel est vérifié par l'équipe Isyro, en tenant compte de
      l'avis de la communauté. <br />
      En continuant, vous ajouterez votre ressource à la liste des demandes de
      publications, et nous en prendrons connaissance au plus vite.
    </v-card-text>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn color="green" @click="emit('createReview')">
        Créer une demande de publication
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
// interface Props {}

// const props = defineProps<{}>();

const emit = defineEmits<{
  (e: "createReview"): void;
}>();
</script>
