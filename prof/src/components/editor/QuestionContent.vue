<template>
  <div
    v-if="rows.length == 0"
    @drop="onDropJSON"
    @dragover="onDragoverJSON"
    @dragleave="hasDragOverJSON = false"
    class="d-flex ma-2"
    :style="{
      border: hasDragOverJSON ? '2px dashed lightblue' : '2px solid lightblue',
      'border-radius': '10px',
      height: '95%',
      'justify-content': 'center',
      'align-items': 'center',
    }"
  >
    Importer une question en faisant glisser un fichier (.isyro.json) ...
  </div>

  <div
    v-else
    style="height: 72vh; overflow-y: auto"
    @dragstart="onDragStart"
    @dragend="onDragEnd"
  >
    <drop-zone
      v-if="showDropZone"
      @drop="(origin: number) => swapBlocks(origin, 0)"
    ></drop-zone>
    <div
      v-for="(row, index) in rows"
      :key="index"
      :ref="el => (blockWidgets[index] = el as Element)"
    >
      <BlockContainer
        @delete="removeBlock(index)"
        @add-syntax-hint="addSyntaxHint(index)"
        :index="index"
        :kind="row.Props.Kind"
        :hide-content="showDropZone || !areExpanded[index]"
        :has-error="errorBlockIndex == index"
        @toggle-content="areExpanded[index] = !areExpanded[index]"
        @copy="copyBlock(index)"
      >
        <component
          :model-value="row.Props.Data"
          :available-parameters="props.availableParameters"
          @update:model-value="(v: any) => updateBlock(index, v)"
          :is="row.Component"
        ></component>
      </BlockContainer>
      <drop-zone
        v-if="showDropZone"
        @drop="(origin: number) => swapBlocks(origin, index + 1)"
      ></drop-zone>
    </div>
  </div>
</template>

<script setup lang="ts">
import DropZone from "@/components/DropZone.vue";
import type {
  Block,
  ExpressionFieldBlock,
  Variable,
} from "@/controller/api_gen";
import { BlockKind } from "@/controller/api_gen";
import { newBlock } from "@/controller/editor";
import { swapItems } from "@/controller/utils";
import BlockContainer from "./blocks/BlockContainer.vue";
import ExpressionFieldVue from "./blocks/ExpressionField.vue";
import FigureBlockVue from "./blocks/FigureBlock.vue";
import FormulaBVue from "./blocks/FormulaB.vue";
import FunctionPointsFieldVue from "./blocks/FunctionPointsField.vue";
import FunctionsGraphVue from "./blocks/FunctionsGraph.vue";
import GeometricConstructionFieldVue from "./blocks/GeometricConstructionField.vue";
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
import TreeB from "./blocks/TreeB.vue";
import SetFieldVue from "./blocks/SetField.vue";
import ImageVue from "./blocks/ImageB.vue";
import { computed } from "vue";
import { ref } from "vue";
import { type Component } from "vue";
import { markRaw } from "vue";
import { watch } from "vue";
import { nextTick } from "vue";
import { controller } from "@/controller/controller";

