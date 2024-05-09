<template>
  <v-card :color="rankColors[props.stage.Rank]" link @click="emit('click')">
    <v-card-text class="text-center pa-1">
      <v-tooltip
        v-if="props.header.HasTODO"
        text="Certaines questions sont en cours d'édition."
      >
        <template v-slot:activator="{ isActive, props }">
          <v-icon v-on="{ isActive }" v-bind="props" class="mr-2"
            >mdi-progress-alert</v-icon
          >
        </template>
      </v-tooltip>

      <v-tooltip>
        <template v-slot:activator="{ isActive, props: innerProps }">
          <span v-on="{ isActive }" v-bind="innerProps">
            {{ props.header.Questions?.length || 0 }} qu.
          </span>
        </template>

        <i v-if="!props.header.Questions?.length">"Aucune question."</i>
        <div v-for="(qu, index) in props.header.Questions" :key="index">
          {{ qu.Title }}
        </div>
      </v-tooltip>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { Stage, StageHeader } from "@/controller/api_gen";
import { rankColors } from "@/controller/utils";

interface Props {
  stage: Stage;
  header: StageHeader;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "click"): void;
}>();
</script>
