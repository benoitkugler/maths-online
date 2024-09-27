<template>
  <v-dialog
    :model-value="props.modelValue"
    @update:model-value="(v) => emit('update:modelValue', v)"
    width="1000px"
  >
    <v-card
      title="Paramètres aléatoires"
      subtitle="Description de la syntaxe"
      style="max-height: 70vh"
      class="overflow-y-auto"
    >
      <v-card-text>
        Les exercices et questions peuvent dépendre de paramètres aléatoires,
        qui sont des variables dont les valeurs sont générées à chaque fois
        qu'une question est posée à l'élève. <br />

        Les paramètres sont définis dans le panneau des Paramètres aléatoires,
        puis utilisés dans les blocs de contenu d'une question.
        <v-expansion-panels class="my-2">
          <v-expansion-panel title="Syntaxe">
            <v-expansion-panel-text>
              Chaque variable est définie par une égalité reliant le nom de la
              variable à une expression, pouvant contenir des nombres, d'autres
              variables, des opérateurs, des fonctions, ... <br />
              Un commentaire explicatif peut etre inséré en le préfixant par
              <C>#</C>.<br />
              Exemple : <br />
              <C># Deuxième solution</C><br /><C>y_2 = y_1+randint(3;15)</C>
              <br /><br />

              Les nombres peuvent utiliser la virgule ou le point comme
              séparateur décimal. Les constantes
              <C>e</C>, <C>π</C> , <C>pi</C> sont utilisables. <br /><br />
              Les variables sont de la forme <C>x</C>, <C>λ</C> ou <C>x_A</C>,
              <C>x_AB</C> pour ajouter un indice. <br /><br />
              En plus des opérations usuelles, les opérateurs modulo
              <C>%</C> et quotient <C>//</C> sont supportés :
              <C>x % 3</C> renvoie le reste de la division euclidienne de x par
              3; <C>x // 3</C> renvoie le quotient (entier) de la division
              euclidienne de x par 3. <br /><br />
              Un symbole spécial, un mot ou un code LaTeX peuvent être insérés
              directement en utilisant des guillemets :
              <C>x = "moyenne"</C> , <C>A = ">"</C> ou <C>B = "\ge"</C>. <br />
              Une matrice est définie par une liste de lignes :
              <C>A = [[1;2;3];[4;5;6]]</C>
            </v-expansion-panel-text>
          </v-expansion-panel>
          <v-expansion-panel title="Fonctions usuelles">
            <v-expansion-panel-text>
              Les fonctions mathématiques usuelles sont supportées, ainsi que
              des fonctions aléatoires : <br />

              <v-list color="info" rounded>
                <v-list-item
                  v-for="(content, index) in fonctionsDesc"
                  :key="index"
                >
                  <v-row>
                    <v-col cols="6">
                      <v-list-item-title>
                        <C>{{ content[0] }}</C>
                      </v-list-item-title>
                    </v-col>
                    <v-col align-self="center">
                      <div class="text-grey">
                        {{ content[1] }}
                      </div>
                    </v-col>
                  </v-row>
                </v-list-item>
              </v-list>
            </v-expansion-panel-text>
          </v-expansion-panel>
          <v-expansion-panel title="Comparaisons">
            <v-expansion-panel-text>
              Il est possible de modéliser des affectations conditionnelles
              grâce aux opérateurs suivants, qui renvoient 1 si la condition est
              vérifiée, 0 sinon. <br />

              <v-list color="info" rounded>
                <v-list-item v-for="(content, index) in compsDesc" :key="index">
                  <v-row>
                    <v-col cols="6">
                      <v-list-item-title>
                        <C>{{ content[0] }}</C>
                      </v-list-item-title>
                    </v-col>
                    <v-col align-self="center">
                      <div class="text-grey">
                        {{ content[1] }}
                      </div>
                    </v-col>
                  </v-row>
                </v-list-item>
              </v-list>
            </v-expansion-panel-text>
          </v-expansion-panel>
          <v-expansion-panel title="Fonctions spéciales">
            <v-expansion-panel-text>
              Les fonctions spéciales permettent de définir des paramètres
              aléatoires complexes de manière rapide et simple. La syntaxe d'une
              définition suit le format : <br />
              <p class="my-2">
                <i>a,b,c,d,... </i> = <b>fonction</b>(<i
                  >argument1, argument2, ...</i
                >)
              </p>
              Les fonctions utilisables sont les suivantes :
              <v-list color="info" rounded>
                <v-list-item>
                  <v-row>
                    <v-col cols="6">
                      <v-list-item-title
                        ><C>a, b, c = pythagorians(bound)</C></v-list-item-title
                      >
                    </v-col>
                    <v-col align-self="center">
                      <div class="text-grey">
                        Génère trois entiers <i>a</i>,<i>b</i>,<i>c</i>
                        vérifiant a^2 + b^2 = c^2. <br />
                        <i>bound</i> est un argument optionnel qui controle le
                        maximum de <i>a</i> par <i>2 bound^2</i>
                      </div>
                    </v-col>
                  </v-row>
                </v-list-item>

                <v-list-item>
                  <v-row>
                    <v-col cols="6">
                      <v-list-item-title
                        ><C>H = projection(A, B, C)</C> ou
                        <C>x_H, y_H = projection(A, B, C)</C></v-list-item-title
                      >
                    </v-col>
                    <v-col align-self="center">
                      <div class="text-grey">
                        Calcule le projeté orthogonal du point <i>A</i> sur
                        <i>(BC)</i>, de coordonnées (<C>x_H</C>, <C>y_H</C>)
                        <br />
                      </div>
                    </v-col>
                  </v-row>
                </v-list-item>

                <v-list-item>
                  <v-row>
                    <v-col cols="6">
                      <v-list-item-title
                        ><C
                          >a, b = number_pair_sum(difficulty)</C
                        ></v-list-item-title
                      >
                    </v-col>
                    <v-col align-self="center">
                      <div class="text-grey">
                        Renvoie deux entiers aléatoires à utiliser dans une
                        addition.
                        <i>difficulty</i> est en entier entre 1 et 5 permettant
                        d'ajuster la difficulté du calcul
                        <br />
                      </div>
                    </v-col>
                  </v-row>
                </v-list-item>
                <v-list-item>
                  <v-row>
                    <v-col cols="6">
                      <v-list-item-title
                        ><C>
                          a, b = number_pair_prod(difficulty)</C
                        ></v-list-item-title
                      >
                    </v-col>
                    <v-col align-self="center">
                      <div class="text-grey">
                        Idem que <i>number_pair_sum</i>, mais pensée pour un
                        produit.
                        <br />
                      </div>
                    </v-col>
                  </v-row>
                </v-list-item>
              </v-list>
            </v-expansion-panel-text>
          </v-expansion-panel>
          <v-expansion-panel title="Instanciation et évaluation">
            <v-expansion-panel-text>
              La résolution des expressions dépend légèrement de leur contexte
              de définition. Une expression définie dans le panneau des
              paramètres aléatoires est systématiquement évaluée. En revanche,
              une expression apparaissant directement dans le contenu d'une
              question ne l'est pas : les variables sont simplement substituées.
              <br />

              Par exemple, si
              <C>a = 5</C> et <C>b = a + 2</C> sont définies comme paramètres,
              l'expression <C>bx</C> sera instantiée (pour l'élève) en
              <C>7x</C>. En revanche, l'expression <C>(a+2)x</C> sera instantiée
              en <C>(5+2)x</C>.
            </v-expansion-panel-text>
          </v-expansion-panel>
        </v-expansion-panels>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import C from "./CodeSpan.vue";

