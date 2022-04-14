import QuestionEditor from "@/views/QuestionEditor.vue";
import TrivialPoursuit from "@/views/TrivialPoursuit.vue";
import { createRouter, createWebHistory } from "vue-router";
import HomeView from "../views/HomeView.vue";

const router = createRouter({
  history: createWebHistory("/prof/"),
  routes: [
    {
      path: "/",
      name: "home",
      component: HomeView,
      meta: { Label: "Accueil" }
    },
    {
      path: "/trivial",
      name: "trivial",
      component: TrivialPoursuit,
      meta: { Label: "Configuration du TrivialPoursuit" }
    },
    {
      path: "/editor",
      name: "editor",
      component: QuestionEditor,
      meta: { Label: "Editeur de question" }
    }
    // {
    //   path: '/about',
    //   name: 'Activités',
    //   // route level code-splitting
    //   // this generates a separate chunk (About.[hash].js) for this route
    //   // which is lazy-loaded when the route is visited.
    //   component: () => import('../views/AboutView.vue')
    // }
  ]
});

export default router;
