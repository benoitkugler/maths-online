<template>
  <v-dialog v-model="showConfirmDelete" max-width="600px">
    <v-card title="Confirmer la suppression">
      <v-card-text>
        Confirmez-vous la suppression de votre commentaire ? <br />

        Cette opération est irréversible.
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          variant="outlined"
          color="red"
          @click="
            emit('delete');
            showConfirmDelete = false;
          "
          >Supprimer</v-btn
        >
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card class="my-2">
    <v-row no-gutters :class="'rounded ma-0 py-2 ' + colorClass">
      <v-col align-self="center">
        <div class="mx-2">
          {{ props.comment.AuthorMail }}
        </div>
      </v-col>
      <v-spacer></v-spacer>
      <v-col align-self="center" style="text-align: right" class="mx-2">
        <i>
          {{ formatTime(props.comment.Comment.Time, true) }}
        </i>
      </v-col>
    </v-row>
    <v-card-text>
      <EditComment
        v-if="isEditing"
        :initial-message="props.comment.Comment.Message"
        @back="isEditing = false"
        @send="
          (m) => {
            isEditing = false;
            emit('update', m);
          }
        "
      ></EditComment>
      <v-row no-gutters v-else>
        <v-col>
          <div v-html="messageHTML"></div>
        </v-col>
        <v-col cols="auto" v-if="props.comment.IsOwned" align-self="center">
          <v-btn
            icon
            title="Editer"
            size="small"
            class="mx-1"
            @click="isEditing = true"
          >
            <v-icon color="primary">mdi-pencil</v-icon>
          </v-btn>
          <v-btn
            icon
            title="Supprimer"
            size="small"
            @click="showConfirmDelete = true"
          >
            <v-icon color="red">mdi-delete</v-icon>
          </v-btn>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { ReviewComment } from "@/controller/api_gen";
import { formatTime } from "@/controller/utils";
import { ref, computed } from "vue";
import EditComment from "./EditComment.vue";

interface Props {
  comment: ReviewComment;
}

const props = defineProps<Props>();

const messageHTML = computed(() =>
  props.comment.Comment.Message.replace(new RegExp(`\n`, "g"), "<br/>")
);

const emit = defineEmits<{
  (e: "update", message: string): void;
  (ed: "delete"): void;
}>();

const colorClass = computed(() =>
  props.comment.IsOwned ? "bg-secondary-lighten-1" : "bg-grey-lighten-3"
);

const isEditing = ref(false);

const showConfirmDelete = ref(false);
</script>
