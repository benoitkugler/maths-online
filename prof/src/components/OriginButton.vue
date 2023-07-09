<template>
  <v-dialog v-model="confirmeCreate" max-width="800px">
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

  <v-menu offset-y close-on-content-click>
    <template v-slot:activator="{ isActive, props: innerProps }">
      <v-btn
        v-on="{ isActive }"
        v-bind="innerProps"
        class="mx-1"
        :size="props.variant == 'icon' ? 'x-small' : 'small'"
        :icon="props.variant == 'icon'"
        :variant="props.variant == 'text' ? 'flat' : undefined"
        title="Options de partage"
        @click.stop
        :color="isPersonnalAndShared ? 'blue' : undefined"
      >
        <v-icon
          icon="mdi-share-variant"
          size="small"
          :class="props.variant == 'text' ? 'mr-2' : ''"
        ></v-icon>
        <template v-if="props.variant == 'text'">Publier</template>
      </v-btn>
    </template>
    <OriginCard
      :origin="props.origin"
      @update="(b) => emit('updatePublic', b)"
      @create-review="confirmeCreate = true"
      @go-to-review="goToReview"
    ></OriginCard>
  </v-menu>
</template>

<script setup lang="ts">
import { Visibility, type Origin } from "@/controller/api_gen";
import { computed } from "vue";
import { useRouter } from "vue-router";
import { $ref } from "vue/macros";
import OriginCard from "./OriginCard.vue";

interface Props {
  origin: Origin;
  variant: "icon" | "text";
}
const props = defineProps<Props>();
const emit = defineEmits<{
  (e: "updatePublic", isPublic: boolean): void;
  (e: "createReview"): void;
}>();

const router = useRouter();

let confirmeCreate = $ref(false);

const isPersonnalAndShared = computed(
  () => props.origin.Visibility == Visibility.Personnal && props.origin.IsPublic
);

function goToReview(id: number) {
  router.push({ name: "reviews", query: { id: id } });
}
</script>
