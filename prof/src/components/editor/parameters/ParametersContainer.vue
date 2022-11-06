<template>
  <v-card>
    <v-row
      :style="{
        'background-color': props.isValidated ? 'lightgreen' : 'lightgray',
      }"
      class="rounded"
      no-gutters
    >
      <v-col md="8" align-self="center">
        <v-card-subtitle class="py-2">{{ props.title }}</v-card-subtitle>
      </v-col>
      <v-col align-self="center" style="text-align: right">
        <v-progress-circular
          v-show="props.isLoading"
          indeterminate
          class="mx-1"
        ></v-progress-circular>
        <v-btn
          icon
          @click="emit('add')"
          :title="addTitle"
          size="x-small"
          class="mr-2 my-2"
        >
          <v-icon icon="mdi-plus" color="green"></v-icon>
        </v-btn>
        <v-btn
          icon
          @click="emit('showHelp')"
          title="Aide"
          size="x-small"
          class="mr-2 my-2"
        >
          <v-icon icon="mdi-help" color="info"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-row no-gutters>
      <v-col> <slot></slot> </v-col>
    </v-row>
  </v-card>
</template>

<script setup lang="ts">
interface Props {
  title: string;
  addTitle: string;
  isLoading: boolean;
  isValidated: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "add"): void;
  (e: "showHelp"): void;
}>();
</script>

<style scoped></style>
