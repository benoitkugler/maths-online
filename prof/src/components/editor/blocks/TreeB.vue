<template>
  Arbre de probabilités régulier : le nombre de branches d'un niveau est
  constant.
  <v-card color="secondary" class="mt-2">
    <v-row no-gutters>
      <v-col md="9" align-self="center">
        <v-card-subtitle> Evènements possibles </v-card-subtitle>
      </v-col>
      <v-spacer></v-spacer>
      <v-col align-self="center" style="text-align: right">
        <v-btn
          icon
          @click="addEventsProposals"
          title="Ajouter un évènement"
          size="x-small"
          class="mr-2 my-2"
        >
          <v-icon icon="mdi-plus" color="green" small></v-icon>
        </v-btn>
      </v-col>
    </v-row>

    <v-list>
      <v-list-item
        v-for="(param, index) in props.modelValue.EventsProposals"
        :key="index"
        class="pr-0"
      >
        <v-row no-gutters>
          <v-col>
            <interpolated-text
              v-model="props.modelValue.EventsProposals![index]"
            ></interpolated-text>
          </v-col>
          <v-col cols="auto">
            <v-btn
              icon
              size="small"
              flat
              @click="removeEventsProposals(index)"
              title="Supprimer cet évènement"
            >
              <v-icon icon="mdi-delete" color="red"></v-icon>
            </v-btn>
          </v-col>
        </v-row>
      </v-list-item>
    </v-list>
  </v-card>

  <v-card class="my-2">
    <v-row no-gutters class="bg-secondary">
      <v-col align-self="center">
        <v-card-subtitle class="py-2"> Réponse </v-card-subtitle>
      </v-col>
      <v-col cols="auto">
        <v-menu>
          <template v-slot:activator="{ isActive, props: innerProps }">
            <v-btn
              v-on="{ isActive }"
              v-bind="innerProps"
              class="mx-1 my-2"
              size="small"
              title="Ajouter un niveau..."
              @click.stop
            >
              <template v-slot:prepend>
                <v-icon color="green">mdi-plus</v-icon>
              </template>
              Niveau
            </v-btn>
          </template>
          <v-list>
            <v-list-item density="compact" @click="addLevel(2)">
              Niveau à <v-chip>2</v-chip> branches
            </v-list-item>
            <v-list-item density="compact" @click="addLevel(3)">
              Niveau à <v-chip>3</v-chip> branches
            </v-list-item>
            <v-list-item density="compact" @click="addLevel(4)">
              Niveau à <v-chip>4</v-chip> branches
            </v-list-item>
          </v-list>
        </v-menu>

        <v-btn
          class="mx-2"
          size="small"
          title="Supprimer le dernier niveau"
          @click="deleteLastLevel"
          :disabled="!props.modelValue.AnswerRoot.Children?.length"
        >
          <template v-slot:prepend>
            <v-icon color="red">mdi-delete</v-icon>
          </template>
          Dernier niveau
        </v-btn>
      </v-col>
    </v-row>
    <v-card-text>
      <tree-node
        style="overflow-x: scroll; max-width: 80vh"
        :is-root="true"
        v-model="props.modelValue.AnswerRoot"
        :event-proposals="props.modelValue.EventsProposals || []"
      ></tree-node>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { TreeBlock, TreeNodeAnswer, Variable } from "@/controller/api_gen";
import TreeNode from "./TreeNode.vue";
import { copy } from "@/controller/utils";
import InterpolatedText from "../utils/InterpolatedText.vue";

interface Props {
  modelValue: TreeBlock;
  availableParameters: Variable[];
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: TreeBlock): void;
}>();

function addEventsProposals() {
  props.modelValue.EventsProposals = (
    props.modelValue.EventsProposals || []
  ).concat("");
}

function removeEventsProposals(index: number) {
  props.modelValue.EventsProposals?.splice(index, 1);
  const L = props.modelValue.EventsProposals?.length || 0;
  // make sure indexes in tree are still valid
  const clamp = (node: TreeNodeAnswer) => {
    if (node.Value >= L && L > 0) {
      node.Value = L - 1;
    }
    // recurse
    node.Children?.forEach(clamp);
  };
  clamp(props.modelValue.AnswerRoot);
}

function deleteLastLevel() {
  const del = (node: TreeNodeAnswer) => {
    if (!node.Children?.length) return;
    if (!node?.Children[0]?.Children?.length) {
      // we are at the correct level
      node.Children = [];
      node.Probabilities = [];
    } else {
      // recurse
      node.Children?.forEach(del);
    }
  };
  del(props.modelValue.AnswerRoot);
  emit("update:modelValue", props.modelValue);
}

const defaultProbas = [
  [],
  ["1"],
  ["0.5", "0.5"],
  ["0.3", "0.3", "0.4"],
  ["0.25", "0.25", "0.25", "0.25"]
];

function addLevel(branchesCount: number) {
  const add = (node: TreeNodeAnswer) => {
    if (node.Children?.length) {
      // recurse
      node.Children.forEach(add);
    } else {
      // we are at the last level
      node.Probabilities = copy(defaultProbas[branchesCount]);
      node.Children = node.Probabilities.map(() => ({
        Value: 0,
        Probabilities: [],
        Children: []
      }));
    }
  };
  add(props.modelValue.AnswerRoot);
  emit("update:modelValue", props.modelValue);
}
</script>

<style></style>
