<template>
  <v-col>
    <v-text-field
      label="Jour"
      variant="outlined"
      density="compact"
      :model-value="day"
      @update:model-value="onUpdateDay"
      hide-details
    ></v-text-field>
  </v-col>
  <v-col>
    <v-text-field
      label="Mois"
      variant="outlined"
      density="compact"
      v-model.num="month"
      @update:model-value="onUpdateMonth"
      hide-details
      ref="monthField"
    ></v-text-field>
  </v-col>
  <v-col>
    <v-text-field
      label="Année"
      type="number"
      variant="outlined"
      density="compact"
      v-model.num="year"
      @update:model-value="onUpdateYear"
      hide-details
      ref="yearField"
    ></v-text-field>
  </v-col>
</template>

<script setup lang="ts">
import type { Date_ } from "@/controller/api_gen";
import { formatDate } from "@/controller/utils";
import { onMounted, watch } from "vue";
import { $ref } from "vue/macros";
import type { VTextField } from "vuetify/lib/components";
interface Props {
  modelValue: Date_;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:modelValue", modelValue: Date_): void;
}>();

let day = $ref("01");
let month = $ref("01");
let year = $ref(2005);

let monthField = $ref<InstanceType<typeof VTextField> | null>(null);
let yearField = $ref<InstanceType<typeof VTextField> | null>(null);

onMounted(parseDate);

watch(props, () => {
  parseDate();
});

function parseDate() {
  const date = formatDate(props.modelValue);
  day = date.substring(0, 2);
  month = date.substring(3, 5);
  year = Number(date.substring(6, 10));
}

// emit if the current date is valid
function tryEmit() {
  const day_ = Number(day);
  const month_ = Number(month);
  if (day_ >= 1 && day_ <= 31 && month_ >= 1 && month_ <= 12 && year >= 1900) {
    emit(
      "update:modelValue",
      new Date(Date.UTC(year, month_ - 1, day_)).toISOString() as Date_
    );
  }
}

function onUpdateDay(s: string) {
  if (s == "") {
    return;
  }
  day = s;
  if (s.length >= 2 && monthField != null) {
    (monthField.$el as HTMLInputElement).querySelector("input")!.select();
  }
  tryEmit();
}
function onUpdateMonth(s: string) {
  if (s == "") {
    return;
  }
  month = s;
  if (s.length >= 2 && yearField != null) {
    (yearField.$el as HTMLInputElement).querySelector("input")!.select();
  }
  tryEmit();
}

function onUpdateYear() {
  tryEmit();
}
</script>
