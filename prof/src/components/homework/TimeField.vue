<template>
  <v-row class="my-1 px-1">
    <v-col cols="12">
      <DateField
        label="Date"
        :model-value="date"
        @update:model-value="updateDate"
      ></DateField>
    </v-col>
    <v-col cols="6" align-self="center">
      <v-select
        variant="outlined"
        density="compact"
        label="Heure"
        :items="hours"
        hide-details
        v-model="hour"
        @update:model-value="updateHour"
      >
      </v-select>
    </v-col>
    <v-col cols="6" align-self="center">
      <v-select
        variant="outlined"
        density="compact"
        label="Minute"
        :items="minutes"
        hide-details
        v-model="minute"
        @update:model-value="updateMinute"
      >
      </v-select>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { Date_, Time } from "@/controller/api_gen";
import { ref, computed } from "vue";
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
const minutes = Array.from({ length: 12 }, (_, i) => ({
  title: (i * 5).toString().padStart(2, "0"),
  value: i * 5,
}));

const hour = ref(hourFromTime(props.modelValue));
const minute = ref(minuteFromTime(props.modelValue));

function hourFromTime(time: Time) {
  const d = new Date(time);
  return d.getHours();
}
function minuteFromTime(time: Time) {
  const d = new Date(time);
  return d.getMinutes();
}

function updateDate(date: Date_) {
  const d = new Date(date);
  d.setHours(hour.value);
  d.setMinutes(minute.value);
  emit("update:model-value", d.toISOString() as Time);
}

function updateHour(hour: number) {
  const d = new Date(props.modelValue);
  d.setHours(hour);
  emit("update:model-value", d.toISOString() as Time);
}
function updateMinute(minute: number) {
  const d = new Date(props.modelValue);
  d.setMinutes(minute);
  emit("update:model-value", d.toISOString() as Time);
}
</script>

<style scoped></style>
