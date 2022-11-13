<template>
  <v-dialog v-model="showConfirmAccept" max-width="600px">
    <v-card title="Confirmer la validation">
      <v-card-text>
        Confirmez-vous la validation de la demande de publication ? <br /><br />

        La ressource sera ajoutée à la base officielle, et les commentaires
        associés à la demande seront effacés. <br />
        Le propriétaire de la ressource sera notifié par email.
      </v-card-text>
      <v-card-actions>
        <v-btn @click="showConfirmAccept = false" color="warning">Retour</v-btn>
        <v-spacer></v-spacer>
        <v-btn color="green" @click="acceptReview"> Accepter</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog v-model="showConfirmDelete" max-width="600px">
    <v-card title="Confirmer la suppression">
      <v-card-text>
        Confirmez-vous la suppression de la demande de publication ? <br />

        Les commentaires associés seront supprimés, et cette opération est
        irréversible.
      </v-card-text>
      <v-card-actions>
        <v-btn @click="showConfirmDelete = false" color="warning">Retour</v-btn>
        <v-spacer></v-spacer>
        <v-btn color="red" @click="deleteReview"> Supprimer</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card class="mx-auto" v-if="reviewExt != null" max-width="800px">
    <v-row>
      <v-col cols="auto" align-self="center" class="pr-2">
        <v-btn
          class="ma-2"
          size="small"
          icon
          title="Retour à la liste des publications"
          @click="emit('back')"
        >
          <v-icon icon="mdi-arrow-left"></v-icon>
        </v-btn>
      </v-col>
      <v-col>
        <v-card-title>{{ props.review.Title }}</v-card-title>
        <v-card-subtitle>{{ labels[props.review.Kind] }}</v-card-subtitle>
      </v-col>
      <v-col cols="auto" align-self="center" class="pr-6">
        <v-btn
          v-if="reviewExt.IsAcceptable"
          variant="outlined"
          color="success"
          class="mr-1"
          @click="showConfirmAccept = true"
        >
          Accepter la publication
        </v-btn>
        <v-btn
          v-if="reviewExt.IsDeletable"
          variant="outlined"
          color="red"
          @click="showConfirmDelete = true"
        >
          Retirer la publication
        </v-btn>
      </v-col>
    </v-row>
    <v-card-text class="py-0">
      <ApprovalArea
        class="mt-4 mb-1"
        :review="reviewExt"
        @update="updateApproval"
      ></ApprovalArea>
      <v-list>
        <v-list-item
          v-if="!reviewExt.Comments?.length"
          style="text-align: center"
        >
          <i>Soyez le premier à commentez...</i>
        </v-list-item>
        <CommentRow
          v-for="(comment, index) in reviewExt.Comments"
          :key="index"
          :comment="comment"
          @update="(m) => updateComment(m, index)"
          @delete="deleteComment(index)"
        ></CommentRow>

        <v-divider class="mt-3"></v-divider>
        <NewComment :disabled="isSending" @send="sendComment"></NewComment>
      </v-list>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import {
  Approval,
  ReviewKindLabels,
  type Comments,
  type ReviewExt,
  type ReviewHeader,
  type Time,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { onActivated, onMounted } from "vue";
import { $ref } from "vue/macros";
import CommentRow from "./CommentRow.vue";
import NewComment from "./NewComment.vue";
import ApprovalArea from "./ApprovalArea.vue";

interface Props {
  review: ReviewHeader;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "back"): void;
}>();

const labels = ReviewKindLabels;

let reviewExt = $ref<ReviewExt | null>(null);

onMounted(fetchData);
onActivated(fetchData);

async function fetchData() {
  const res = await controller.ReviewLoad({ id: props.review.Id });
  if (res == undefined) return;
  reviewExt = res;
}

function ownComments() {
  if (reviewExt == null) return [];
  return (
    reviewExt.Comments?.filter((cm) => cm.IsOwned).map((cm) => cm.Comment) || []
  );
}

let isSending = $ref(false);

async function _updateComments(comments: Comments) {
  isSending = true;
  const res = await controller.ReviewUpdateCommnents({
    IdReview: props.review.Id,
    Comments: comments,
  });
  isSending = false;
  if (res == undefined) return;

  fetchData();
}

function sendComment(comment: string) {
  const cms = ownComments();
  cms.push({
    Time: new Date(Date.now()).toISOString() as Time,
    Message: comment,
  });
  _updateComments(cms);
}

async function updateComment(message: string, index: number) {
  if (reviewExt == null) return;
  const comment = reviewExt.Comments![index];
  comment.Comment.Message = message;
  _updateComments(ownComments());
}

async function deleteComment(index: number) {
  if (reviewExt == null) return;
  reviewExt.Comments?.splice(index, 1);
  _updateComments(ownComments());
}

async function updateApproval(appro: Approval) {
  const res = await controller.ReviewUpdateApproval({
    IdReview: props.review.Id,
    Approval: appro,
  });
  if (res == undefined) return;
  fetchData();
}

let showConfirmDelete = $ref(false);

async function deleteReview() {
  showConfirmDelete = false;
  const res = await controller.ReviewDelete({ id: props.review.Id });
  if (res == undefined) return;
  emit("back");
}

let showConfirmAccept = $ref(false);
async function acceptReview() {
  showConfirmAccept = false;
  const res = await controller.ReviewAccept({ id: props.review.Id });
  if (res == undefined) return;
  emit("back");
}
</script>
