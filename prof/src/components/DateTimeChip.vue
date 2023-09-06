<template>
  <v-menu
    offset-y
    :close-on-content-click="false"
    :model-value="timeToEdit != null"
    @update:model-value="timeToEdit = null"
  >
    <template v-slot:activator="{ isActive, props: innerProps }">
      <v-chip
        v-on="{ isActive }"
        v-bind="innerProps"
        style="text-align: right"
        class="ml-1"
        color="primary-darken-1"
        variant="outlined"
        @click="timeToEdit = props.modelValue"
      >
        {{ props.prefix || "" }} {{ formatTime(props.modelValue, true, true) }}
      </v-chip>
    </template>
    <v-card width="300px" :subtitle="props.title">
      <v-card-text class="pb-0" v-if="timeToEdit != null">
        <TimeField v-model="timeToEdit"></TimeField>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          color="success"
          @click="
            emit('update:modelValue', timeToEdit!);
            timeToEdit = null;
          "
          >Enregistrer</v-btn
        >
      </v-card-actions>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import type { Time } from "@/controller/api_gen";
import { formatTime } from "@/controller/utils";
import TimeField from "./homework/TimeField.vue";
import { ref } from "vue";
interface Props {
  modelValue: Time;
  title: string;
  prefix?: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:modelValue", modelValue: Time): void;
}>();

const timeToEdit = ref<Time | null>(null);
</script>
