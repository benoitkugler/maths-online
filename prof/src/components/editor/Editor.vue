<template>
  <v-card class="ma-1">
    <v-row>
      <v-col>
        <v-card-title>Editeur de question</v-card-title>
        <v-card-subtitle
          >Editer et visualiser une question pour l'élève</v-card-subtitle
        >
      </v-col>
      <v-spacer></v-spacer>
      <v-col align-self="center" style="text-align: right">
        <v-btn icon class="mx-2" @click="save" :disabled="!session_id">
          <v-icon icon="mdi-content-save"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
  </v-card>
</template>

<script setup lang="ts">
import { controller } from "@/controller/controller";

const props = defineProps({
  session_id: { type: String, required: true }
});

async function save() {
  await controller.EditSaveAndPreview({ SessionID: props.session_id || "" });
}
</script>
