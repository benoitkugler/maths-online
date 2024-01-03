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
        <v-col md="6">
          <v-checkbox
            density="compact"
            label="Afficher l'origine"
            v-model="props.modelValue.ShowOrigin"
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
      <div
        v-for="(point, index) in props.modelValue.Drawings.Points"
        :key="index"
      >
        <v-list-item>
          <v-row class="fix-input-width">
            <v-col md="3" align-self="center" class="mt-2">
              <v-row no-gutters>
                <v-col cols="12">
                  <v-text-field
                    density="compact"
                    variant="outlined"
                    label="Nom"
                    v-model="point.Name"
                    @update:model-value="(v) => onTypePointName(index, v)"
                    :color="expressionColor"
                    hint="Expression."
                  >
                  </v-text-field>
                </v-col>
              </v-row>
              <v-row no-gutters>
                <v-col cols="12" align-self="center">
                  <btn-color-picker
                    v-model="point.Point.Color"
                    @update:model-value="emitUpdate"
                  ></btn-color-picker> </v-col
              ></v-row>
            </v-col>
            <v-col md="7" class="mt-2">
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
                  <label-pos-field v-model="point.Point.Pos"></label-pos-field>
                </v-col>
              </v-row>
            </v-col>
            <v-col md="2" align-self="center" class="px-0">
              <v-btn icon size="x-small" @click="deletePoint(index)">
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
        <v-card-subtitle> Segments, vecteurs et droites </v-card-subtitle>
      </v-col>
      <v-col md="3" style="text-align: right">
        <v-btn
          icon
          @click="addSegment"
          title="Ajouter un segment ou un vecteur défini par deux points"
          size="x-small"
          class="mr-2 my-2"
          :disabled="pointsNamesHints.length < 2"
        >
          <v-icon icon="mdi-plus" color="green" size="small"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-list>
      <div
        v-for="(segment, index) in props.modelValue.Drawings.Segments"
        :key="index"
      >
        <v-list-item>
          <v-row class="fix-input-width">
            <v-col align-self="center" md="4" class="mt-2">
              <v-row>
                <v-col md="12">
                  <v-combobox
                    density="compact"
                    variant="outlined"
                    hide-details
                    hide-no-data
                    label="Origine"
                    :items="pointsNamesHints"
                    v-model="segment.From"
                    :color="expressionColor"
                  ></v-combobox>
                </v-col>
                <v-col md="12">
                  <v-combobox
                    density="compact"
                    variant="outlined"
                    hide-details
                    hide-no-data
                    label="Extrémité"
                    :items="pointsNamesHints"
                    v-model="segment.To"
                    :color="expressionColor"
                  ></v-combobox>
                </v-col>
                <v-col md="12" style="text-align: center">
                  <btn-color-picker v-model="segment.Color"></btn-color-picker>
                </v-col>
              </v-row>
            </v-col>
            <v-col md="7" class="mt-2">
              <v-row no-gutters class="px-0">
                <v-col md="12">
                  <segment-kind-field
                    v-model="segment.Kind"
                  ></segment-kind-field>
                </v-col>
                <v-col md="12">
                  <interpolated-text
                    label="Légende (Optionnelle)"
                    v-model="segment.LabelName"
                  >
                  </interpolated-text>
                </v-col>
                <v-col md="12" class="pt-4">
                  <label-pos-field v-model="segment.LabelPos"></label-pos-field>
                </v-col>
              </v-row>
            </v-col>
            <v-col md="1" align-self="center" class="px-0">
              <v-btn icon size="x-small" @click="deleteSegment(index)">
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
      <div
        v-for="(line, index) in props.modelValue.Drawings.Lines"
        :key="index"
      >
        <v-list-item>
          <v-row class="fix-input-width mt-1">
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
                    hide-details
                  ></v-text-field>
                </v-col>
                <v-col>
                  <btn-color-picker v-model="line.Color"></btn-color-picker>
                </v-col>
              </v-row>
            </v-col>

            <v-col md="2" align-self="center">
              <v-btn icon size="x-small" @click="deleteLine(index)">
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
        <v-card-subtitle> Cercles </v-card-subtitle>
      </v-col>
      <v-col md="3" style="text-align: right">
        <v-btn
          icon
          @click="addCircle"
          title="Ajouter un cercle"
          size="x-small"
          class="mr-2 my-2"
        >
          <v-icon icon="mdi-plus" color="green" size="small"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-list>
      <div
        v-for="(circle, index) in props.modelValue.Drawings.Circles"
        :key="index"
      >
        <v-list-item>
          <v-row class="fix-input-width mt-1">
            <v-col cols="3" align-self="center">
              <interpolated-text
                label="Légende (optionnelle)"
                v-model="circle.Legend"
              >
              </interpolated-text>
            </v-col>
            <v-col>
              <v-row>
                <v-col align-self="center">
                  <v-text-field
                    density="compact"
                    variant="outlined"
                    label="Centre: X"
                    hint="Expression"
                    v-model="circle.Center.X"
                    :color="expressionColor"
                    class="no-hint-padding"
                  ></v-text-field>
                </v-col>
                <v-col align-self="center">
                  <v-text-field
                    density="compact"
                    variant="outlined"
                    label="Centre: Y"
                    v-model="circle.Center.Y"
                    hint="Expression"
                    :color="expressionColor"
                    class="no-hint-padding"
                  ></v-text-field>
                </v-col>
                <v-col align-self="center">
                  <v-text-field
                    density="compact"
                    variant="outlined"
                    label="Rayon"
                    v-model="circle.Radius"
                    hint="Expression"
                    :color="expressionColor"
                    class="no-hint-padding"
                  ></v-text-field>
                </v-col>
              </v-row>

              <v-row class="mb-2">
                <v-col>
                  <div class="text-grey mb-1">
                    <small>Couleur de ligne :</small>
                  </div>
                  <btn-color-picker
                    v-model="circle.LineColor"
                  ></btn-color-picker>
                </v-col>
                <v-col>
                  <div class="text-grey mb-1">
                    <small>Couleur de remplissage :</small>
                  </div>
                  <btn-color-picker
                    v-model="circle.FillColor"
                  ></btn-color-picker>
                </v-col>
              </v-row>
            </v-col>

            <v-col cols="auto" align-self="center">
              <v-btn icon size="x-small" @click="deleteCircle(index)">
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
        <v-card-subtitle> Surfaces colorées </v-card-subtitle>
      </v-col>
      <v-col md="3" style="text-align: right">
        <v-btn
          icon
          @click="addArea"
          title="Ajouter une surface colorée délimitée par des points"
          size="x-small"
          class="mr-2 my-2"
          :disabled="pointsNamesHints.length < 3"
        >
          <v-icon icon="mdi-plus" color="green" size="small"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-list>
      <div
        v-for="(area, index) in props.modelValue.Drawings.Areas"
        :key="index"
      >
        <v-row class="mt-1">
          <v-col cols="2" align-self="center">
            <btn-color-picker v-model="area.Color"></btn-color-picker>
          </v-col>
          <v-col cols="9" align-self="center">
            <expression-list-field
              label="Extrémités"
              hint="Défini la surface à colorier (l'ordre compte)"
              :model-value="area.Points || []"
              @update:model-value="(v) => (area.Points = v)"
            ></expression-list-field>
          </v-col>
          <v-col cols="auto" align-self="center" class="pl-0 pr-0">
            <v-btn icon size="x-small" @click="deleteArea(index)">
              <v-icon icon="mdi-delete" color="red"></v-icon>
            </v-btn>
          </v-col>
        </v-row>
        <v-divider></v-divider>
      </div>
    </v-list>
  </v-card>
