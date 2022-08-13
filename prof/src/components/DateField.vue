<template>
  <v-input hide-details>
    <v-row>
      <v-col>
        <small class="text-grey-darken-3 ml-2">
          {{ props.label }}
        </small>
      </v-col>
    </v-row>
    <v-row>
      <date-field-raw
        :model-value="props.modelValue"
        @update:model-value="(v) => emit('update:modelValue', v)"
      ></date-field-raw>
    </v-row>
  </v-input>
</template>

<script setup lang="ts">
import type { Date_ } from "@/controller/api_gen";
import DateFieldRaw from "./DateFieldRaw.vue";
interface Props {
  label: string;
  modelValue: Date_;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:modelValue", modelValue: Date_): void;
}>();
</script>
