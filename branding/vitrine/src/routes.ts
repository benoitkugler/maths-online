import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";
import HomePageVue from "./views/HomePage.vue";
import PrivacyPageVue from "./views/PrivacyPage.vue";

declare module "vue-router" {
  interface RouteMeta {
    title: string;
  }
}

const routes: RouteRecordRaw[] = [
  {
    name: "home",
    path: "/",
    component: HomePageVue,
    meta: {
      title: "Isyro"
    }
  },
  {
    path: "/privacy",
    component: PrivacyPageVue,
    meta: {
      title: "Isyro - Mentions légales et RGPD"
    }
  },
  {
    path: "/:pathMatch(.*)*",
    redirect: { name: "home" }
  }
];

export const router = createRouter({
  // 4. Provide the history implementation to use. We are using the hash history for simplicity here.
  history: createWebHistory(),
  routes // short for `routes: routes`
});
