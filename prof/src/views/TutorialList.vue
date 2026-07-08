<template>
  <v-card>
    <v-card-text>
      <v-row>
        <v-col class="text-center" v-for="src in srcs" :key="src.id" cols="4">
          <!-- miniature -->
          <v-hover v-if="launchedVideo != src.id">
            <template #default="{ isHovering, props: hoverProps }">
              <v-card v-bind="hoverProps" :subtitle="src.title">
                <v-img
                  :src="`https://tube-action-educative.apps.education.fr/lazy-static/thumbnails/${src.thumbnail}`"
                  cover
                >
                  <v-row class="h-100"
                    ><v-col align-self="center">
                      <v-btn
                        icon="mdi-play"
                        v-show="isHovering"
                        @click="launchedVideo = src.id"
                      >
                      </v-btn> </v-col
                  ></v-row>
                </v-img>
              </v-card>
            </template>
          </v-hover>
          <!-- actual video content -->
          <iframe
            v-else
            loading="lazy"
            :title="src.title"
            width="100%"
            height="315"
            :src="`https://tube-action-educative.apps.education.fr/videos/embed/${src.videoId}?warningTitle=0&amp;peertubeLink=0&autoplay=1`"
            style="border: 0px"
            allow="fullscreen"
            sandbox="allow-same-origin allow-scripts allow-popups allow-forms"
          ></iframe>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { ref } from "vue";

const launchedVideo = ref<number | null>(null);

