<template>
  <v-card class="my-1">
    <v-card-subtitle class="bg-secondary py-3"
      >Options de la figure</v-card-subtitle
    >
    <v-card-text>
      <v-row>
        <v-col md="6">
          <v-text-field
            density="compact"
            variant="outlined"
            v-model.number="props.modelValue.Bounds.Width"
            label="Largeur"
            hint="Largeur de la figure, en nombre de carreaux"
            hide-details
          ></v-text-field>
        </v-col>
        <v-col md="6">
          <v-text-field
            hide-details
            density="compact"
            variant="outlined"
            v-model.number="props.modelValue.Bounds.Height"
            label="Hauteur"
            hint="Hauteur de la figure, en nombre de carreaux"
          ></v-text-field>
        </v-col>
      </v-row>
      <v-row>
        <v-col md="6">
          <v-text-field
            density="compact"
            variant="outlined"
            v-model.number="props.modelValue.Bounds.Origin.X"
            label="Origine : abscisse"
            hint="Abscisse de l'origine par rapport au coin inférieur gauche."
          ></v-text-field>
        </v-col>
        <v-col md="6">
          <v-text-field
            density="compact"
            variant="outlined"
            v-model.number="props.modelValue.Bounds.Origin.Y"
            label="Origine : ordonnée"
            hint="Ordonnée de l'origine par rapport au coin inférieur gauche."
          ></v-text-field>
        </v-col>
      </v-row>
      <v-row no-gutters>
        <v-col md="6">
          <v-checkbox
            density="compact"
            label="Afficher la grille"
            v-model="props.modelValue.ShowGrid"
            hide-details
          ></v-checkbox>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>

  <v-card color="secondary" class="my-1">
    <v-row no-gutters>
      <v-col align-self="center" md="9">
        <v-card-subtitle> Points </v-card-subtitle>
      </v-col>
      <v-col md="3" style="text-align: right">
        <v-btn
          icon
          @click="addPoint"
          title="Ajouter un point"
          size="x-small"
          class="mr-2 my-2"
        >
          <v-icon icon="mdi-plus" color="green" size="small"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-list>
      <div v-for="(point, index) in props.modelValue.Drawings.Points">
        <v-list-item>
          <v-row class="fix-input-width">
            <v-col md="3" align-self="center">
              <v-text-field
                density="compact"
                variant="outlined"
                label="Nom"
                v-model="point.Name"
                @update:model-value="v => onTypePointName(index, v)"
                hide-details
              >
              </v-text-field>
            </v-col>
            <v-col md="7">
              <v-row no-gutters>
                <v-col md="6">
                  <v-text-field
                    density="compact"
                    variant="outlined"
                    v-model="point.Point.Coord.X"
                    label="X"
                    hint="Expression"
                    class="mr-2"
                    :color="expressionColor"
                  ></v-text-field>
                </v-col>
                <v-col md="6">
                  <v-text-field
                    density="compact"
                    variant="outlined"
                    v-model="point.Point.Coord.Y"
                    label="Y"
                    hint="Expression"
                    :color="expressionColor"
                  ></v-text-field>
                </v-col>
                <v-col md="12">
                  <v-select
                    label="Position de la légende"
                    density="compact"
                    variant="outlined"
                    :items="labelPosItems.map(i => i.text)"
                    :model-value="
                      labelPosItems.find(v => v.value == point.Point.Pos)?.text
                    "
                    @update:model-value="v => onPosChange(v, index)"
                    hide-details
                  ></v-select>
                </v-col>
              </v-row>
            </v-col>
            <v-col md="2" align-self="center">
              <v-btn icon size="x-small" flat @click="deletePoint(index)">
                <v-icon icon="mdi-delete" color="red"></v-icon>
              </v-btn>
            </v-col>
          </v-row>
        </v-list-item>
        <v-divider></v-divider>
      </div>
    </v-list>
  </v-card>

  <v-card color="secondary" class="my-1">
    <v-row no-gutters>
      <v-col align-self="center" md="9">
        <v-card-subtitle> Segments et vecteurs </v-card-subtitle>
      </v-col>
      <v-col md="3" style="text-align: right">
        <v-btn
          icon
          @click="addSegment"
          title="Ajouter un segment ou un vecteur défini par deux points"
          size="x-small"
          class="mr-2 my-2"
          :disabled="segmentsPointItems.length < 2"
        >
          <v-icon icon="mdi-plus" color="green" size="small"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-list>
      <div v-for="(segment, index) in props.modelValue.Drawings.Segments">
        <v-list-item>
          <v-row class="fix-input-width">
            <v-col align-self="center" md="4">
              <v-row>
                <v-col md="12">
                  <v-select
                    density="compact"
                    variant="outlined"
                    hide-details
                    label="Origine"
                    :items="segmentsPointItems"
                    v-model="segment.From"
                  ></v-select>
                </v-col>
                <v-col md="12">
                  <v-select
                    density="compact"
                    variant="outlined"
                    hide-details
                    label="Extrémité"
                    :items="segmentsPointItems"
                    v-model="segment.To"
                  ></v-select>
                </v-col>
              </v-row>
            </v-col>
            <v-col md="7">
              <v-row no-gutters>
                <v-col md="12">
                  <v-switch
                    label="Représenter un vecteur"
                    hide-details
                    v-model="segment.AsVector"
                    color="secondary"
                  ></v-switch>
                </v-col>
                <v-col md="12">
                  <v-text-field
                    density="compact"
                    variant="outlined"
                    label="Légende"
                    hint="Optionnelle"
                    v-model="segment.LabelName"
                  ></v-text-field>
                </v-col>
                <v-col md="12">
                  <v-select
                    label="Position de la légende"
                    density="compact"
                    variant="outlined"
                    :items="labelPosItems.map(i => i.text)"
                    :model-value="
                      labelPosItems.find(v => v.value == segment.LabelPos)?.text
                    "
                    @update:model-value="v => onSegmentLabelPosChange(v, index)"
                    hide-details
                  ></v-select>
                </v-col>
              </v-row>
            </v-col>
            <v-col md="1" align-self="center">
              <v-btn icon size="x-small" flat @click="deleteSegment(index)">
                <v-icon icon="mdi-delete" color="red"></v-icon>
              </v-btn>
            </v-col>
          </v-row>
        </v-list-item>
        <v-divider></v-divider>
      </div>
    </v-list>
  </v-card>

  <v-card color="secondary" class="my-1">
    <v-row no-gutters>
      <v-col align-self="center" md="9">
        <v-card-subtitle> Droites </v-card-subtitle>
      </v-col>
      <v-col md="3" style="text-align: right">
        <v-btn
          icon
          @click="addLine"
          title="Ajouter une droite définie par une équation"
          size="x-small"
          class="mr-2 my-2"
        >
          <v-icon icon="mdi-plus" color="green" size="small"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-list>
      <div v-for="(line, index) in props.modelValue.Drawings.Lines">
        <v-list-item>
          <v-row class="fix-input-width">
            <v-col cols="10">
              <v-row>
                <v-col align-self="center">
                  <v-text-field
                    density="compact"
                    variant="outlined"
                    label="A"
                    hint="Expression du coefficient directeur"
                    v-model="line.A"
                    :color="expressionColor"
                    class="no-hint-padding"
                  ></v-text-field>
                </v-col>
                <v-col align-self="center">
                  <v-text-field
                    density="compact"
                    variant="outlined"
                    label="B"
                    v-model="line.B"
                    hint="Expression de l'ordonnée à l'origine"
                    :color="expressionColor"
                    class="no-hint-padding"
                  ></v-text-field>
                </v-col>
              </v-row>
              <v-row>
                <v-col align-self="center" cols="8">
                  <v-text-field
                    density="compact"
                    variant="outlined"
                    label="Légende"
                    v-model="line.Label"
                  ></v-text-field>
                </v-col>
                <v-col>
                  <btn-color-picker v-model="line.Color"></btn-color-picker>
                </v-col>
              </v-row>
            </v-col>

            <v-col md="2" align-self="center">
              <v-btn icon size="x-small" flat @click="deleteLine(index)">
                <v-icon icon="mdi-delete" color="red"></v-icon>
              </v-btn>
            </v-col>
          </v-row>
        </v-list-item>
        <v-divider></v-divider>
      </div>
    </v-list>
  </v-card>
