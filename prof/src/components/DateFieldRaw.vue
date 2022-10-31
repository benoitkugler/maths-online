<template>
  <v-col>
    <v-text-field
      @focus="onFocus"
      placeholder="JJ/MM/AAAA"
      variant="outlined"
      density="compact"
      :model-value="userInput"
      @update:model-value="onType"
      hide-details
    >
    </v-text-field>
  </v-col>
</template>

<script setup lang="ts">
import type { Date_ } from "@/controller/api_gen";
import { formatDate } from "@/controller/utils";
import { onMounted, watch } from "vue";
import { $ref } from "vue/macros";
interface Props {
  modelValue: Date_;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:modelValue", modelValue: Date_): void;
}>();

let userInput = $ref("");

onMounted(parsePropsDate);

watch(props, () => {
  parsePropsDate();
});

function onType(s: string) {
  if (s.length == 2) {
    userInput = s + "/";
  } else if (s.length == 5) {
    userInput = s + "/";
  } else {
    userInput = s;
  }
  const parsed = tryParseDate();
  if (parsed != null) {
    emit("update:modelValue", parsed);
  }
}

function tryParseDate() {
  if (userInput.length != 10) return null;
  const day = Number(userInput.substring(0, 2));
  const month = Number(userInput.substring(3, 5));
  const year = Number(userInput.substring(6, 10));
  if (day >= 1 && day <= 31 && month >= 1 && month <= 12 && year >= 1900) {
    return new Date(Date.UTC(year, month - 1, day)).toISOString() as Date_;
  }
  return null;
}

function parsePropsDate() {
  userInput = formatDate(props.modelValue);
}

function onFocus(ev: FocusEvent) {
  (ev.target as HTMLInputElement).select();
}
</script>
