<template>
  <v-app>
    <v-navigation-drawer app v-model="showSideBar" temporary>
      <v-list-item>
        <v-list-item-title class="title"> Contenu </v-list-item-title>
      </v-list-item>
      <v-divider></v-divider>

      <v-list-item>
        <v-btn link :to="{ name: 'home' }" color="purple-lighten-4"
          >Accueil</v-btn
        >
      </v-list-item>
      <v-list-item>
        <v-btn link :to="{ name: 'classrooms' }" color="purple-lighten-4"
          >Classes et élèves</v-btn
        >
      </v-list-item>
      <v-divider></v-divider>

      <v-list-item>
        <v-btn link :to="{ name: 'trivial' }" color="pink-lighten-3"
          >Triv'Maths</v-btn
        >
      </v-list-item>
      <v-divider></v-divider>

      <v-list-item>
        <v-btn link :to="{ name: 'editor-question' }" color="teal-lighten-4"
          >Editeur de question</v-btn
        >
      </v-list-item>
      <v-list-item>
        <v-btn link :to="{ name: 'editor-exercice' }" color="teal-lighten-4"
          >Editeur d'exercice</v-btn
        >
      </v-list-item>

      <v-divider></v-divider>
      <v-list-item>
        <v-btn @click="logout">
          <v-icon icon="mdi-exit-to-app" class="mr-2"></v-icon>
          Déconnexion</v-btn
        >
      </v-list-item>
    </v-navigation-drawer>

    <v-app-bar app dense color="secondary">
      <v-app-bar-nav-icon
        @click="showSideBar = !showSideBar"
        v-if="isLoggedIn"
      ></v-app-bar-nav-icon>
      <v-app-bar-title tag="h5">
        Isyro -
        <b>{{ $route.meta.Label }}</b>
      </v-app-bar-title>
      <v-spacer></v-spacer>
      <small>(Version {{ version }})</small>
    </v-app-bar>

    <v-main>
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
        top
        right
        absolute
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
import { $ref } from "vue/macros";
import { controller, IsDev } from "./controller/controller";
import LogginScreen from "./views/LogginScreen.vue";

let showSideBar = $ref(false);
const version = process.env.VERSION;

let message = $ref("");
let messageColor = $ref("secondary");

let errorKind = $ref("");
let errorHtml = $ref("");

controller.onError = (s, m) => {
  errorKind = s;
  errorHtml = m;
};

controller.showMessage = (s, color) => {
  message = s;
  messageColor = color || "success";
};

let isLoggedIn = $ref(IsDev);

function onLoggin() {
  isLoggedIn = true;
}

function logout() {
  isLoggedIn = false;
  showSideBar = false;
  controller.logout();
}
</script>