</template>

<script setup lang="ts">
import {
  LabelPos,
  SegmentKind,
  TextKind,
  type FigureBlock,
  type Variable,
} from "@/controller/api_gen";
import { colorByKind, extractPoints, lastColorUsed } from "@/controller/editor";
import BtnColorPicker from "../utils/BtnColorPicker.vue";
import ExpressionListField from "../utils/ExpressionListField.vue";
import InterpolatedText from "../utils/InterpolatedText.vue";
import LabelPosField from "../utils/LabelPosField.vue";
import SegmentKindField from "../utils/SegmentKindField.vue";
import { computed } from "vue";

interface Props {
  modelValue: FigureBlock;
  availableParameters: Variable[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: FigureBlock): void;
}>();

function emitUpdate() {
  emit("update:modelValue", props.modelValue);
}

const expressionColor = colorByKind[TextKind.Expression];

const pointsNamesHints = computed(() =>
  (props.modelValue.Drawings.Points || []).map((p) => p.Name)
);

function addPoint() {
  const points = props.modelValue.Drawings.Points || [];
  points.push({
    Name: "",
    Point: {
      Coord: { X: "", Y: "" },
      Pos: LabelPos.TopRight,
      Color: lastColorUsed.color,
    },
  });
  props.modelValue.Drawings.Points = points;
  emitUpdate();
}

function deletePoint(index: number) {
  props.modelValue.Drawings.Points!.splice(index, 1);
  emitUpdate();
}

function addSegment() {
  const from = pointsNamesHints.value[0];
  const to = pointsNamesHints.value[1];
  const segments = props.modelValue.Drawings.Segments || [];
  segments.push({
    From: from,
    To: to,
    Kind: SegmentKind.SKSegment,
    LabelName: "",
    LabelPos: LabelPos.Top,
    Color: lastColorUsed.color,
  });
  props.modelValue.Drawings.Segments = segments;
  emitUpdate();
}

function deleteSegment(index: number) {
  props.modelValue.Drawings.Segments!.splice(index, 1);
  emitUpdate();
}

function addLine() {
  const lines = props.modelValue.Drawings.Lines || [];
  lines.push({ Label: "(d)", A: "1", B: "0", Color: lastColorUsed.color });
  props.modelValue.Drawings.Lines = lines;
  emitUpdate();
}

function deleteLine(index: number) {
  props.modelValue.Drawings.Lines!.splice(index, 1);
  emitUpdate();
}

function addCircle() {
  const circles = props.modelValue.Drawings.Circles || [];
  circles.push({
    Center: { X: "1", Y: "2" },
    Radius: "5",
    LineColor: lastColorUsed.color,
    FillColor: "",
    Legend: "$C_1$",
  });
  props.modelValue.Drawings.Circles = circles;
  emitUpdate();
}

function deleteCircle(index: number) {
  props.modelValue.Drawings.Circles!.splice(index, 1);
  emitUpdate();
}

function addArea() {
  const areas = props.modelValue.Drawings.Areas || [];
  areas.push({
    Points: pointsNamesHints.value.slice(0, 3),
    Color: lastColorUsed.color,
  });
  props.modelValue.Drawings.Areas = areas;
  emitUpdate();
}

function deleteArea(index: number) {
  props.modelValue.Drawings.Areas!.splice(index, 1);
  emitUpdate();
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
