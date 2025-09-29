import ClassroomsListVue from "@/views/ClassroomsList.vue";
import ExerciceEditorVue from "@/views/ExerciceEditor.vue";
import HomeworkActivityVue from "@/views/HomeworkActivity.vue";
import QuestionEditorVue from "@/views/QuestionEditor.vue";
import ReviewListVue from "@/views/ReviewList.vue";
import SettingsPanelVue from "@/views/SettingsPanel.vue";
import TrivialPoursuit from "@/views/TrivialPoursuit.vue";
import HomeViewVue from "@/views/HomeView.vue";
import CeinturesActivityVue from "@/views/CeinturesActivity.vue";
import { createRouter, createWebHistory } from "vue-router";
import ResetPasswordVue from "@/views/ResetPassword.vue";

const router = createRouter({
  history: createWebHistory("/prof/"),
  routes: [
    {
      path: "/",
      name: "home",
      component: HomeViewVue,
      meta: { Label: "Accueil" },
    },
    {
      path: "/classrooms",
      name: "classrooms",
      component: ClassroomsListVue,
      meta: { Label: "Classes et élèves" },
    },
    {
      path: "/editor-question",
      name: "editor-question",
      component: QuestionEditorVue,
      meta: { Label: "Editeur de question" },
    },
    {
      path: "/editor-exercice",
      name: "editor-exercice",
      component: ExerciceEditorVue,
      meta: { Label: "Editeur d'exercices" },
    },
    {
      path: "/trivial",
      name: "trivial",
      component: TrivialPoursuit,
      meta: { Label: "Isy'Triv" },
    },
    {
      path: "/homework",
      name: "homework",
      component: HomeworkActivityVue,
      meta: { Label: "Travail à la maison" },
    },
    {
      path: "/ceintures",
      name: "ceintures",
      component: CeinturesActivityVue,
      meta: { Label: "Ceintures de calcul" },
    },
    {
      path: "/reviews",
      name: "reviews",
      component: ReviewListVue,
      meta: { Label: "Publications" },
    },
    {
      path: "/settings",
      name: "settings",
      component: SettingsPanelVue,
      meta: { Label: "Paramètres" },
    },
    {
      path: "/reset-password",
      name: "resetPassword",
      component: ResetPasswordVue,
      meta: { Label: "Mot de passe oublié" },
    },
    // {
    //   path: '/about',
    //   name: 'Activités',
    //   // route level code-splitting
    //   // this generates a separate chunk (About.[hash].js) for this route
    //   // which is lazy-loaded when the route is visited.
    //   component: () => import('../views/AboutView.vue')
    // },
    {
      path: "/:catchAll(.*)",
      redirect: { name: "home" },
    },
  ],
});

export default router;
