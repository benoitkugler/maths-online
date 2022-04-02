<template>
  <div class="ma-2">
    <v-row style="height: 89vh">
      <v-col>
        <Editor :session_id="sessionID"></Editor>
      </v-col>
      <v-col md="auto">
        <Preview :session_id="sessionID"></Preview>
      </v-col>
    </v-row>
  </div>
</template>

<script setup lang="ts">
import { controller } from "@/controller/controller";
import { onMounted } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import Editor from "../components/editor/Editor.vue";
import Preview from "../components/editor/Preview.vue";

let sessionID = $ref("");

onMounted(async () => {
  const out = await controller.EditStartSession(null);
  if (out === undefined) {
    return;
  }
  sessionID = out.ID;
});
</script>

<style></style>
