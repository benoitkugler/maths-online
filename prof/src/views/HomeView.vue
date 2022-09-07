<template>
  <v-container class="pb-1 fill-height">
    <v-row class="fill-height" justify="space-between">
      <v-col sm="12" lg="4" align-self="center">
        <v-card>
          <v-card-title class="bg-secondary rounded pa-4"
            >Bienvenue !</v-card-title
          >
          <v-card-subtitle class="my-2"
            >Démarrer en choisissant une activité.</v-card-subtitle
          >
          <v-card-text>
            <v-row justify="center">
              <v-col>
                <v-btn
                  link
                  block
                  :to="{ name: 'classrooms' }"
                  color="purple-lighten-4"
                >
                  Gérer mes classes et mes élèves
                </v-btn>
              </v-col>
              <v-col>
                <v-btn
                  link
                  block
                  color="teal-lighten-4"
                  :to="{ name: 'editor-question' }"
                >
                  Créer ou modifier une question
                </v-btn>
              </v-col>
              <!-- <v-col>
                <v-btn
                  link
                  block
                  color="teal-lighten-4"
                  :to="{ name: 'editor-exercice' }"
                >
                  Créer ou modifier un exercice
                </v-btn>
              </v-col> -->
              <v-col>
                <v-btn
                  link
                  block
                  :to="{ name: 'trivial' }"
                  color="pink-lighten-3"
                >
                  Lancer une partie de Triv'Maths
                </v-btn>
              </v-col>
              <v-col>
                <v-btn
                  link
                  block
                  :to="{ name: 'homework' }"
                  color="pink-lighten-3"
                >
                  Lancer un travail à la maison
                </v-btn>
              </v-col>
              <v-col>
                <v-btn
                  link
                  block
                  href="https://education.alwaysdata.net/test-eleve"
                  color="primary"
                >
                  Tester l'application élève
                </v-btn>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col sm="12" lg="7" align-self="center">
        <v-img :src="cahierImg" class="rounded"></v-img>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import cahierImg from "@/assets/cahier-stylo-calc.png";
</script>
