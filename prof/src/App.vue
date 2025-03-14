<template>
  <v-app>
    <v-navigation-drawer app v-model="showSideBar" temporary>
      <v-list-item>
        <v-row>
          <v-col cols="auto">
            <v-img :src="logoSrc" width="50px"></v-img>
          </v-col>
          <v-col align-self="center"> Accéder à la page...</v-col>
        </v-row>
      </v-list-item>
      <v-divider></v-divider>

      <v-list-item>
        <v-btn link class="my-1" :to="{ name: 'home' }" color="purple-lighten-4"
          >Accueil</v-btn
        >
      </v-list-item>
      <v-list-item>
        <v-btn
          link
          class="my-1"
          :to="{ name: 'classrooms' }"
          color="purple-lighten-4"
          >Classes et élèves</v-btn
        >
      </v-list-item>
      <v-divider></v-divider>

      <v-list-item>
        <v-btn
          link
          class="my-1"
          :to="{ name: 'editor-question' }"
          color="teal-lighten-4"
          >Editeur de question</v-btn
        >
      </v-list-item>
      <v-list-item>
        <v-btn
          link
          class="my-1"
          :to="{ name: 'editor-exercice' }"
          color="teal-lighten-4"
          >Editeur d'exercice</v-btn
        >
      </v-list-item>

      <v-divider></v-divider>

      <v-list-item>
        <v-btn
          link
          class="my-1"
          :to="{ name: 'trivial' }"
          color="pink-lighten-3"
          >Isy'Triv</v-btn
        >
      </v-list-item>

      <v-list-item>
        <v-btn
          link
          class="my-1"
          :to="{ name: 'homework' }"
          color="pink-lighten-3"
          >Travail à la maison</v-btn
        >
      </v-list-item>
      <v-list-item>
        <v-btn
          link
          class="my-1"
          :to="{ name: 'ceintures' }"
          color="pink-lighten-3"
          >Ceintures de calcul</v-btn
        >
      </v-list-item>
      <v-divider></v-divider>
      <v-list-item>
        <v-btn link class="my-1" :to="{ name: 'reviews' }" color="secondary"
          >Publications</v-btn
        >
      </v-list-item>
      <v-divider></v-divider>

      <v-divider></v-divider>
      <v-list-item>
        <v-btn class="my-1" link :to="{ name: 'settings' }">
          <v-icon icon="mdi-cog" class="mr-2"></v-icon>
          Paramètres</v-btn
        >
      </v-list-item>
      <v-list-item>
        <v-btn class="my-1" @click="logout">
          <v-icon icon="mdi-exit-to-app" class="mr-2"></v-icon>
          Déconnexion</v-btn
        >
      </v-list-item>
    </v-navigation-drawer>

    <v-app-bar
      app
      :density="appBarDense ? 'compact' : undefined"
      color="secondary"
    >
      <v-app-bar-nav-icon
        @click="showSideBar = !showSideBar"
        v-if="isLoggedIn"
      ></v-app-bar-nav-icon>
      <v-app-bar-title tag="h5">
        Isyro -
        <b>{{ $route.meta.Label }}</b>
      </v-app-bar-title>
      <v-spacer></v-spacer>
      <div class="mr-2">
        <small>(Version {{ version }})</small>
      </div>
    </v-app-bar>

    <v-main :class="{ 'background-logo': !isLoggedIn }">
      <loggin-screen v-if="!isLoggedIn" @loggin="onLoggin"></loggin-screen>
      <router-view v-else v-slot="{ Component }">
        <transition>
          <keep-alive>
            <component :is="Component" />
          </keep-alive>
        </transition>
      </router-view>

      <v-snackbar
        style="z-index: 10000"
        app
        :model-value="message != ''"
        @update:model-value="message = ''"
        :timeout="4000"
        :color="messageColor"
        location="bottom left"
        close-on-content-click
      >
        {{ message }}
      </v-snackbar>

      <v-snackbar
        app
        :model-value="errorKind != ''"
        @update:model-value="errorKind = ''"
        :timeout="4000"
        color="red"
      >
        <b>{{ errorKind }}</b>
        <div v-html="errorHtml"></div>
      </v-snackbar>
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
import logoSrc from "@/assets/logo.png";
import { useRoute } from "vue-router";
import { controller, IsDev } from "./controller/controller";
import LogginScreen from "./views/LogginScreen.vue";
import { ref } from "vue";
import { computed } from "vue";

const showSideBar = ref(false);
const version = process.env.VERSION;

const message = ref("");
const messageColor = ref("secondary");

const errorKind = ref("");
const errorHtml = ref("");

controller.onError = (s, m) => {
  errorKind.value = s;
  errorHtml.value = m;
};

controller.showMessage = (s, color) => {
  message.value = s;
  messageColor.value = color || "success";
};

const isLoggedIn = ref(IsDev);

const route = useRoute();
const appBarDense = computed(() => route.name != "home");

function onLoggin() {
  isLoggedIn.value = true;
}

function logout() {
  isLoggedIn.value = false;
  showSideBar.value = false;
  controller.logout();
}
</script>

<style>
.background-logo {
  background: url("@/assets/logo-alpha30.png") no-repeat fixed !important;
  background-position: bottom 10px right 10px !important;
  background-size: 15% !important;
}
</style>
