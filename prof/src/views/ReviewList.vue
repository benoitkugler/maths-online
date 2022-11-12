<template>
  <v-container class="pb-1 fill-height">
    <v-card v-if="currentReview == null" class="my-5 mx-auto" width="90%">
      <v-row class="mx-0">
        <v-col cols="9">
          <v-card-title>Publications</v-card-title>
          <v-card-subtitle>
            Discuter des demandes de publications dans la base officielle
          </v-card-subtitle>
        </v-col>
      </v-row>

      <v-alert color="info" class="my-2 mx-6" density="compact" closable>
        La base officielle contient des resources dont le contenu est vérifié,
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
      class="mx-auto"
      :review="currentReview"
      @back="currentReview = null"
      v-else
    ></review-pannel>
  </v-container>
</template>

<script setup lang="ts">
import type { ReviewHeader } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { onActivated, onMounted } from "vue";
import { $ref } from "vue/macros";
import ReviewRow from "../components/review/ReviewRow.vue";
import ReviewPannel from "../components/review/ReviewPannel.vue";
import { useRoute, useRouter } from "vue-router";

let reviews = $ref<ReviewHeader[]>([]);
let currentReview = $ref<ReviewHeader | null>(null);

onMounted(init);
onActivated(init);

const route = useRoute();

async function init() {
  await fetchReviews();
  parseQuery(); // require the updated list of reviews
}

async function fetchReviews() {
  const res = await controller.ReviewsList();
  if (res == undefined) return;
  reviews = res || [];
}

function goToReview(review: ReviewHeader) {
  currentReview = review;
}

function parseQuery() {
  const id = Number(route.query["id"]);
  if (isNaN(id)) return;
  const review = reviews.find((re) => re.Id == id);
  if (review == undefined) return;
  goToReview(review);
}
</script>
