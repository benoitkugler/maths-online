<template>
  <v-card class="ma-1 border-red">
    <v-row style="background-color: lightgray">
      <v-col md="8">
        <v-card-subtitle class="py-2">Paramètres aléatoires</v-card-subtitle>
      </v-col>
      <v-spacer></v-spacer>
      <v-col align-self="center" style="text-align: right">
        <v-btn
          icon
          @click="emit('add')"
          title="Ajouter un paramètre"
          size="small"
          class="mr-2"
        >
          <v-icon icon="mdi-plus" color="green" small></v-icon>
        </v-btn>

        <!-- <v-menu offset-y close-on-content-click>
          <template v-slot:activator="{ isActive, props }">
            <v-btn
              small
              icon
              title="Ajouter un paramètre"
              v-on="{ isActive }"
              v-bind="props"
            >
              <v-icon icon="mdi-plus" color="green"></v-icon>
            </v-btn>
          </template>
          <block-bar @add="addBlock"></block-bar>
        </v-menu> -->

        <!-- <v-divider vertical></v-divider>
        <v-btn icon class="mx-2" @click="save" :disabled="!session_id">
          <v-icon icon="mdi-content-save"></v-icon>
        </v-btn> -->
      </v-col>
    </v-row>
    <v-row no-gutters>
      <v-col>
        <v-list>
          <v-list-item v-for="(param, index) in props.parameters" class="pr-0">
            <v-list-item-title>
              <v-text-field
                variant="plain"
                suffix=":"
                density="compact"
                hide-details
                :model-value="String.fromCodePoint(param.variable)"
                @update:model-value="s => onVariableChange(s, index)"
              ></v-text-field>
            </v-list-item-title>
            <v-text-field
              class="ml-2"
              variant="outlined"
              density="compact"
              hide-details
              :model-value="param.expression"
              @update:model-value="s => onExpressionChange(s, index)"
            ></v-text-field>
            <v-btn icon size="small" flat @click="emit('delete', index)">
              <v-icon icon="mdi-delete" color="red"></v-icon>
            </v-btn>
          </v-list-item>
        </v-list>
      </v-col>
    </v-row>
  </v-card>
</template>

<script setup lang="ts">
import type {
  randomParameter,
  randomParameters
} from "@/controller/exercice_gen";

interface Props {
  parameters: randomParameters;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "add"): void;
  (e: "update", index: number, param: randomParameter): void;
  (e: "delete", index: number): void;
}>();

function onVariableChange(s: string, index: number) {
  if (s != "") {
    const param = props.parameters![index];
    param.variable = s.codePointAt(0)!;
    emit("update", index, param);
  }
}

function onExpressionChange(s: string, index: number) {
  const param = props.parameters![index];
  param.expression = s;
  emit("update", index, param);
}
</script>
