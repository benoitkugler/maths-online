<template>
  <iframe v-if="session_id" :src="src" width="350px" height="100%"></iframe>
  <div
    v-else
    class="mx-4 d-flex"
    style="
      width: 350px;
      height: 100%;
      justify-content: center;
      align-items: center;
    "
  >
    <v-progress-circular indeterminate color="secondary"></v-progress-circular>
  </div>
</template>

<script setup lang="ts">
import { controller } from "@/controller/controller";
import { computed } from "@vue/runtime-core";

const props = defineProps({
  session_id: { type: String, required: true }
});
let src = computed(() =>
  controller.getURL(`/prof-loopback-app?sessionID=${props.session_id}`)
);
</script>