interface Props {
  modelValue: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:modelValue", v: boolean): void;
}>();

const fonctionsDesc = [
  [
    "randInt(-2; 10)",
    "Renvoie un entier aléatoire entre un minimium et un maximum (inclus), ici -2 et 10.",
  ],
  [
    "randChoice(-4;12;99)",
    "Renvoie un nombre (ou une expression) aléatoire parmi ceux proposés par l'utilisateur, ici {-4, 12, 99}.",
  ],
  [
    "choiceFrom(A; B; C; D; k)",
    "Renvoie l'expression à l'index k parmi celles proposées (ici, k est entre 1 et 4).",
  ],
  ["randPrime(15;28)", "Renvoie un nombre premier entre 15 et 28 (inclus)."],
  [
    "randDecDen(min; max)",
    "Renvoie un entier aléatoire entre min et max tel que diviser n'importe quel autre entier par ce nombre donne un nombre décimal (min et max sont optionnels et valent 1 et 100 par défaut).",
  ],
  [
    "randMatrix(n, p, min, max)",
    "Renvoie une matrice de taille n x p à coefficients entiers aléatoires compris entre min et max (inclus).",
  ],
  ["round(x; 3)", "Arrondi x à trois chiffres après la virgule."],
  [
    "forceDecimal(x)",
    "Affiche x sous forme décimale, même si x est rationnel.",
  ],
  ["floor(x)", "Renvoie la partie entière de x."],
  ["binom(k; n)", "Renvoie le coefficient binomial k parmi n."],
  ["isPrime(n)", "Renvoie 1 is n est un nombre premier, 0 sinon."],
  ["sgn(x)", "Renvoie le signe de x : 1 si x > 0, -1 si x < 0, 0 si x = 0."],
  ["min(x; 1.2; -4)", "Renvoie le minimum d'une série de valeurs."],
  ["max(x; 1.2; -4)", "Renvoie le maximum d'une série de valeurs."],
  [
    `sum(k; 1; n; k^2; "expand")`,
    `Renvoie la somme d'une expression évaluée pour une plage d'indice. Le dernier paramètre est optionel : s'il est présent, le symbole Sigma n'est pas utilisé.
    Utiliser "expand-eval" pour aussi évaluer chaque terme.`,
  ],
  [`prod(k; 1; n; k^2; "expand")`, "Idem que sum, mais pour un produit."],
  [`union(k; 1; n; A_{k}; "expand")`, "Idem que sum, mais pour une union."],
  [
    `inter(k; 1; n; A_{k}; "expand")`,
    "Idem que sum, mais pour une intersection.",
  ],
  ["coeff(A; i; j)", "Renvoie le coefficient en case (i, j) de la matrice A."],
  ["set(A; i; j, v)", "Renvoie la matrice A avec en case (i, j) la valeur v."],
  ["trans(A) ou transpose(A)", "Renvoie la transposée de A."],
  ["trace(A)", "Renvoie la trace de A."],
  ["det(A)", "Renvoie le déterminant de A."],
  ["inv(A)", "Renvoie l'inverse de A."],
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

const compsDesc = [
  ["(a == b)", "a et b sont égaux"],
  ["(a > b)", "a est strictement supérieur à b"],
  ["(a >= b)", "a est supérieur ou égal à b"],
  ["(a < b)", "a est strictement inférieur à b"],
  ["(a <= b)", "a est inférieur ou égal à b"],
] as const;
</script>

<style scoped></style>
