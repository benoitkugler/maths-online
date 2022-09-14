<template>
  <div
    @drop="onDropJSON"
    @dragover="onDragoverJSON"
    class="d-flex ma-2 px-2"
    style="
      border: 1px solid blue;
      border-radius: 10px;
      height: 96%;
      justify-content: center;
      align-items: center;
    "
    v-if="rows.length == 0"
  >
    Importer une question en faisant glisser un fichier (.isyro.json) ...
  </div>
  <div
    v-else
    style="height: 66vh; overflow-y: auto"
    @dragstart="onDragStart"
    @dragend="onDragEnd"
  >
    <drop-zone
      v-if="showDropZone"
      @drop="(origin) => swapBlocks(origin, 0)"
    ></drop-zone>
    <div
      v-for="(row, index) in rows"
      :key="index"
      :ref="el => (blockWidgets[index] = el as Element)"
    >
      <BlockContainer
        @delete="removeBlock(index)"
        :index="index"
        :kind="row.Props.Kind"
        :hide-content="showDropZone"
        :has-error="errorBlockIndex == index"
      >
        <component
          :model-value="row.Props.Data"
          @update:model-value="(v: any) => updateBlock(index, v)"
          :is="row.Component"
          :available-parameters="props.availableParameters"
        ></component>
      </BlockContainer>
      <drop-zone
        v-if="showDropZone"
        @drop="(origin) => swapBlocks(origin, index + 1)"
      ></drop-zone>
    </div>
  </div>
</template>

<script setup lang="ts">
import DropZone from "@/components/DropZone.vue";
import type { Block, Question, Variable } from "@/controller/api_gen";
import { BlockKind } from "@/controller/api_gen";
import { newBlock } from "@/controller/editor";
import { swapItems } from "@/controller/utils";
import { markRaw, ref } from "@vue/reactivity";
import { computed, nextTick, watch, type Component } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import BlockContainer from "./blocks/BlockContainer.vue";
import FigureAffineLineFieldVue from "./blocks/FigureAffineLineField.vue";
import FigureBlockVue from "./blocks/FigureBlock.vue";
import FigurePointFieldVue from "./blocks/FigurePointField.vue";
import FigureVectorFieldVue from "./blocks/FigureVectorField.vue";
import FigureVectorPairFieldVue from "./blocks/FigureVectorPairField.vue";
import FormulaBVue from "./blocks/FormulaB.vue";
import FormulaFieldVue from "./blocks/FormulaField.vue";
import FunctionPointsFieldVue from "./blocks/FunctionPointsField.vue";
import FunctionsGraphVue from "./blocks/FunctionsGraph.vue";
import NumberFieldVue from "./blocks/NumberField.vue";
import OrderedListFieldVue from "./blocks/OrderedListField.vue";
import ProofFieldVue from "./blocks/ProofField.vue";
import RadioFieldVue from "./blocks/RadioField.vue";
import SignTableVue from "./blocks/SignTable.vue";
import SignTableFieldVue from "./blocks/SignTableField.vue";
import TableVue from "./blocks/TableB.vue";
import TableFieldVue from "./blocks/TableField.vue";
import TextVue from "./blocks/TextB.vue";
import TreeFieldVue from "./blocks/TreeField.vue";
import VariationTableVue from "./blocks/VariationTable.vue";
import VariationTableFieldVue from "./blocks/VariationTableField.vue";
import VectorFieldVue from "./blocks/VectorField.vue";