</template>

<script setup lang="ts">
import { colorByKind, extractPoints } from "@/controller/editor";
import type { FigureBlock, Variable } from "@/controller/exercice_gen";
import { LabelPos, LabelPosLabels, TextKind } from "@/controller/exercice_gen";
import { computed } from "@vue/runtime-core";
import BtnColorPicker from "./BtnColorPicker.vue";

interface Props {
  modelValue: FigureBlock;
  availableParameters: Variable[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: FigureBlock): void;
}>();

const expressionColor = colorByKind[TextKind.Expression];

const labelPosItems = Object.entries(LabelPosLabels).map(k => ({
  value: Number(k[0]) as LabelPos,
  text: k[1]
}));

const segmentsPointItems = computed(() =>
  (props.modelValue.Drawings.Points || []).map(p => p.Name)
);

function addPoint() {
  const points = props.modelValue.Drawings.Points || [];
  points.push({
    Name: "",
    Point: { Coord: { X: "", Y: "" }, Pos: LabelPos.TopRight }
  });
  props.modelValue.Drawings.Points = points;
}

function deletePoint(index: number) {
  props.modelValue.Drawings.Points!.splice(index, 1);
}

function deleteSegment(index: number) {
  props.modelValue.Drawings.Segments!.splice(index, 1);
}

function deleteLine(index: number) {
  props.modelValue.Drawings.Lines!.splice(index, 1);
}

function onPosChange(v: string, index: number) {
  const pos = labelPosItems.find(item => item.text == v)!.value;
  props.modelValue.Drawings.Points![index].Point.Pos = pos;
}

function onSegmentLabelPosChange(v: string, index: number) {
  const pos = labelPosItems.find(item => item.text == v)!.value;
  props.modelValue.Drawings.Segments![index].LabelPos = pos;
}

function addSegment() {
  const points = props.modelValue.Drawings.Points || [];
  if (points.length < 2) {
    return;
  }
  const from = points[0];
  const to = points[1];
  const segments = props.modelValue.Drawings.Segments || [];
  segments.push({
    From: from.Name,
    To: to.Name,
    AsVector: false,
    LabelName: "",
    LabelPos: LabelPos.Top
  });
  props.modelValue.Drawings.Segments = segments;
}

function addLine() {
  const lines = props.modelValue.Drawings.Lines || [];
  lines.push({ Label: "(d)", A: "1", B: "0", Color: "#5500FF" });
  props.modelValue.Drawings.Lines = lines;
}

const availablePoints = computed(() =>
  extractPoints(props.availableParameters)
);

function onTypePointName(index: number, name: string) {
  if (!availablePoints.value.includes(name)) {
    return;
  }

  const point = props.modelValue.Drawings.Points![index];
  // do not autocomplete if fields are already taken
  if (!point.Point.Coord.X) {
    point.Point.Coord.X = "x_" + name;
  }
  if (!point.Point.Coord.Y) {
    point.Point.Coord.Y = "y_" + name;
  }
}
</script>

<style scoped>
.no-hint-padding:deep(.v-input__details) {
  padding-inline: 0px;
}

.fix-input-width:deep(input) {
  width: 100%;
}
</style>
