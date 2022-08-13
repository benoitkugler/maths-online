<template>
  <v-card class="my-2">
    <v-row class="bg-secondary pa-2 rounded" no-gutters>
      <v-col md="9" align-self="center"> Réponse ordonnée attendue. </v-col>
      <v-spacer></v-spacer>
      <v-col align-self="center" style="text-align: right">
        <v-btn
          icon
          @click="addAnswer"
          title="Ajouter un élément"
          size="x-small"
          class="mr-2"
        >
          <v-icon icon="mdi-plus" color="green" small></v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-row no-gutters class="">
      <v-col>
        <v-list class="overflow-y-auto" style="max-height: 50vh">
          <v-list-item
            v-for="(param, index) in props.modelValue.Answer"
            :key="index"
            class="pr-0"
            :ref="(el:any) => (answerPropsRefs[index] = el as Element)"
          >
            <v-row no-gutters>
              <v-col>
                <interpolated-text v-model="props.modelValue.Answer![index]">
                </interpolated-text>
              </v-col>
              <v-col cols="auto">
                <v-btn
                  icon
                  size="small"
                  flat
                  @click="removeAnswer(index)"
                  title="Supprimer cet élément"
                >
                  <v-icon icon="mdi-delete" color="red"></v-icon>
                </v-btn>
              </v-col>
            </v-row>
          </v-list-item>
        </v-list>
      </v-col>
    </v-row>
  </v-card>

  <v-card>
    <v-row class="bg-secondary pa-2 rounded" no-gutters>
      <v-col md="9" align-self="center"> Champs additionnels. </v-col>
      <v-spacer></v-spacer>
      <v-col align-self="center" style="text-align: right">
        <v-btn
          icon
          @click="addAdditionalProposal"
          title="Ajouter un élément"
          size="x-small"
          class="mr-2"
        >
          <v-icon icon="mdi-plus" color="green" small></v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-row no-gutters class="mt-3">
      <v-col>
        <v-list>
          <v-list-item
            v-for="(param, index) in props.modelValue.AdditionalProposals"
            :key="index"
            class="pr-0"
          >
            <v-row no-gutters>
              <v-col>
                <interpolated-text
                  v-model="props.modelValue.AdditionalProposals![index]"
                >
                </interpolated-text>
              </v-col>
              <v-col cols="auto">
                <v-btn
                  icon
                  size="small"
                  flat
                  @click="removeAdditionalProposal(index)"
                  title="Supprimer cet élément"
                >
                  <v-icon icon="mdi-delete" color="red"></v-icon>
                </v-btn>
              </v-col>
            </v-row>
          </v-list-item>
        </v-list>
      </v-col>
    </v-row>
  </v-card>

  <v-row class="my-2">
    <v-col>
      <interpolated-text
        force-latex
        v-model="props.modelValue.Label"
        label="Préfixe"
        hint="Code LaTeX ajouté devant le champ de réponse. (Optionnel)."
      ></interpolated-text>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { OrderedListFieldBlock } from "@/controller/api_gen";
import { ref } from "@vue/reactivity";
import { nextTick } from "@vue/runtime-core";
import InterpolatedText from "../utils/InterpolatedText.vue";

interface Props {
  modelValue: OrderedListFieldBlock;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: OrderedListFieldBlock): void;
}>();

const answerPropsRefs = ref<any>([]);

function addAnswer() {
  props.modelValue.Answer?.push("$x$");
  emit("update:modelValue", props.modelValue);
  nextTick(() => {
    answerPropsRefs.value[
      answerPropsRefs.value.length - 1
    ]?.$el.scrollIntoView();
  });
}

function removeAnswer(index: number) {
  props.modelValue.Answer?.splice(index, 1);
  emit("update:modelValue", props.modelValue);
}

function addAdditionalProposal() {
  props.modelValue.AdditionalProposals?.push("$y$");
  emit("update:modelValue", props.modelValue);
}

function removeAdditionalProposal(index: number) {
  props.modelValue.AdditionalProposals?.splice(index, 1);
  emit("update:modelValue", props.modelValue);
}
</script>

<style></style>
