<template>
  <v-dialog v-model="confirmeCreate" max-width="800px" :eager="false">
    <v-card title="Confirmer la demande de publication">
      <v-card-text>
        Merci pour votre participation ! <br /><br />
        Le contenu officiel est vérifié par l'équipe Isyro, en tenant compte de
        l'avis de la communauté. <br />
        En continuant, vous ajouterez votre ressource à la liste des demandes de
        publications, et nous en prendrons connaissance au plus vite.
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          color="green"
          @click="
            emit('createReview');
            confirmeCreate = false;
          "
        >
          Créer une demande de publication
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-btn class="mx-1" size="small" variant="flat" @click.stop="onClick">
    <template v-slot:prepend>
      <v-icon icon="mdi-share-variant" size="small" class="mr-4"></v-icon>
    </template>

    {{ title }}
  </v-btn>
</template>

<script setup lang="ts">
import { type Origin, PublicStatus } from "@/controller/api_gen";
import { computed, ref } from "vue";
import { useRouter } from "vue-router";

interface Props {
  origin: Origin;
}
const props = defineProps<Props>();
const emit = defineEmits<{
  (e: "updatePublic", isPublic: boolean): void;
  (e: "createReview"): void;
}>();

const router = useRouter();

const confirmeCreate = ref(false);

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
        confirmeCreate.value = true;
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
