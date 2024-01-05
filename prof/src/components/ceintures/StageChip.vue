<template>
  <v-chip
    variant="elevated"
    :color="rankColors[props.stage.Rank]"
    class="ma-2"
    @[props.link&&`click`]="emit('click')"
  >
    {{ DomainLabels[props.stage.Domain] }}
  </v-chip>
</template>

<script setup lang="ts">
import { Stage } from "@/controller/api_gen";
import { DomainLabels } from "@/controller/api_gen";
import { rankColors } from "@/controller/utils";

interface Props {
  stage: Stage;
  link?: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "click"): void;
}>();
</script>
