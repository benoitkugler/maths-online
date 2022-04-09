<template>
  <v-card class="ma-2 border-red">
    <v-row style="background-color: lightgray">
      <v-col md="9" align-self="center">
        <v-card-subtitle class="py-2">Triplets pythagoriciens</v-card-subtitle>
      </v-col>
      <v-spacer></v-spacer>
      <v-col align-self="center" style="text-align: right">
        <v-btn
          icon
          @click="emit('add')"
          title="Ajouter un triplet"
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
          <v-list-item v-for="(param, index) in props.parameters" class="pr-0">
            <v-row>
              <v-col class="pb-0">
                <v-row>
                  <v-col cols="4">
                    <variable-field
                      label="A"
                      v-model="param.A"
                      @update:model-value="emit('update', index, param)"
                    >
                    </variable-field>
                  </v-col>
                  <v-col cols="4">
                    <variable-field
                      label="B"
                      v-model="param.B"
                      @update:model-value="emit('update', index, param)"
                    >
                    </variable-field>
                  </v-col>
                  <v-col cols="4">
                    <variable-field
                      label="C"
                      v-model="param.C"
                      @update:model-value="emit('update', index, param)"
                    >
                    </variable-field>
                  </v-col>
                </v-row>
                <v-row no-gutters class="pt-3">
                  <v-col cols="12" class="pb-0">
                    <v-text-field
                      class="ml-2"
                      variant="outlined"
                      density="compact"
                      v-model.number="param.Bound"
                      label="Borne"
                      hint="a est environ entre 1 et 2*borne^2"
                      @update:model-value="emit('update', index, param)"
                    ></v-text-field>
                  </v-col>
                </v-row>
              </v-col>
              <v-col cols="3" align-self="center">
                <v-btn
                  icon
                  size="small"
                  flat
                  @click="emit('delete', index)"
                  title="Supprimer ce triplet"
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
import type { PythagorianTriplet } from "@/controller/exercice_gen";
import VariableField from "./VariableField.vue";

interface Props {
  parameters: PythagorianTriplet[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "add"): void;
  (e: "update", index: number, param: PythagorianTriplet): void;
  (e: "delete", index: number): void;
}>();
</script>
