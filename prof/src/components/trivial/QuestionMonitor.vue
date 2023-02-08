<template>
  <v-card>
    <v-row>
      <v-col style="text-align: right">
        <v-btn icon flat @click="emit('close')">
          <v-icon icon="mdi-close"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-card-text style="height: 100%">
      <v-row no-gutters style="height: 100%" justify="center">
        <v-col cols="auto">
          <iframe
            :src="src"
            width="400px"
            height="100%"
            ref="iframe"
            @load="setupListener"
          ></iframe>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { QuestionContent } from "@/controller/api_gen";
import { controller, PreviewMode } from "@/controller/controller";
import { computed } from "vue";
import { $ref } from "vue/macros";

interface Props {
  question: QuestionContent;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

let iframe = $ref<HTMLIFrameElement | null>(null);

function setupListener() {
  if (!iframe?.contentWindow) return;
  window.addEventListener("message", (ev) => {
    const data = JSON.parse(ev.data);
    if (data["PREVIEW_READY"]) {
      sendEvent();
    }
  });
}

/** transfer the given payload to the flutter embedded app */
function sendEvent() {
  if (iframe == null) return;
  iframe.contentWindow?.postMessage(JSON.stringify(props.question), "*");
}

let src = computed(() =>
  controller.getURL(`/prof-preview-app?mode=${PreviewMode}`)
);
</script>

<style scoped></style>
