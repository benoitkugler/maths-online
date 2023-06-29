<template>
  <v-row class="my-1 px-1">
    <v-col cols="12">
      <DateField
        label="Date de rendu"
        :model-value="date"
        @update:model-value="updateDate"
      ></DateField>
    </v-col>
    <v-col cols="12" align-self="center">
      <v-select
        variant="outlined"
        density="compact"
        label="Heure de rendu"
        :items="hours"
        hide-details
        v-model="hour"
        @update:model-value="updateHour"
      >
      </v-select>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { Date_, Time } from "@/controller/api_gen";
import { computed } from "vue";
import { $ref } from "vue/macros";
import DateField from "../DateField.vue";

interface Props {
  modelValue: Time;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: Time): void;
}>();

const date = computed(() => props.modelValue as unknown as Date_);

const hours = Array.from({ length: 24 }, (_, i) => ({
  title: i.toString().padStart(2, "0") + "h",
  value: i,
}));

let hour = $ref(hourFromTime(props.modelValue));

function hourFromTime(time: Time) {
  const d = new Date(time);
  return d.getHours();
}

function updateDate(date: Date_) {
  const d = new Date(date);
  d.setHours(hour);
  emit("update:model-value", d.toISOString() as Time);
}

function updateHour(hour: number) {
  const d = new Date(props.modelValue);
  console.log(d.toISOString(), hour);
  d.setHours(hour);
  console.log(d.toISOString());

  emit("update:model-value", d.toISOString() as Time);
}
</script>

<style scoped></style>
