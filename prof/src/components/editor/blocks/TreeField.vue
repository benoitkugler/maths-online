<template>
  Arbre de probabilités régulier : tous les niveaux ont la même structure.
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
            <v-text-field
              variant="outlined"
              density="compact"
              v-model="props.modelValue.EventsProposals![index]"
              hide-details
            >
            </v-text-field>
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
      <v-col md="9" align-self="center">
        <v-card-subtitle class="py-2">
          Réponse : profondeur de l'arbre
        </v-card-subtitle>
      </v-col>
    </v-row>
    <v-card-text>
      <v-text-field
        variant="outlined"
        density="compact"
        label="Profondeur"
        hint="Nombre de niveaux de l'arbre"
        :model-value="depth"
        @update:model-value="onChooseDepth"
        persistent-hint
      ></v-text-field>
    </v-card-text>
  </v-card>

  <v-card color="secondary" class="my-2">
    <v-row no-gutters>
      <v-col md="9" align-self="center">
        <v-card-subtitle>
          Réponse : évènements et probabilités
        </v-card-subtitle>
      </v-col>
      <v-spacer></v-spacer>
      <v-col align-self="center" style="text-align: right">
        <v-btn
          icon
          @click="addRootChildren"
          title="Ajouter une issue"
          size="x-small"
          class="mr-2 my-2"
        >
          <v-icon icon="mdi-plus" color="green" small></v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-row no-gutters>
      <v-col>
        <v-list>
          <v-list-item
            v-for="(param, index) in rootChildren"
            class="pr-0"
            :key="index"
          >
            <v-row no-gutters class="mt-2">
              <v-col>
                <v-text-field
                  variant="outlined"
                  density="compact"
                  :model-value="param.proba"
                  @update:model-value="s => props.modelValue.AnswerRoot.Probabilities![index] = s"
                  hide-details
                  label="Probabilité"
                  :color="color"
                  class="mr-2"
                >
                </v-text-field>
              </v-col>
              <v-col>
                <v-select
                  variant="outlined"
                  density="compact"
                  label="Evènement"
                  :items="props.modelValue.EventsProposals || []"
                  :model-value="props.modelValue.EventsProposals![param.value]"
                  @update:model-value="s => props.modelValue.AnswerRoot.Children![index].Value = props.modelValue.EventsProposals?.indexOf(s)!"
                  hide-details
                ></v-select>
              </v-col>
              <v-col cols="auto">
                <v-btn
                  icon
                  size="small"
                  flat
                  @click="removeRootChildren(index)"
                  title="Supprimer cette issue"
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
</template>

<script setup lang="ts">
// TODO: for now, we only support regular trees
// custom trees require v-treeview component, not available yet

import type {
  TreeFieldBlock,
  TreeNodeAnswer,
  Variable,
} from "@/controller/api_gen";
import { ExpressionColor } from "@/controller/editor";
import { computed } from "@vue/runtime-core";

interface Props {
  modelValue: TreeFieldBlock;
  availableParameters: Variable[];
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: TreeFieldBlock): void;
}>();

const color = ExpressionColor;

const rootChildren = computed(() => {
  return (props.modelValue.AnswerRoot.Probabilities || []).map((v, i) => ({
    proba: v,
    value: props.modelValue.AnswerRoot.Children![i].Value,
  }));
});

function computeDepth(n: TreeNodeAnswer): number {
  if (!n.Children?.length) {
    return 0;
  }
  return 1 + computeDepth(n.Children[0]);
}

const depth = computed(() => computeDepth(props.modelValue.AnswerRoot));

function onChooseDepth(s: string) {
  const sInt = Number(s);
  if (isNaN(sInt)) {
    return;
  }

  props.modelValue.AnswerRoot = buildRegularTree(
    props.modelValue.AnswerRoot.Probabilities || [],
    props.modelValue.AnswerRoot.Children!.map((c) => c.Value),
    sInt
  );
  emit("update:modelValue", props.modelValue);
}

function buildRegularTree(
  probalities: string[],
  events: number[],
  maxDepth: number
): TreeNodeAnswer {
  const buildWithValue = (value: number, depth: number): TreeNodeAnswer => {
    if (depth == 0) {
      return { Probabilities: [], Children: [], Value: value };
    }
    return {
      Probabilities: probalities,
      Children: events.map((e) => buildWithValue(e, depth - 1)),
      Value: value,
    };
  };

  return buildWithValue(0, maxDepth);
}

function addEventsProposals() {
  props.modelValue.EventsProposals?.push("");
}

function removeEventsProposals(index: number) {
  props.modelValue.EventsProposals?.splice(index, 1);
}

function addRootChildren() {
  props.modelValue.AnswerRoot.Children?.push({
    Children: [],
    Probabilities: [],
    Value: 0,
  });
  props.modelValue.AnswerRoot.Probabilities?.push("0.5");
}

function removeRootChildren(index: number) {
  props.modelValue.AnswerRoot.Children?.splice(index, 1);
  props.modelValue.AnswerRoot.Probabilities?.splice(index, 1);
}
</script>

<style></style>
