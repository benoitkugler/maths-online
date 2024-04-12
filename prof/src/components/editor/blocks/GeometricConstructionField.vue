<template>
  <v-select
    class="my-2"
    density="compact"
    variant="outlined"
    :items="[
      {
        title: 'Construire un point',
        value: GeoFieldKind.GFPoint,
      },
      {
        title: 'Construire un vecteur',
        value: GeoFieldKind.GFVector,
      },
      {
        title: 'Construire une droite',
        value: GeoFieldKind.GFAffineLine,
      },
      {
        title: 'Construire une paire de vecteurs',
        value: GeoFieldKind.GFVectorPair,
      },
    ]"
    label="Choisir le type du champ de réponse"
    hide-details
    :model-value="props.modelValue.Field.Kind"
    @update:model-value="onChangeFieldKind"
  >
  </v-select>

  <GFPointW
    v-if="props.modelValue.Field.Kind == GeoFieldKind.GFPoint"
    :model-value="(props.modelValue.Field.Data as GFPoint)"
    @update:model-value="
      (v) => {
        props.modelValue.Field.Data = v;
        emitUpdate();
      }
    "
  ></GFPointW>
  <GFVectorW
    v-else-if="props.modelValue.Field.Kind == GeoFieldKind.GFVector"
    :model-value="(props.modelValue.Field.Data as GFVector)"
    @update:model-value="
      (v) => {
        props.modelValue.Field.Data = v;
        emitUpdate();
      }
    "
  ></GFVectorW>
  <GFAffineLineW
    v-else-if="props.modelValue.Field.Kind == GeoFieldKind.GFAffineLine"
    :model-value="(props.modelValue.Field.Data as GFAffineLine)"
    @update:model-value="
      (v) => {
        props.modelValue.Field.Data = v;
        emitUpdate();
      }
    "
  ></GFAffineLineW>
  <GFVectorPairW
    v-else-if="props.modelValue.Field.Kind == GeoFieldKind.GFVectorPair"
    :model-value="(props.modelValue.Field.Data as GFVectorPair)"
    @update:model-value="
      (v) => {
        props.modelValue.Field.Data = v;
        emitUpdate();
      }
    "
  ></GFVectorPairW>

  <v-select
    class="mt-4 mb-2"
    density="compact"
    variant="outlined"
    :items="[
      {
        title: 'Figures (droites, cercles)',
        value: FiguresOrGraphsKind.FigureBlock,
      },
      {
        title: 'Graphes (fonctions, suites)',
        value: FiguresOrGraphsKind.FunctionsGraphBlock,
      },
    ]"
    label="Choisir le type de figure de fond"
    hide-details
    :model-value="props.modelValue.Background.Kind"
    @update:model-value="onChangeBackgroundKind"
  >
  </v-select>
  <figure-block-vue
    v-if="props.modelValue.Background.Kind == FiguresOrGraphsKind.FigureBlock"
    :model-value="(props.modelValue.Background.Data as FigureBlock)"
    @update:model-value="
      (v) => {
        props.modelValue.Background.Data = v;
        emitUpdate();
      }
    "
    :available-parameters="props.availableParameters"
  ></figure-block-vue>
  <functions-graph
    v-else-if="
      props.modelValue.Background.Kind ==
      FiguresOrGraphsKind.FunctionsGraphBlock
    "
    :model-value="(props.modelValue.Background.Data as FunctionsGraphBlock)"
    @update:model-value="
      (v) => {
        props.modelValue.Background.Data = v;
        emitUpdate();
      }
    "
    :available-parameters="props.availableParameters"
  ></functions-graph>
</template>

<script setup lang="ts">
import {
  FiguresOrGraphsKind,
  VectorPairCriterion,
  type GeometricConstructionFieldBlock,
  type Variable,
  type FigureBlock,
  type GFPoint,
  type GFVector,
  type GFAffineLine,
  type GFVectorPair,
} from "@/controller/api_gen";
import FigureBlockVue from "./FigureBlock.vue";
import { GeoFieldKind } from "@/controller/api_gen";
import FunctionsGraph from "./FunctionsGraph.vue";
import type { FunctionsGraphBlock, Int } from "@/controller/api_gen";
import { xRune } from "@/controller/editor";
import GFPointW from "./GFPointW.vue";
import GFVectorW from "./GFVectorW.vue";
import GFAffineLineW from "./GFAffineLineW.vue";
import GFVectorPairW from "./GFVectorPairW.vue";

interface Props {
  modelValue: GeometricConstructionFieldBlock;
  availableParameters: Variable[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: GeometricConstructionFieldBlock): void;
}>();

function emitUpdate() {
  emit("update:modelValue", props.modelValue);
}

function onChangeBackgroundKind(kind: FiguresOrGraphsKind) {
  props.modelValue.Background.Kind = kind;
  switch (kind) {
    case FiguresOrGraphsKind.FigureBlock:
      props.modelValue.Background.Data = {
        ShowGrid: true,
        ShowOrigin: true,
        Bounds: {
          Width: 10 as Int,
          Height: 10 as Int,
          Origin: { X: 3, Y: 3 },
        },
        Drawings: {
          Lines: [],
          Points: [],
          Segments: [],
          Circles: [],
          Areas: [],
        },
      };
      break;
    case FiguresOrGraphsKind.FunctionsGraphBlock:
      props.modelValue.Background.Data = {
        FunctionExprs: [
          {
            Function: "abs(x) + sin(x)",
            Decoration: {
              Label: "f",
              Color: "",
            },
            Variable: { Name: xRune, Indice: "" },
            From: "-5",
            To: "5",
          },
        ],
        SequenceExprs: [],
        FunctionVariations: [],
        Areas: [],
        Points: [],
      };
  }
}

function onChangeFieldKind(kind: GeoFieldKind) {
  props.modelValue.Field.Kind = kind;
  switch (kind) {
    case GeoFieldKind.GFPoint: {
      const data: GFPoint = {
        Answer: {
          X: "3",
          Y: "4",
        },
      };
      props.modelValue.Field.Data = data;
      break;
    }
    case GeoFieldKind.GFVector: {
      const data: GFVector = {
        Answer: {
          X: "3",
          Y: "4",
        },
        MustHaveOrigin: true,
        AnswerOrigin: {
          X: "0",
          Y: "1",
        },
      };
      props.modelValue.Field.Data = data;
      break;
    }
    case GeoFieldKind.GFAffineLine: {
      const data: GFAffineLine = {
        Label: "D",
        A: "1",
        B: "-2",
      };
      props.modelValue.Field.Data = data;
      break;
    }
    case GeoFieldKind.GFVectorPair: {
      const data: GFVectorPair = {
        Criterion: VectorPairCriterion.VectorColinear,
      };
      props.modelValue.Field.Data = data;
      break;
    }
  }
}
</script>

<style scoped></style>
