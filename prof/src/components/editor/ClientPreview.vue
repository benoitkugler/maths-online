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
import {
  LoopbackServerEventKind,
  type LoopbackServerEvent,
} from "@/controller/loopback_gen";
import { computed, ref } from "vue";

interface Props {
  hide?: boolean;
}

const props = defineProps<Props>();

defineExpose({ pause, preview });

const iframe = ref<HTMLIFrameElement | null>(null);

/** transfer the given payload to the flutter embedded app */
function sendEvent(previewEvent: LoopbackServerEvent) {
  if (iframe.value == null) return;
  iframe.value.contentWindow?.postMessage(JSON.stringify(previewEvent), "*");
}

function pause() {
  sendEvent({
    Kind: LoopbackServerEventKind.LoopbackPaused,
    Data: {},
  });
}

function preview(data: LoopbackServerEvent) {
  sendEvent(data);
}

const src = computed(() =>
  controller.getURL(`/prof-loopback-app?mode=${PreviewMode}`)
);
</script>
