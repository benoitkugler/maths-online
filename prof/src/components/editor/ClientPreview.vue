<template>
  <iframe
    :src="src"
    height="100%"
    ref="iframe"
    class="rounded"
    :style="
      props.hide
        ? 'position: absolute;width:0;height:0;border:0;'
        : 'width:350px'
    "
  ></iframe>
</template>

<script setup lang="ts">
import { controller, PreviewMode } from "@/controller/controller";
import { computed, defineExpose } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import {
  LoopbackServerEventKind,
  type LoopbackServerEvent,
} from "@/controller/loopback_gen";
import type {
  LoopbackShowExercice,
  LoopbackShowQuestion,
} from "@/controller/api_gen";

interface Props {
  hide?: boolean;
}

const props = defineProps<Props>();

defineExpose({ pause, showQuestion, showExercice });

let iframe = $ref<HTMLIFrameElement | null>(null);

/** transfer the given payload to the flutter embedded app */
function sendEvent(previewEvent: LoopbackServerEvent) {
  if (iframe == null) return;
  iframe.contentWindow?.postMessage(JSON.stringify(previewEvent), "*");
}

function pause() {
  sendEvent({
    Kind: LoopbackServerEventKind.LoopbackPaused,
    Data: {},
  });
}

function showQuestion(qu: LoopbackShowQuestion) {
  sendEvent({
    Kind: LoopbackServerEventKind.LoopbackShowQuestion,
    Data: qu,
  });
}

function showExercice(qu: LoopbackShowExercice) {
  sendEvent({
    Kind: LoopbackServerEventKind.LoopbackShowExercice,
    Data: qu,
  });
}

let src = computed(() =>
  controller.getURL(`/prof-loopback-app?mode=${PreviewMode}`)
);
</script>
