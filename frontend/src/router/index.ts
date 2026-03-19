import { createRouter, createWebHistory } from "vue-router";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/",
      name: "landing",
      component: () => import("../components/Landing.vue"),
    },
    {
      path: "/success",
      name: "success",
      component: () => import("../components/Success.vue"),
    },
    {
      path: "/cancel",
      name: "cancel",
      component: () => import("../components/Cancel.vue"),
    },
  ],
  scrollBehavior(_to, _from, savedPosition) {
    return savedPosition || { top: 0 };
  },
});

export default router;
