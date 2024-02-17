<template>
  <v-row class="my-2">
    <v-col>
      <v-combobox
        variant="outlined"
        density="compact"
        multiple
        label="Taille des groupes"
        chips
        :model-value="props.modelValue.Groups || []"
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
import type { GroupsStrategyAuto, Int } from "@/controller/api_gen";
import { ref } from "vue";
import { VCombobox } from "vuetify/lib/components/index.mjs";

interface Props {
  modelValue: GroupsStrategyAuto;
}

const emit = defineEmits<{
  (e: "update:model-value", v: GroupsStrategyAuto): void;
}>();

const props = defineProps<Props>();

const field = ref<InstanceType<typeof VCombobox> | null>(null);

// onUpdate is called on delete or when typing Enter
function onUpdate(v: unknown) {
  const final: Int[] = [];
  (v as (string | Int)[]).forEach((item) => {
    if (typeof item == "string") {
      final.push(...parseEntry(item));
    } else {
      final.push(item);
    }
  });

  emit("update:model-value", { Groups: final });
}

function onSearch(s: string) {
  s = s.trim();
  if (!s.endsWith(",")) return; // wait for separator
  s = s.substring(0, s.length - 1);
  const numbers = parseEntry(s);
  if (numbers.length) {
    const newGroups = (props.modelValue.Groups || []).concat(...numbers);
    emit("update:model-value", { Groups: newGroups });
    field.value!.search = "";
  }
}

function parseEntry(s: string) {
  let chunks = s.split("*");
  if (chunks.length == 2) {
    const a = Number(chunks[0].trim());
    const b = Number(chunks[1].trim());
    if (isNaN(a) || isNaN(b)) {
      return [];
    }
    return Array.from({ length: a }, () => b as Int);
  }

  return s
    .split(",")
    .map((c) => Number(c) as Int)
    .filter((v) => !isNaN(v));
}
</script>

<style scoped></style>
