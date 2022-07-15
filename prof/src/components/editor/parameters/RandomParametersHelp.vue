<template>
  <v-dialog
    :model-value="props.modelValue"
    @update:model-value="(v) => emit('update:modelValue', v)"
    min-width="1200px"
    width="max-content"
  >
    <v-card
      title="Paramètres aléatoires"
      subtitle="Description des fonctions"
      style="width: 800px; max-height: 70vh"
      class="overflow-y-auto"
    >
      <v-card-text>
        Les paramètres aléatoires sont des variables dont les valeurs sont
        générées à chaque fois qu'une question est posée à l'élève. <br />
        Leur définition se fait avec une syntaxe de type calculatrice. En
        particulier, les fonctions suivantes peuvent être utilisées.

        <v-list color="info" rounded>
          <v-list-item v-for="(content, index) in helpContent" :key="index">
            <v-row>
              <v-col cols="6">
                <v-list-item-title> {{ content[0] }} </v-list-item-title>
              </v-col>
              <v-col align-self="center">
                <div class="text-grey">
                  {{ content[1] }}
                </div>
              </v-col>
            </v-row>
          </v-list-item>
        </v-list>

        <v-alert color="info">
          Les variables peuvent être indicées en ajoutant _ , comme dans x_A.
          Pour insérer un symbol complexe, on peut utiliser la variable spéciale
          @, comme dans @_\ge, qui affichera le code LaTeX \ge, au lieu de
          placer \ge en indice.
        </v-alert>
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

const helpContent = [
  [
    "randChoice(-4;12;99)",
    "Renvoie un nombre aléatoire parmi ceux proposés par l'utilisateur, ici {-4, 12, 99}.",
  ],
  [
    "choiceSymbol( (A; B; C; D); k)",
    "Renvoie le symbol à l'index k parmi ceux proposés (ici, k est entre 1 et 4).",
  ],
  [
    "randSymbol(A; B; C)",
    "Renvoie un symbol choisi uniformément parmi ceux proposés, ici {A, B, C}. \n" +
      "randSymbol est en fait un raccourci pour choiceSymbol( ... ; randInt(1; ...)).",
  ],
  ["randPrime(15;28)", "Renvoie un nombre premier entre 15 et 28 (inclus)."],
  [
    "randDecDen()",
    "Renvoie un entier aléatoire parmi 1, 2, 4, 5, 8, 10, 16, 20, 25, 40, 50, 80, 100 (diviser n'importe quel entier par l'un de ces nombres permettra d'obtenir un nombre décimal)",
  ],
  ["round(x; 3)", "Arrondi x à trois chiffres après la virgule"],
  ["floor(x)", "Renvoie la partie entière de x"],
  ["isPrime(n)", "Renvoie 1 is n est un nombre premier, 0 sinon"],
  ["sgn(x)", "Renvoie le signe de x : 1 si x > 0, -1 si x < 0, 0 si x = 0"],
  ["isZero(x)", "Renvoie 1 si x vaut 0, 0 sinon"],
  ["min(x; 1.2; -4)", "Renvoie le minimum d'une série de valeurs"],
  ["max(x; 1.2; -4)", "Renvoie le maximum d'une série de valeurs"],
  ["exp(x)", "Fonction exponentielle"],
  ["ln(x)", "Fonction logarithme"],
  ["sin(x)", "Fonction sinus"],
  ["cos(x)", "Fonction cosinus"],
  ["tan(x)", "Fonction tangente"],
  ["asin(x)", "Fonction arcsinus"],
  ["acos(x)", "Fonction arccos"],
  ["atan(x)", "Fonction arctan"],
  ["abs(x)", "Fonction valeur absolue"],
  ["sqrt(x)", "Fonction racine carrée"],
] as const;
</script>

<style scoped>
.small-input:deep(input) {
  font-size: 14px;
  width: 100%;
}
</style>
