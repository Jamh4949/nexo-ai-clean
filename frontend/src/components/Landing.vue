<script setup lang="ts">
import { ref } from "vue";
import { createCheckout } from "../services/api";
import { authService } from "../services/auth";

const isAnnual = ref(false);
const loading = ref(false);
const email = ref("");
const showEmailModal = ref(false);
const selectedPlan = ref<"monthly" | "annual">("monthly");
const user = ref(authService.getUser());

authService.onAuthStateChanged((u) => {
  user.value = u;
});

function openCheckout(plan: "monthly" | "annual") {
  selectedPlan.value = plan;
  showEmailModal.value = true;
}

async function handleSubscribe() {
  if (!email.value) return;
  loading.value = true;
  try {
    const res = await createCheckout({
      plan: selectedPlan.value,
      email: email.value,
    });
    window.open(res.checkout_url, "_blank");
  } catch (e: any) {
    alert(e.message || "Error al procesar el pago");
  } finally {
    loading.value = false;
    showEmailModal.value = false;
    email.value = "";
  }
}

async function handleAuth() {
  if (user.value) {
    await authService.signOut();
  } else {
    await authService.signIn();
  }
}

const features = [
  {
    icon: "⚡",
    title: "Generación instantánea",
    description:
      "Crea posts, captions y threads en segundos con inteligencia artificial avanzada.",
  },
  {
    icon: "🎯",
    title: "Multi-plataforma",
    description:
      "Contenido optimizado para Instagram, Twitter/X, LinkedIn, TikTok y más.",
  },
  {
    icon: "📊",
    title: "Tono personalizado",
    description:
      "Ajusta el tono, estilo y formato según tu marca y audiencia objetivo.",
  },
  {
    icon: "🔄",
    title: "Calendario de contenido",
    description:
      "Planifica y programa tus publicaciones con un calendario visual integrado.",
  },
];
</script>

