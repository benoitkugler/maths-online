<template>
  <v-text-field
    autofocus
    @focus="onFocus"
    placeholder="JJ/MM/AAAA"
    variant="outlined"
    density="compact"
    :model-value="userInput"
    @update:model-value="onType"
    hide-details
    :label="props.label"
  >
  </v-text-field>
</template>

<script setup lang="ts">
import type { Date_ } from "@/controller/api_gen";
import { formatDate } from "@/controller/utils";
import { onMounted, ref, watch } from "vue";
interface Props {
  modelValue: Date_;
  label?: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:modelValue", modelValue: Date_): void;
}>();

const userInput = ref("");

onMounted(parsePropsDate);

watch(props, () => {
  parsePropsDate();
});

function onType(s: string) {
  if (s.length == 2) {
    userInput.value = s + "/";
  } else if (s.length == 5) {
    userInput.value = s + "/";
  } else {
    userInput.value = s;
  }
  const parsed = tryParseDate();
  if (parsed != null) {
    emit("update:modelValue", parsed);
  }
}

function tryParseDate() {
  if (userInput.value.length != 10) return null;
  const day = Number(userInput.value.substring(0, 2));
  const month = Number(userInput.value.substring(3, 5));
  const year = Number(userInput.value.substring(6, 10));
  if (day >= 1 && day <= 31 && month >= 1 && month <= 12 && year >= 1900) {
    return new Date(Date.UTC(year, month - 1, day)).toISOString() as Date_;
  }
  return null;
}

function parsePropsDate() {
  userInput.value = formatDate(props.modelValue);
}

function onFocus(ev: FocusEvent) {
  (ev.target as HTMLInputElement).select();
}
</script>