const srcs = [
  {
    id: 38,
    title: "38. Lier une classe à un autre professeur",
    videoId: "3tqTujZofXVReApmEpozZQ",
    thumbnail: "a9e43881-b54e-4944-ac1b-8a64d87393f7.jpg",
  },
  {
    id: 37,
    title: "37. Limiter le temps sur un devoir noté",
    videoId: "ek2HS3FarEJuGSZZTiQz2q",
    thumbnail: "2729b4d6-622d-4b1c-b109-8db318f451d6.jpg",
  },
  {
    id: 36,
    title: "36. Dupliquer un isy'triv",
    videoId: "kMurkAP1Nrzh4EsekYWtfT",
    thumbnail: "b024f8c5-c870-4160-8d19-ef2b1f155f14.jpg",
  },
  {
    id: 35,
    title: "35. Supprimer un isy'triv",
    videoId: "8jv3ahxMt5WopKkaopU8vA",
    thumbnail: "b3c5ba16-2237-4a37-8d6b-a8304c338a90.jpg",
  },
  {
    id: 34,
    title: "34. Modifier un isy'triv",
    videoId: "9e45bfU6DgHFiezMwWgCjM",
    thumbnail: "23d964fe-5c5f-40f6-bb1b-b5c827d396d8.jpg",
  },
  {
    id: 33,
    title: "33. Attribuer un isy'triv à une classe",
    videoId: "eaRdtj5xFui9EjWxdHLJoK",
    thumbnail: "c0169d53-a2a3-49d9-9642-4c93bb437c59.jpg",
  },
  {
    id: 32,
    title: "32. Arrêter une partie isy'triv en cours",
    videoId: "qRPUp8QJAYtnA12J1Vhq6j",
    thumbnail: "649ae50e-9608-4a29-ae46-462abcd590af.jpg",
  },
  {
    id: 31,
    title: "31. Suivre une partie isy'triv",
    videoId: "kTuPoeDeREa7To2uZPCDAF",
    thumbnail: "05e9e0a8-e327-47ed-a0bb-9358daa4e95f.jpg",
  },
  {
    id: 30,
    title: "30. Lancer un isy'triv en live",
    videoId: "m9nZ9jZxWr9K4oh5W1oUij",
    thumbnail: "29486e7f-2d55-44f2-91c2-ea7b60b2bed7.jpg",
  },
  {
    id: 29,
    title: "29. Créer un isy'triv de révisions",
    videoId: "hNJnVw9zpdSqoCpxxXokAn",
    thumbnail: "f98d0925-b284-4381-80c7-366283e57b0c.jpg",
  },
  {
    id: 28,
    title: "28. Participer à la publication de questions",
    videoId: "cGxk4USEbxU9Jo4hyXUVzX",
    thumbnail: "d69d66ff-95fd-4b7e-bcfc-af818a7d2729.jpg",
  },
  {
    id: 27,
    title: "27. Publier une question créée",
    videoId: "4rJdGbyyK4TYvbnFiK232d",
    thumbnail: "4c7d1396-83bc-465a-9ec8-986627aab093.jpg",
  },
  {
    id: 26,
    title: "26. Supprimer un bloc d'une question",
    videoId: "9VZ4eoNFc38U7mjRAbMM4A",
    thumbnail: "1ad3ceae-cec5-4191-a75d-3007abca0b2e.jpg",
  },
  {
    id: 25,
    title: "25. Dupliquer un bloc dans une question",
    videoId: "83CrVNX6NvffhCPyqs8xA1",
    thumbnail: "26b487cb-23df-4283-8e01-bb4abeac0c21.jpg",
  },
  {
    id: 24,
    title: "24. Supprimer une variante d'une question",
    videoId: "eUx3CCck5ZjT3F9w7PbhTS",
    thumbnail: "3c15caa2-c48d-4241-b468-34ddd4a48071.jpg",
  },
  {
    id: 23,
    title: "23. Créer une variante d'une question",
    videoId: "tYg4BersAqaUXyLcYBHNa6",
    thumbnail: "61ac3d7e-b3a4-4de0-8232-77bf3e922794.jpg",
  },
  {
    id: 22,
    title: "22. Supprimer une question",
    videoId: "oEC6KsXEpA13M4kcHzuEjL",
    thumbnail: "6545f47d-dfd1-4eee-9de9-a822d8e64d2b.jpg",
  },
  {
    id: 21,
    title: "21. Dupliquer une question",
    videoId: "uepzpNg7ktN5JZ2htvnkXt",
    thumbnail: "65ea37f1-719b-4ff2-a6dc-2c30408612e9.jpg",
  },
  {
    id: 20,
    title: "20. Insérer une image dans un énoncé",
    videoId: "wwT5NsVEuxnxkkcKiziFvb",
    thumbnail: "cf80243f-2a51-4736-a73c-53902373505a.jpg",
  },
  {
    id: 19,
    title: "19. Insérer un tableau dans un énoncé",
    videoId: "vd23bsQUNi3T3puTbqFMbQ",
    thumbnail: "6586030b-ba92-4336-80eb-5b6ad674db3a.jpg",
  },
  {
    id: 18,
    title: "18. Créer une réponse texte nombre",
    videoId: "9Dd5Vwjia3X1cAvAuNRUBD",
    thumbnail: "7f05cda3-ac56-458c-bacb-9b734a83d1d4.jpg",
  },
  {
    id: 17,
    title: "17. Créer une réponse expression",
    videoId: "d6J9Q24RRb33VWxLvULNut",
    thumbnail: "68d25794-21ae-4139-b4fc-5d51f677e0fe.jpg",
  },
  {
    id: 16,
    title: "16. Créer un glisser-déposer",
    videoId: "pxk2j27R6r61Fw5rJmcCQB",
    thumbnail: "c6343fd3-9ccf-4d57-a55b-02236fe11ee5.jpg",
  },
  {
    id: 15,
    title: "15. Modifier une question créée",
    videoId: "ntvrGaLoxMj8HDeur48uZE",
    thumbnail: "96ac4141-9f4d-46e8-84ea-f89bbc69f275.jpg",
  },
  {
    id: 14,
    title: "14. Créer un QCM menu déroulant",
    videoId: "uAaEJtvTeHJ4tLnVkAukKc",
    thumbnail: "b6517dd3-5775-4667-b0cf-d7d29b61ed32.jpg",
  },
  {
    id: 13,
    title: "13. Créer une question QCM classique",
    videoId: "tKWp1W6MHcwNEsLzF8MUVk",
    thumbnail: "7211a855-dfab-4fc0-a49c-5a52acf77b9c.jpg",
  },
  {
    id: 12,
    title: "12. Créer une question",
    videoId: "3tgwpok8ACrBxm2prbfmya",
    thumbnail: "bc3991c3-e0e6-4341-8dd5-b35297de934a.jpg",
  },
  {
    id: 11,
    title: "11. Supprimer une classe",
    videoId: "wvXt22XzAZvaQzohRvYiAj",
    thumbnail: "8dbfbfe5-d463-4fd3-a17c-c28041b9f377.jpg",
  },
  {
    id: 10,
    title: "10. Supprimer un travail à la maison",
    videoId: "sTdwJim2H9SX5EQyjBNTDh",
    thumbnail: "d731e186-ccfd-41a3-9312-0bc6c30ec5b7.jpg",
  },
  {
    id: 9,
    title: "9. Supprimer un élève d'une classe",
    videoId: "ehkjsv4GMnzKe2bZnkcb6S",
    thumbnail: "f14a9903-031d-46ce-bafd-313288fa0e2a.jpg",
  },
  {
    id: 8,
    title: "8. Ajouter un élève à une classe",
    videoId: "mEudFAWcAojbDJegsphpYY",
    thumbnail: "fb5b8883-54db-4921-8ad7-e63d5b29f14e.jpg",
  },
  {
    id: 7,
    title: "7. Dispenser un élève d'un travail noté",
    videoId: "79Kxnxw6CJw4U4iMyByatc",
    thumbnail: "807a89cc-2b31-4922-9f85-d0e9a08fc902.jpg",
  },
  {
    id: 6,
    title: "6. Relever les résultats d'un travail noté",
    videoId: "syYgNs6Hk9GvcScijTW4Yk",
    thumbnail: "5a2ebf35-7177-4052-9115-9775cf96d9a6.jpg",
  },
  {
    id: 5,
    title: "5. Préparer un travail à la maison",
    videoId: "8MtyKqoiFpJ7mtBXGs9oNW",
    thumbnail: "769e08e3-80d7-47d0-9e01-f55f000829d3.jpg",
  },
  {
    id: 4,
    title: "4. Rejoindre une classe",
    videoId: "qTx2GagL55XCQw5xjBmE3u",
    thumbnail: "025da569-600d-41c4-ac51-0f48bda0ad5a.jpg",
  },
  {
    id: 3,
    title: "3. Générer un code classe",
    videoId: "f9Saj6DAzciGtUCeKJMvLg",
    thumbnail: "7b1b3b11-ea6a-41a0-93b3-fd6bc43c65e7.jpg",
  },
  {
    id: 2,
    title: "2. Renommer une classe",
    videoId: "iXHvFf4yoJLsj2cYnC669n",
    thumbnail: "e7d55a72-0574-4532-8f4a-42fe78882538.jpg",
  },
  {
    id: 1,
    title: "1. Importer une classe depuis Pronote",
    videoId: "5axiaTK9EdYu57mzjcKhQy",
    thumbnail: "3c8bc316-ea96-41e3-9e32-736a0a652af2.jpg",
  },
].sort((a, b) => a.id - b.id);
</script>

<style></style>
