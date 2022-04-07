<template>
  <v-list style="height: 70vh">
    <v-list-subheader>Enoncés</v-list-subheader>
    <v-list-item dense class="py-0" @click="emit('add', BlockKind.TextBlock)">
      Texte
    </v-list-item>
    <v-list-item
      dense
      class="py-0"
      @click="emit('add', BlockKind.FormulaBlock)"
    >
      Formule
    </v-list-item>
    <v-list-item dense class="py-0" @click="emit('add', BlockKind.FigureBlock)">
      Figure
    </v-list-item>
    <v-list-item dense class="py-0"> Fonction </v-list-item>
    <v-list-item dense class="py-0"> Fonction par variations </v-list-item>
    <v-list-item dense class="py-0"> Tableau de signes </v-list-item>
    <v-list-item dense class="py-0"> Tableau de variations </v-list-item>
    <v-list-item dense class="py-0"> Tableau </v-list-item>
    <v-divider></v-divider>
    <v-list-subheader>Champs de réponse</v-list-subheader>
    <v-list-item dense class="py-0"> Nombre </v-list-item>
    <v-list-item dense class="py-0"> Expression </v-list-item>
    <v-list-item dense class="py-0"> QCM </v-list-item>
  </v-list>
</template>

<script setup lang="ts">
import { BlockKind } from "@/controller/exercice_gen";

const emit = defineEmits<{
  (e: "add", kind: BlockKind): void;
}>();
</script>
