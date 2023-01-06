<template>
  <v-card
    title="Champ de réponse Expression"
    subtitle="La réponse est une expression mathématique dans l'un des formats suivants."
  >
    <v-card-text>
      <v-list>
        <v-list-item>
          <v-row>
            <v-col align-self="center" cols="3">
              <v-list-item-title>Expression simple</v-list-item-title>
            </v-col>
            <v-col align-self="center" cols="6" class="text-grey"
              >Une expression correspondant à un nombre réel.</v-col
            >
            <v-col align-self="center" style="text-align: right">
              2x - y + sin(u)
            </v-col>
          </v-row>
          <v-row>
            <v-col align-self="center" cols="3">
              <v-list-item-title>Vecteur</v-list-item-title>
            </v-col>
            <v-col align-self="center" cols="6" class="text-grey">
              Un vecteur délimité par des paranthèses et des point-virgules.
            </v-col>
            <v-col align-self="center" style="text-align: right">
              (2x -y ; 4 + z; t^2)
            </v-col>
          </v-row>
          <v-row>
            <v-col align-self="center" cols="3">
              <v-list-item-title>Ensemble</v-list-item-title>
            </v-col>
            <v-col align-self="center" cols="6" class="text-grey">
              Un ensemble délimité par des accolades et des point-virgules.
            </v-col>
            <v-col align-self="center" style="text-align: right">
              { 2x -y ; 4 + z; t^2 }
            </v-col>
          </v-row>
          <v-row>
            <v-col align-self="center" cols="3">
              <v-list-item-title>Intervalle</v-list-item-title>
            </v-col>
            <v-col align-self="center" cols="6" class="text-grey">
              Un intervalle délimité par des crochets et un point-virgule.
            </v-col>
            <v-col align-self="center" style="text-align: right">
              ] -inf; 4 + z ]
            </v-col>
          </v-row>
        </v-list-item>
      </v-list>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
// const emit = defineEmits<{
//   (e: "close"): void;
// }>();

// interface ContainerProps {

// }

// const props = defineProps<ContainerProps>();
</script>

<style scoped></style>