interface Props {
  modelValue: Block[];
  errorBlockIndex?: number;
  availableParameters: Variable[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:modelValue", content: Block[]): void;
  (e: "importQuestion", content: Question): void;
}>();

const rows = computed(() => props.modelValue.map(dataToBlock));

interface block {
  Props: Block;
  Component: Component;
}

function dataToBlock(data: Block): block {
  switch (data.Kind) {
    case BlockKind.TextBlock:
      return { Props: data, Component: markRaw(TextVue) };
    case BlockKind.FormulaBlock:
      return { Props: data, Component: markRaw(FormulaBVue) };
    case BlockKind.FigureBlock:
      return { Props: data, Component: markRaw(FigureBlockVue) };
    case BlockKind.FunctionsGraphBlock:
      return { Props: data, Component: markRaw(FunctionsGraphVue) };
    case BlockKind.VariationTableBlock:
      return { Props: data, Component: markRaw(VariationTableVue) };
    case BlockKind.SignTableBlock:
      return { Props: data, Component: markRaw(SignTableVue) };
    case BlockKind.TableBlock:
      return { Props: data, Component: markRaw(TableVue) };
    case BlockKind.NumberFieldBlock:
      return { Props: data, Component: markRaw(NumberFieldVue) };
    case BlockKind.ExpressionFieldBlock:
      return { Props: data, Component: markRaw(FormulaFieldVue) };
    case BlockKind.RadioFieldBlock:
      return { Props: data, Component: markRaw(RadioFieldVue) };
    case BlockKind.OrderedListFieldBlock:
      return { Props: data, Component: markRaw(OrderedListFieldVue) };
    case BlockKind.FigurePointFieldBlock:
      return { Props: data, Component: markRaw(FigurePointFieldVue) };
    case BlockKind.FigureVectorFieldBlock:
      return { Props: data, Component: markRaw(FigureVectorFieldVue) };
    case BlockKind.VariationTableFieldBlock:
      return { Props: data, Component: markRaw(VariationTableFieldVue) };
    case BlockKind.SignTableFieldBlock:
      return { Props: data, Component: markRaw(SignTableFieldVue) };
    case BlockKind.FunctionPointsFieldBlock:
      return { Props: data, Component: markRaw(FunctionPointsFieldVue) };
    case BlockKind.FigureVectorPairFieldBlock:
      return { Props: data, Component: markRaw(FigureVectorPairFieldVue) };
    case BlockKind.FigureAffineLineFieldBlock:
      return { Props: data, Component: markRaw(FigureAffineLineFieldVue) };
    case BlockKind.TreeFieldBlock:
      return { Props: data, Component: markRaw(TreeFieldVue) };
    case BlockKind.TableFieldBlock:
      return { Props: data, Component: markRaw(TableFieldVue) };
    case BlockKind.VectorFieldBlock:
      return { Props: data, Component: markRaw(VectorFieldVue) };
    case BlockKind.ProofFieldBlock:
      return { Props: data, Component: markRaw(ProofFieldVue) };
    default:
      throw "dataToBlock: unexpected Kind";
  }
}

defineExpose({ addBlock });

const blockWidgets = ref<(Element | null)[]>([]);

watch(props, () => {
  if (props.errorBlockIndex != null) {
    blockWidgets.value[props.errorBlockIndex]?.scrollIntoView();
  }
});

function addBlock(kind: BlockKind) {
  props.modelValue.push(newBlock(kind));
  emit("update:modelValue", props.modelValue);

  nextTick(() => {
    const L = blockWidgets.value?.length;
    if (L) {
      blockWidgets.value[L - 1]?.scrollIntoView();
    }
  });
}

function updateBlock(index: number, data: Block["Data"]) {
  props.modelValue[index].Data = data;
  emit("update:modelValue", props.modelValue);
}

function removeBlock(index: number) {
  props.modelValue.splice(index, 1);
  emit("update:modelValue", props.modelValue);
}

/** take the block at the index `origin` and insert it right before
the block at index `target` (which is between 0 and nbBlocks)
 */
function swapBlocks(origin: number, target: number) {
  const out = swapItems(origin, target, props.modelValue);
  emit("update:modelValue", out);
}

let showDropZone = $ref(false);

function onDragStart() {
  setTimeout(() => (showDropZone = true), 100); // workaround bug
}

function onDragEnd(ev: DragEvent) {
  showDropZone = false;
}

async function onDropJSON(ev: DragEvent) {
  if (ev.dataTransfer?.files.length) {
    ev.preventDefault();
    const content = (await ev.dataTransfer?.files[0].text()) || "";
    const question = JSON.parse(content);
    emit("importQuestion", question);
  }
}

function onDragoverJSON(ev: DragEvent) {
  if (ev.dataTransfer?.files.length || ev.dataTransfer?.items.length) {
    ev.preventDefault();
  }
}
</script>

<style scoped></style>
