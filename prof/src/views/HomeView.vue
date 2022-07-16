<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card>
          <v-card-title class="bg-secondary rounded">Bienvenue !</v-card-title>
          <v-card-subtitle
            >Démarrer en choisissant une activité.</v-card-subtitle
          >
          <v-card-text>
            <v-row justify="center" style="text-align: center">
              <v-col>
                <v-btn
                  link
                  :to="{ name: 'classrooms' }"
                  color="purple-lighten-4"
                >
                  Classes et élèves
                </v-btn>
              </v-col>
              <v-col>
                <v-btn link :to="{ name: 'trivial' }" color="pink-lighten-3">
                  Triv'Maths
                </v-btn>
              </v-col>
              <v-col>
                <v-btn
                  link
                  color="teal-lighten-4"
                  :to="{ name: 'editor-question' }"
                >
                  Editeur de question
                </v-btn>
              </v-col>
              <v-col>
                <v-btn link color="teal-lighten-4" :to="{ name: 'editor-exercice' }">
                  Editeur d'exercice
                </v-btn>
              </v-col>
              <v-col>
                <v-btn
                  link
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
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "HomeView",

  components: {},
});
</script>
