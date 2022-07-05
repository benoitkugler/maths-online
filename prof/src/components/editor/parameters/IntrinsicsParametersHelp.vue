<template>
  <v-dialog
    :model-value="props.modelValue"
    @update:model-value="(v) => emit('update:modelValue', v)"
    min-width="1200px"
    width="max-content"
  >
    <v-card
      title="Fonctions spéciales"
      subtitle="Description"
      style="width: 800px"
    >
      <v-card-text>
        Les fonctions spéciales permettent de définir des paramètres aléatoires
        complexes de manière rapide et simple. La syntaxe d'une définition suit
        le format : <br />
        <p class="my-2">
          <i>a,b,c,d,... </i> = <b>fonction</b>(<i>argument1, argument2, ...</i
          >)
        </p>
        Les fonctions utilisables sont les suivantes :
        <v-list color="info" rounded>
          <v-list-item>
            <v-row>
              <v-col cols="6">
                <v-list-item-title
                  >a, b, c = pythagorians(bound)</v-list-item-title
                >
              </v-col>
              <v-col align-self="center">
                <div class="text-grey">
                  Génère trois entiers <i>a</i>,<i>b</i>,<i>c</i> vérifiant a^2
                  + b^2 = c^2. <br />
                  <i>bound</i> est un argument optionnel qui controle le maximum
                  de <i>a</i> par <i>2 bound^2</i>
                </div>
              </v-col>
            </v-row>
          </v-list-item>
          <v-list-item>
            <v-row>
              <v-col cols="6">
                <v-list-item-title>H = projection(A, B, C)</v-list-item-title>
              </v-col>
              <v-col align-self="center">
                <div class="text-grey">
                  Calcule le projeté orthogonal du point <i>A</i> sur
                  <i>(BC)</i>.
                </div>
              </v-col>
            </v-row>
          </v-list-item>
        </v-list>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
interface Props {
  modelValue: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:modelValue", v: boolean): void;
}>();
</script>

<style scoped>
.small-input:deep(input) {
  font-size: 14px;
  width: 100%;
}
</style>