interface Props {
  modelValue: Block[];
  errorBlockIndex?: number;
  availableParameters: Variable[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:modelValue", content: Block[]): void;
  (e: "importQuestion", json: string): void;
  (e: "addSyntaxHint", expr: ExpressionFieldBlock): void;
}>();

const rows = computed(() => props.modelValue.map(dataToBlock));
const areExpanded = ref(props.modelValue.map(() => true));

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
      return { Props: data, Component: markRaw(ExpressionFieldVue) };
    case BlockKind.RadioFieldBlock:
      return { Props: data, Component: markRaw(RadioFieldVue) };
    case BlockKind.OrderedListFieldBlock:
      return { Props: data, Component: markRaw(OrderedListFieldVue) };
    case BlockKind.GeometricConstructionFieldBlock:
      return { Props: data, Component: markRaw(GeometricConstructionFieldVue) };
    case BlockKind.VariationTableFieldBlock:
      return { Props: data, Component: markRaw(VariationTableFieldVue) };
    case BlockKind.SignTableFieldBlock:
      return { Props: data, Component: markRaw(SignTableFieldVue) };
    case BlockKind.FunctionPointsFieldBlock:
      return { Props: data, Component: markRaw(FunctionPointsFieldVue) };
    case BlockKind.TreeBlock:
      return { Props: data, Component: markRaw(TreeB) };
    case BlockKind.TreeFieldBlock:
      return { Props: data, Component: markRaw(TreeFieldVue) };
    case BlockKind.TableFieldBlock:
      return { Props: data, Component: markRaw(TableFieldVue) };
    case BlockKind.VectorFieldBlock:
      return { Props: data, Component: markRaw(VectorFieldVue) };
    case BlockKind.ProofFieldBlock:
      return { Props: data, Component: markRaw(ProofFieldVue) };
    case BlockKind.SetFieldBlock:
      return { Props: data, Component: markRaw(SetFieldVue) };
    case BlockKind.ImageBlock:
      return { Props: data, Component: markRaw(ImageVue) };
  }
}

defineExpose({ addBlock, addExistingBlock });

const blockWidgets = ref<(Element | null)[]>([]);

watch(props, () => {
  if (props.errorBlockIndex != null) {
    blockWidgets.value[props.errorBlockIndex]?.scrollIntoView();
  }
});

function addBlock(kind: BlockKind) {
  const l = props.modelValue;
  l.push(newBlock(kind));
  areExpanded.value.push(true);
  emit("update:modelValue", l);

  nextTick(() => {
    const L = blockWidgets.value?.length;
    if (L) {
      blockWidgets.value[L - 1]?.scrollIntoView();
    }
  });
}

function addExistingBlock(block: Block) {
  const l = props.modelValue;
  l.push(block);
  areExpanded.value.push(true);
  emit("update:modelValue", l);

  nextTick(() => {
    const L = blockWidgets.value?.length;
    if (L) {
      blockWidgets.value[L - 1]?.scrollIntoView();
    }
  });
}

function updateBlock(index: number, data: Block["Data"]) {
  const l = props.modelValue;
  l[index].Data = data;
  emit("update:modelValue", l);
}

function removeBlock(index: number) {
  const l = props.modelValue;
  l.splice(index, 1);
  areExpanded.value.splice(index, 1);
  emit("update:modelValue", l);
}

/** take the block at the index `origin` and insert it right before
the block at index `target` (which is between 0 and nbBlocks)
 */
function swapBlocks(origin: number, target: number) {
  const out = swapItems(origin, target, props.modelValue);
  areExpanded.value = swapItems(origin, target, areExpanded.value);
  emit("update:modelValue", out);
}

const showDropZone = ref(false);

function onDragStart() {
  setTimeout(() => (showDropZone.value = true), 100); // workaround bug
}

function onDragEnd(_: DragEvent) {
  showDropZone.value = false;
}

async function onDropJSON(ev: DragEvent) {
  if (ev.dataTransfer?.files.length) {
    ev.preventDefault();
    const content = (await ev.dataTransfer?.files[0].text()) || "";
    emit("importQuestion", content);
    hasDragOverJSON.value = false;
  }
}

const hasDragOverJSON = ref(false);

function onDragoverJSON(ev: DragEvent) {
  if (ev.dataTransfer?.files.length || ev.dataTransfer?.items.length) {
    hasDragOverJSON.value = true;
    ev.preventDefault();
  }
}

async function addSyntaxHint(index: number) {
  const block = props.modelValue[index].Data as ExpressionFieldBlock;
  emit("addSyntaxHint", block);
}

async function copyBlock(index: number) {
  const json = JSON.stringify(props.modelValue[index]);
  await navigator.clipboard.writeText(json);
  controller.showMessage("Bloc copié dans le presse-papier.");
}
</script>

<style scoped></style>