<template>
  <div class="min-h-screen bg-white text-gray-900">
    <!-- Navbar -->
    <nav
      class="sticky top-0 z-50 border-b border-gray-100 bg-white/80 backdrop-blur-lg"
    >
      <div class="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
        <div class="flex items-center gap-2">
          <div
            class="flex h-8 w-8 items-center justify-center rounded-lg bg-primary-600 text-sm font-bold text-white"
          >
            N
          </div>
          <span class="text-xl font-bold tracking-tight">NexoAI</span>
        </div>

        <div class="hidden items-center gap-8 text-sm font-medium text-gray-600 md:flex">
          <a href="#features" class="transition hover:text-gray-900">Funciones</a>
          <a href="#pricing" class="transition hover:text-gray-900">Precios</a>
        </div>

        <button
          @click="handleAuth"
          class="rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium transition hover:border-gray-300 hover:bg-gray-50"
        >
          {{ user ? "Cerrar Sesión" : "Iniciar Sesión" }}
        </button>
      </div>
    </nav>

    <!-- Hero -->
    <section class="relative overflow-hidden">
      <div
        class="pointer-events-none absolute inset-0 bg-gradient-to-b from-primary-50/50 to-transparent"
      ></div>
      <div class="relative mx-auto max-w-4xl px-6 pb-20 pt-24 text-center md:pt-32">
        <div
          class="mb-6 inline-flex items-center gap-2 rounded-full border border-primary-200 bg-primary-50 px-4 py-1.5 text-sm font-medium text-primary-700"
        >
          <span class="relative flex h-2 w-2">
            <span
              class="absolute inline-flex h-full w-full animate-ping rounded-full bg-primary-400 opacity-75"
            ></span>
            <span class="relative inline-flex h-2 w-2 rounded-full bg-primary-500"></span>
          </span>
          Potenciado por Inteligencia Artificial
        </div>

        <h1
          class="mb-6 text-4xl font-extrabold leading-tight tracking-tight text-gray-900 md:text-6xl"
        >
          Genera contenido para
          <span
            class="bg-gradient-to-r from-primary-600 to-accent-500 bg-clip-text text-transparent"
          >
            redes sociales
          </span>
          en segundos
        </h1>

        <p class="mx-auto mb-10 max-w-2xl text-lg leading-relaxed text-gray-600 md:text-xl">
          NexoAI usa inteligencia artificial para crear posts, captions y estrategias
          de contenido adaptadas a tu marca. Ahorra horas de trabajo cada semana.
        </p>

        <div class="flex flex-col items-center justify-center gap-4 sm:flex-row">
          <a
            href="#pricing"
            class="inline-flex items-center gap-2 rounded-xl bg-primary-600 px-8 py-3.5 text-base font-semibold text-white shadow-lg shadow-primary-600/25 transition hover:bg-primary-700 hover:shadow-primary-600/30"
          >
            Comenzar gratis
            <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M17 8l4 4m0 0l-4 4m4-4H3"
              />
            </svg>
          </a>
          <a
            href="#features"
            class="inline-flex items-center gap-2 rounded-xl border border-gray-200 px-8 py-3.5 text-base font-semibold text-gray-700 transition hover:bg-gray-50"
          >
            Ver funciones
          </a>
        </div>
      </div>
    </section>

    <!-- Features -->
    <section id="features" class="border-t border-gray-100 bg-gray-50/50 py-24">
      <div class="mx-auto max-w-6xl px-6">
        <div class="mb-16 text-center">
          <h2 class="mb-4 text-3xl font-bold tracking-tight md:text-4xl">
            Todo lo que necesitas para crear contenido
          </h2>
          <p class="mx-auto max-w-2xl text-lg text-gray-600">
            Herramientas diseñadas para creadores, marketers y equipos que quieren
            escalar su presencia en redes sociales.
          </p>
        </div>

        <div class="grid gap-8 md:grid-cols-2">
          <div
            v-for="feature in features"
            :key="feature.title"
            class="group rounded-2xl border border-gray-100 bg-white p-8 transition hover:border-primary-200 hover:shadow-lg hover:shadow-primary-600/5"
          >
            <div
              class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-primary-50 text-2xl transition group-hover:bg-primary-100"
            >
              {{ feature.icon }}
            </div>
            <h3 class="mb-2 text-lg font-semibold">{{ feature.title }}</h3>
            <p class="leading-relaxed text-gray-600">{{ feature.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Pricing -->
    <section id="pricing" class="border-t border-gray-100 py-24">
      <div class="mx-auto max-w-5xl px-6">
        <div class="mb-16 text-center">
          <h2 class="mb-4 text-3xl font-bold tracking-tight md:text-4xl">
            Planes simples y transparentes
          </h2>
          <p class="mx-auto mb-10 max-w-xl text-lg text-gray-600">
            Elige el plan que mejor se adapte a tus necesidades. Sin sorpresas.
          </p>

          <!-- Toggle Mensual / Anual -->
          <div class="inline-flex items-center gap-4 rounded-full bg-gray-100 p-1.5">
            <button
              @click="isAnnual = false"
              :class="[
                'rounded-full px-6 py-2 text-sm font-semibold transition',
                !isAnnual
                  ? 'bg-white text-gray-900 shadow-sm'
                  : 'text-gray-500 hover:text-gray-700',
              ]"
            >
              Mensual
            </button>
            <button
              @click="isAnnual = true"
              :class="[
                'rounded-full px-6 py-2 text-sm font-semibold transition',
                isAnnual
                  ? 'bg-white text-gray-900 shadow-sm'
                  : 'text-gray-500 hover:text-gray-700',
              ]"
            >
              Anual
              <span
                class="ml-1.5 inline-block rounded-full bg-green-100 px-2 py-0.5 text-xs font-semibold text-green-700"
              >
                -17%
              </span>
            </button>
          </div>
        </div>

        <div class="grid items-start gap-8 md:grid-cols-2">
          <!-- Plan Starter -->
          <div
            class="rounded-2xl border border-gray-200 bg-white p-8 transition hover:shadow-lg"
          >
            <div class="mb-6">
              <h3 class="text-lg font-semibold text-gray-900">Starter</h3>
              <p class="mt-1 text-sm text-gray-500">
                Ideal para creadores independientes
              </p>
            </div>

            <div class="mb-6 flex items-baseline gap-1">
              <span class="text-5xl font-extrabold tracking-tight">
                ${{ isAnnual ? "0" : "0" }}
              </span>
              <span class="text-gray-500">/mes</span>
            </div>

            <ul class="mb-8 space-y-3">
              <li
                v-for="item in [
                  '10 generaciones por mes',
                  '2 redes sociales',
                  'Plantillas básicas',
                  'Soporte por email',
                ]"
                :key="item"
                class="flex items-center gap-3 text-sm text-gray-600"
              >
                <svg
                  class="h-5 w-5 shrink-0 text-primary-500"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M5 13l4 4L19 7"
                  />
                </svg>
                {{ item }}
              </li>
            </ul>

            <button
              class="w-full rounded-xl border border-gray-200 py-3 text-sm font-semibold text-gray-700 transition hover:bg-gray-50"
            >
              Comenzar gratis
            </button>
          </div>

          <!-- Plan Pro -->
          <div
            class="relative rounded-2xl border-2 border-primary-600 bg-white p-8 shadow-xl shadow-primary-600/10"
          >
            <div
              class="absolute -top-3.5 left-1/2 -translate-x-1/2 rounded-full bg-primary-600 px-4 py-1 text-xs font-semibold text-white"
            >
              Más popular
            </div>

            <div class="mb-6">
              <h3 class="text-lg font-semibold text-gray-900">Pro</h3>
              <p class="mt-1 text-sm text-gray-500">
                Para equipos y profesionales del marketing
              </p>
            </div>

            <div class="mb-6 flex items-baseline gap-1">
              <span class="text-5xl font-extrabold tracking-tight">
                ${{ isAnnual ? "12.5" : "15" }}
              </span>
              <span class="text-gray-500">/mes</span>
              <span v-if="isAnnual" class="ml-2 text-sm text-gray-400 line-through">
                $15
              </span>
            </div>

            <p v-if="isAnnual" class="mb-4 text-sm text-green-600 font-medium">
              $150/año — Ahorras $30
            </p>

            <ul class="mb-8 space-y-3">
              <li
                v-for="item in [
                  'Generaciones ilimitadas',
                  'Todas las redes sociales',
                  'Calendario de contenido',
                  'Tonos y estilos personalizados',
                  'Análisis de rendimiento',
                  'Soporte prioritario',
                ]"
                :key="item"
                class="flex items-center gap-3 text-sm text-gray-600"
              >
                <svg
                  class="h-5 w-5 shrink-0 text-primary-500"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M5 13l4 4L19 7"
                  />
                </svg>
                {{ item }}
              </li>
            </ul>

            <button
              @click="openCheckout(isAnnual ? 'annual' : 'monthly')"
              :disabled="loading"
              class="w-full rounded-xl bg-primary-600 py-3 text-sm font-semibold text-white shadow-lg shadow-primary-600/25 transition hover:bg-primary-700 disabled:opacity-60"
            >
              {{ loading ? "Procesando..." : "Suscribirse" }}
            </button>
          </div>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="border-t border-gray-100 bg-gray-50/50">
      <div
        class="mx-auto flex max-w-6xl flex-col items-center gap-4 px-6 py-10 text-sm text-gray-500 md:flex-row md:justify-between"
      >
        <div class="flex items-center gap-2 font-semibold text-gray-900">
          <div
            class="flex h-6 w-6 items-center justify-center rounded bg-primary-600 text-xs font-bold text-white"
          >
            N
          </div>
          NexoAI
        </div>
        <p>&copy; {{ new Date().getFullYear() }} NexoAI. Todos los derechos reservados.</p>
      </div>
    </footer>

    <!-- Modal de Email para Checkout -->
    <Teleport to="body">
      <div
        v-if="showEmailModal"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
        @click.self="showEmailModal = false"
      >
        <div class="w-full max-w-md rounded-2xl bg-white p-8 shadow-2xl">
          <h3 class="mb-2 text-xl font-bold">Suscribirse al plan Pro</h3>
          <p class="mb-6 text-sm text-gray-500">
            {{
              selectedPlan === "annual"
                ? "$150/año ($12.5/mes)"
                : "$15/mes"
            }}
          </p>

          <form @submit.prevent="handleSubscribe">
            <label class="mb-1.5 block text-sm font-medium text-gray-700">
              Tu correo electrónico
            </label>
            <input
              v-model="email"
              type="email"
              required
              placeholder="tu@email.com"
              class="mb-6 w-full rounded-lg border border-gray-300 px-4 py-3 text-sm outline-none transition focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20"
            />
            <div class="flex gap-3">
              <button
                type="button"
                @click="showEmailModal = false"
                class="flex-1 rounded-lg border border-gray-200 py-2.5 text-sm font-medium transition hover:bg-gray-50"
              >
                Cancelar
              </button>
              <button
                type="submit"
                :disabled="loading"
                class="flex-1 rounded-lg bg-primary-600 py-2.5 text-sm font-semibold text-white transition hover:bg-primary-700 disabled:opacity-60"
              >
                {{ loading ? "Procesando..." : "Continuar al pago" }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>
  </div>
</template>
