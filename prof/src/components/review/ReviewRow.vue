<template>
  <v-list-item link class="my-2">
    <v-row>
      <v-col cols="2" align-self="center">
        <v-chip color="primary-darken-1">
          {{ labels[props.review.Kind] }}
        </v-chip>
      </v-col>
      <v-col cols="8" align-self="center">
        ({{ props.review.Id }}) {{ props.review.Title }}
      </v-col>

      <v-col align-self="center">
        <v-chip> {{ props.review.NbComments }} commentaire(s) </v-chip>
      </v-col>
    </v-row>
  </v-list-item>
</template>

<script setup lang="ts">
import { ReviewKindLabels, type ReviewHeader } from "@/controller/api_gen";

interface Props {
  review: ReviewHeader;
}

const props = defineProps<Props>();

// const emit = defineEmits<{}>();

const labels = ReviewKindLabels;
</script>
