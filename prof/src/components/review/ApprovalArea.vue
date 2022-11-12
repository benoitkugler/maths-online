<template>
  <v-row no-gutters class="bg-grey-lighten-4 rounded pa-1">
    <v-col align-self="center" cols="auto" class="text-grey">
      Approbation :
    </v-col>
    <v-col align-self="center" class="mx-1">
      <v-chip
        :variant="isSelected(Approval.Opposed) ? 'outlined' : undefined"
        :size="isSelected(Approval.Opposed) ? 'large' : undefined"
        @click="emit('update', Approval.Opposed)"
        label
        prepend-icon="mdi-minus"
        color="red"
      >
        {{ props.review.Approvals[Approval.Opposed] }}
      </v-chip>
    </v-col>
    <v-col align-self="center" class="mx-1">
      <v-chip
        :variant="isSelected(Approval.Neutral) ? 'outlined' : undefined"
        :size="isSelected(Approval.Neutral) ? 'large' : undefined"
        @click="emit('update', Approval.Neutral)"
        label
      >
        {{ props.review.Approvals[Approval.Neutral] }}
      </v-chip>
    </v-col>
    <v-col align-self="center" class="mx-1">
      <v-chip
        :variant="isSelected(Approval.InFavor) ? 'outlined' : undefined"
        :size="isSelected(Approval.InFavor) ? 'large' : undefined"
        @click="emit('update', Approval.InFavor)"
        label
        append-icon="mdi-plus"
        color="green"
      >
        {{ props.review.Approvals[Approval.InFavor] }}
      </v-chip>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { Approval, type ReviewExt } from "@/controller/api_gen";
import { computed } from "@vue/reactivity";

interface Props {
  review: ReviewExt;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update", appro: Approval): void;
}>();

function isSelected(appro: Approval) {
  return props.review.UserApproval == appro;
}
</script>
