<template>
  <v-row no-gutters class="px-2">
    <v-col>
      <v-btn icon flat title="Supprimer" small>
        <v-icon small color="red" @click="emit('delete')">mdi-close</v-icon>
      </v-btn>
    </v-col>
    <v-col md="auto" style="text-align: right">
      <slot name="toolbar"></slot>
    </v-col>
  </v-row>
  <v-row no-gutters class="px-4 my-0 py-0">
    <slot></slot>
  </v-row>
</template>

<script setup lang="ts">
const emit = defineEmits<{
  (e: "delete"): void;
}>();
</script>

<style></style>
