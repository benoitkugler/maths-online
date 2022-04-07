<template>
  <container
    @delete="emit('delete')"
    @swap="(o, t) => emit('swap', o, t)"
    :index="props.Pos.index"
    :nb-blocks="props.Pos.nbBlocks"
  >
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
              type="number"
              v-model.num="props.Data.Bounds.Width"
              label="Largeur"
              hint="Largeur de la figure, en nombre de carreaux"
            ></v-text-field>
          </v-col>
          <v-col md="6">
            <v-text-field
              density="compact"
              variant="outlined"
              type="number"
              v-model.num="props.Data.Bounds.Height"
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
              type="number"
              v-model.num="props.Data.Bounds.Origin.X"
              label="Origine : abscisse"
              hint="Abscisse de l'origine par rapport au coin inférieur gauche"
            ></v-text-field>
          </v-col>
          <v-col md="6">
            <v-text-field
              density="compact"
              variant="outlined"
              type="number"
              v-model.num="props.Data.Bounds.Origin.Y"
              label="Origine : ordonnée"
              hint="Ordonnée de l'origine par rapport au coin inférieur gauche"
            ></v-text-field>
          </v-col>
        </v-row>
        <v-row no-gutters>
          <v-col md="6">
            <v-checkbox
              density="compact"
              label="Afficher la grille"
              v-model="props.Data.ShowGrid"
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
        <div v-for="(point, index) in props.Data.Drawings.Points">
          <v-list-item>
            <v-row>
              <v-col md="3" align-self="center">
                <v-text-field
                  density="compact"
                  variant="outlined"
                  label="Nom"
                  v-model="point.Name"
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
                    ></v-text-field>
                  </v-col>
                  <v-col md="6">
                    <v-text-field
                      density="compact"
                      variant="outlined"
                      v-model="point.Point.Coord.Y"
                      label="Y"
                      hint="Expression"
                    ></v-text-field>
                  </v-col>
                  <v-col md="12">
                    <v-select
                      label="Position de la légende"
                      density="compact"
                      variant="outlined"
                      :items="labelPosItems.map(i => i.text)"
                      :model-value="
                        labelPosItems.find(v => v.value == point.Point.Pos)
                          ?.text
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
        <div v-for="(segment, index) in props.Data.Drawings.Segments">
          <v-list-item>
            <v-row>
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
                        labelPosItems.find(v => v.value == segment.LabelPos)
                          ?.text
                      "
                      @update:model-value="
                        v => onSegmentLabelPosChange(v, index)
                      "
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
        <div v-for="(line, index) in props.Data.Drawings.Lines">
          <v-list-item>
            <v-row>
              <v-col align-self="center" md="4">
                <v-text-field
                  density="compact"
                  variant="outlined"
                  label="A"
                  hint="Expression du coefficient directeur"
                  v-model="line.A"
                ></v-text-field>
              </v-col>
              <v-col align-self="center" md="4">
                <v-text-field
                  density="compact"
                  variant="outlined"
                  label="B"
                  hint="Expression de l'ordonnée à l'origine"
                  v-model="line.B"
                ></v-text-field>
              </v-col>
              <v-col align-self="center" md="3">
                <v-text-field
                  density="compact"
                  variant="outlined"
                  label="Légende"
                  v-model="line.Label"
                ></v-text-field>
              </v-col>
              <v-col md="1" align-self="center">
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
  </container>
</template>

<script setup lang="ts">
import type { FigureBlock } from "@/controller/exercice_gen";
import { LabelPos, LabelPosLabels } from "@/controller/exercice_gen";
import { computed } from "@vue/runtime-core";
import type { ContainerProps } from "./Container.vue";
import Container from "./Container.vue";

interface Props {
  Data: FigureBlock;
  Pos: ContainerProps;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "delete"): void;
  (e: "swap", origin: number, target: number): void;
}>();

const labelPosItems = Object.entries(LabelPosLabels).map(k => ({
  value: Number(k[0]) as LabelPos,
  text: k[1]
}));

const segmentsPointItems = computed(() =>
  (props.Data.Drawings.Points || []).map(p => p.Name)
);

function addPoint() {
  const points = props.Data.Drawings.Points || [];
  points.push({
    Name: "",
    Point: { Coord: { X: "0", Y: "0" }, Pos: LabelPos.Top }
  });
  props.Data.Drawings.Points = points;
}

function deletePoint(index: number) {
  props.Data.Drawings.Points!.splice(index, 1);
}

function deleteSegment(index: number) {
  props.Data.Drawings.Segments!.splice(index, 1);
}

function deleteLine(index: number) {
  props.Data.Drawings.Lines!.splice(index, 1);
}

function onPosChange(v: string, index: number) {
  const pos = labelPosItems.find(item => item.text == v)!.value;
  props.Data.Drawings.Points![index].Point.Pos = pos;
}

function onSegmentLabelPosChange(v: string, index: number) {
  const pos = labelPosItems.find(item => item.text == v)!.value;
  props.Data.Drawings.Segments![index].LabelPos = pos;
}

function addSegment() {
  const points = props.Data.Drawings.Points || [];
  if (points.length < 2) {
    return;
  }
  const from = points[0];
  const to = points[1];
  const segments = props.Data.Drawings.Segments || [];
  segments.push({
    From: from.Name,
    To: to.Name,
    AsVector: false,
    LabelName: "",
    LabelPos: LabelPos.Top
  });
  props.Data.Drawings.Segments = segments;
}

function addLine() {
  const lines = props.Data.Drawings.Lines || [];
  lines.push({ Label: "(d)", A: "1", B: "0" });
  props.Data.Drawings.Lines = lines;
}
</script>

<style scoped></style>
