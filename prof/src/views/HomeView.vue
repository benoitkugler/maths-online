<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card color="secondary">
          <v-card-title>Bienvenue !</v-card-title>
          <v-card-subtitle
            >Démarrer en choisissant une activité.</v-card-subtitle
          >
          <v-card-text>
            <v-row
              ><v-col>
                <v-btn link :to="{ name: 'trivial' }" color="primary">
                  Triv'Maths
                </v-btn>
              </v-col>
              <v-col>
                <v-btn link :to="{ name: 'editor' }" color="primary">
                  Editeur de question
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

  components: {}
});
</script>
