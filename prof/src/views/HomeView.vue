<template>
  <v-container class="pb-1 fill-height">
    <v-row class="fill-height" justify="space-between">
      <v-col cols="12" md="7" lg="6" align-self="center">
        <v-card>
          <v-card-title class="bg-secondary rounded pa-4">
            Bienvenue !
          </v-card-title>
          <v-card-subtitle class="my-2">
            Démarrer en choisissant une activité.
          </v-card-subtitle>
          <v-card-text>
            <v-row>
              <v-col cols="4">
                <link-card
                  route-name="classrooms"
                  color="purple-lighten-4"
                  title="Gérer mes classes et mes élèves"
                  icon="mdi-account-school"
                ></link-card>
              </v-col>
              <v-col cols="4">
                <link-card
                  route-name="editor-question"
                  color="teal-lighten-4"
                  title="Créer ou modifier une question"
                  icon="mdi-file-edit"
                ></link-card>
              </v-col>
              <v-col cols="4">
                <link-card
                  route-name="editor-exercice"
                  color="teal-lighten-4"
                  title="Créer ou modifier un exercice"
                  icon="mdi-file-edit"
                ></link-card>
              </v-col>
            </v-row>
            <v-row>
              <v-col align-self="center" cols="4">
                <link-card
                  route-name="trivial"
                  color="pink-lighten-3"
                  title="Lancer une partie Isy'Triv"
                  icon="mdi-cog"
                ></link-card>
              </v-col>
              <v-col align-self="center" cols="4">
                <link-card
                  route-name="homework"
                  color="pink-lighten-3"
                  title="Lancer un travail à la maison"
                  icon="mdi-cog"
                ></link-card>
              </v-col>
              <v-col align-self="center" cols="4">
                <link-card
                  route-name="ceintures"
                  color="pink-lighten-3"
                  title="Accéder au parcours des ceintures de calcul"
                  icon="mdi-cog"
                ></link-card>
              </v-col>
            </v-row>

            <v-row>
              <v-col cols="4">
                <link-card
                  route-name="reviews"
                  color="secondary"
                  title="Discuter des publications en cours"
                  icon="mdi-message-draw"
                ></link-card>
              </v-col>
              <v-col cols="4"></v-col>
              <v-col cols="4">
                <link-card
                  :route-name="studentAppLink"
                  color="primary"
                  title="Tester l'application élève"
                  href
                  icon="mdi-cellphone"
                ></link-card>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col md="5" lg="5" align-self="center" class="d-none d-md-flex">
        <v-img :src="cahierImg" class="rounded"></v-img>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import cahierImg from "@/assets/cahier-stylo-calc.png";
import LinkCard from "@/components/home/LinkCard.vue";
import { PreviewMode } from "@/controller/controller";
import { computed } from "vue";

const studentAppLink = computed(() =>
  PreviewMode == "dev"
    ? "http://localhost:1323/test-eleve?mode=dev"
    : "https://isyro.fr/test-eleve"
);
</script>
