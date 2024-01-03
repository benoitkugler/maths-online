<template>
  <v-container class="pb-1 fill-height">
    <v-card v-if="currentReview == null" class="my-5 mx-auto" width="90%">
      <v-row class="mx-0">
        <v-col cols="9">
          <v-card-title>Publications</v-card-title>
        </v-col>
      </v-row>

      <v-alert
        variant="tonal"
        color="info"
        class="my-2 mx-6"
        density="compact"
        closable
      >
        La base officielle contient des ressources dont le contenu est vérifié,
        afin de garantir une qualité optimale. Vous pouvez participer à
        l'évalution des contenus suivants.
      </v-alert>

      <v-list class="mx-4">
        <review-row
          v-for="(review, index) in reviews"
          :key="index"
          :review="review"
          @click="goToReview(review)"
        ></review-row>
        <v-list-item v-if="!reviews.length" style="text-align: center">
          <i>Aucune demande de publication en cours.</i>
        </v-list-item>
      </v-list>
    </v-card>
    <review-pannel
      :review="currentReview"
      @back="backToList"
      v-else
    ></review-pannel>
  </v-container>
</template>

<script setup lang="ts">
import type { ReviewHeader } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { onActivated, onMounted, ref } from "vue";
import ReviewRow from "../components/review/ReviewRow.vue";
import ReviewPannel from "../components/review/ReviewPannel.vue";
import { useRoute } from "vue-router";

const reviews = ref<ReviewHeader[]>([]);
const currentReview = ref<ReviewHeader | null>(null);

onMounted(init);
onActivated(init);

const route = useRoute();

async function init() {
  await fetchReviews();
  parseQuery(); // require the updated list of reviews
}

async function fetchReviews() {
  const res = await controller.ReviewsList();
  reviews.value = res || [];
}

function goToReview(review: ReviewHeader) {
  currentReview.value = review;
}

function parseQuery() {
  const id = Number(route.query["id"]);
  if (isNaN(id)) return;
  const review = reviews.value.find((re) => re.Id == id);
  if (review == undefined) return;
  goToReview(review);
}

// update the list which may have changed due to accept or delete actions
async function backToList() {
  currentReview.value = null;
  fetchReviews();
}
</script>
