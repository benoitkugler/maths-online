<template>
  <v-list-item
    v-if="props.asListItem"
    prepend-icon="mdi-share-variant"
    :title="title"
    @click="onClick"
  >
  </v-list-item>
  <v-btn v-else class="mx-1" size="small" variant="flat" @click="onClick">
    <template v-slot:prepend>
      <v-icon icon="mdi-share-variant" size="small" class="mr-4"></v-icon>
    </template>

    {{ title }}
  </v-btn>
</template>

<script setup lang="ts">
import { type Origin, PublicStatus } from "@/controller/api_gen";
import { computed } from "vue";
import { useRouter } from "vue-router";

interface Props {
  origin: Origin;
  asListItem?: boolean;
}
const props = defineProps<Props>();
const emit = defineEmits<{
  (e: "updatePublic", isPublic: boolean): void;
  (e: "createReview"): void;
}>();

const router = useRouter();

function goToReview() {
  const id = props.origin.IsInReview.Id;
  router.push({ name: "reviews", query: { id: id } });
}

const title = computed(() => {
  // do we have an active review
  if (props.origin.IsInReview.InReview) {
    return "Aller à la demande en cours";
  }
  switch (props.origin.PublicStatus) {
    case PublicStatus.NotAdmin:
      return "Publier...";
    case PublicStatus.AdminPublic:
      return "Masquer à la communauté";
    case PublicStatus.AdminNotPublic:
      return "Partager à la communauté";
  }
  return "";
});

function onClick() {
  if (props.origin.IsInReview.InReview) {
    goToReview();
  } else {
    switch (props.origin.PublicStatus) {
      case PublicStatus.NotAdmin:
        emit("createReview");
        return;
      case PublicStatus.AdminPublic:
        emit("updatePublic", false);
        return;
      case PublicStatus.AdminNotPublic:
        emit("updatePublic", true);
        return;
    }
  }
}
</script>
