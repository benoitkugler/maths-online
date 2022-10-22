<template>
  <v-row class="my-2">
    <v-col>
      <v-combobox
        variant="outlined"
        density="compact"
        multiple
        label="Taille des groupes"
        chips
        :model-value="props.modelValue"
        @update:model-value="onUpdate"
        @update:search-input="onSearch"
        :hide-no-data="true"
        :hide-selected="true"
        clearable
        append-inner-icon=""
        hint="Entrer les tailles des groupes séparées par des virgules. Ajouter a fois le nombre b avec a*b. Ex: 2, 4, 3*5"
        ref="field"
      >
      </v-combobox>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { $ref } from "vue/macros";
import type { VCombobox } from "vuetify/lib/components";

interface Props {
  modelValue: number[];
}

const emit = defineEmits<{
  (e: "update:model-value", v: number[]): void;
}>();

const props = defineProps<Props>();

let field = $ref<InstanceType<typeof VCombobox> | null>(null);

// onUpdate is called on delete or when typing Enter
function onUpdate(v: unknown) {
  console.log(v);

  const final: number[] = [];
  (v as (string | number)[]).forEach((item) => {
    if (typeof item == "string") {
      final.push(...parseEntry(item));
    } else {
      final.push(item);
    }
  });

  console.log(final);

  emit("update:model-value", final);
}

function onSearch(s: string) {
  s = s.trim();
  if (!s.endsWith(",")) return; // wait for separator
  s = s.substring(0, s.length - 1);
  const numbers = parseEntry(s);
  if (numbers.length) {
    emit("update:model-value", props.modelValue.concat(...numbers));
    field.search = "";
  }
}

function parseEntry(s: string) {
  const chunks = s.split("*");
  if (chunks.length == 2) {
    const a = Number(chunks[0].trim());
    const b = Number(chunks[1].trim());
    if (isNaN(a) || isNaN(b)) {
      return [];
    }
    return Array.from({ length: a }, () => b);
  }

  const v = Number(s);
  if (isNaN(v)) {
    return [];
  }
  return [v];
}
</script>

<style scoped></style>
