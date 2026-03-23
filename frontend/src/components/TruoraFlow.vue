<script setup lang="ts">
import { onMounted, onUnmounted, ref, computed } from "vue";
import { useRouter, useRoute } from "vue-router";
import { currentUser } from "../services/auth";
import { createCheckout } from "../services/api";

const router = useRouter();
const route = useRoute();

const flowStatus = ref<"loading" | "pending" | "success" | "failed" | "error">("loading");
const truoraURL = ref("");
const errorMsg = ref("");
const redirectingToPayment = ref(false);

const selectedPlan = computed<"monthly" | "annual">(() => {
  return route.query.plan === "annual" ? "annual" : "monthly";
});

const planLabel = computed(() => {
  return selectedPlan.value === "annual" ? "$150/año" : "$15/mes";
});

async function fetchToken() {
  flowStatus.value = "loading";
  errorMsg.value = "";

  try {
    const backendURL = import.meta.env.VITE_BACKEND_URL || "http://localhost:8080";
    const body: Record<string, string> = {};
    if (currentUser.value?.id) {
      body.account_id = currentUser.value.id;
    }

    const res = await fetch(`${backendURL}/api/truora/generate-token`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });

    const data = await res.json();

    if (!res.ok || !data.success) {
      throw new Error(data.error || "Error al generar token de verificación");
    }

    truoraURL.value = data.process_url;
    flowStatus.value = "pending";
  } catch (e: any) {
    errorMsg.value = e.message || "No se pudo conectar con el servidor";
    flowStatus.value = "error";
  }
}

async function redirectToStripe() {
  redirectingToPayment.value = true;

  const email = currentUser.value?.email || "";
  if (!email) {
    flowStatus.value = "success";
    redirectingToPayment.value = false;
    return;
  }

  try {
    const session = await createCheckout({
      plan: selectedPlan.value,
      email,
    });
    window.location.href = session.checkout_url;
  } catch (e: any) {
    errorMsg.value = e.message || "Error al crear sesión de pago";
    flowStatus.value = "error";
    redirectingToPayment.value = false;
  }
}

function handleMessage(event: MessageEvent) {
  if (event.origin !== "https://identity.truora.com") {
    return;
  }

  if (event.data === "truora.process.succeeded") {
    console.log("Verificación exitosa — redirigiendo a Stripe Checkout");
    flowStatus.value = "success";
    redirectToStripe();
  } else if (event.data === "truora.process.failed") {
    console.log("El flujo de verificación falló");
    flowStatus.value = "failed";
  }
}

onMounted(() => {
  window.addEventListener("message", handleMessage);
  fetchToken();
});

onUnmounted(() => {
  window.removeEventListener("message", handleMessage);
});
</script>

<template>
  <div class="flex min-h-screen flex-col items-center justify-center bg-gray-50 px-4 py-10">
    <!-- Estado: Cargando token -->
    <div v-if="flowStatus === 'loading'" class="text-center">
      <svg class="mx-auto mb-4 h-10 w-10 animate-spin text-primary-600" viewBox="0 0 24 24" fill="none">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
      </svg>
      <p class="text-sm text-gray-500">Preparando verificación...</p>
    </div>

    <!-- Estado: Error -->
    <div v-else-if="flowStatus === 'error'" class="w-full max-w-md text-center">
      <div class="mx-auto mb-8 flex h-24 w-24 items-center justify-center rounded-full bg-red-100">
        <svg class="h-12 w-12 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3m0 4h.01M5.07 19h13.86a2 2 0 001.74-2.97L13.74 4.03a2 2 0 00-3.48 0L3.33 16.03A2 2 0 005.07 19z" />
        </svg>
      </div>
      <h1 class="mb-3 text-3xl font-bold tracking-tight text-gray-900">Error de conexión</h1>
      <p class="mb-10 text-lg text-gray-600">{{ errorMsg }}</p>
      <div class="flex flex-col items-center gap-3 sm:flex-row sm:justify-center">
        <button
          @click="fetchToken"
          class="inline-flex items-center gap-2 rounded-xl bg-primary-600 px-8 py-3.5 text-base font-semibold text-white shadow-lg shadow-primary-600/25 transition hover:bg-primary-700"
        >
          Reintentar
        </button>
        <button
          @click="router.push('/')"
          class="inline-flex items-center gap-2 rounded-xl border border-gray-200 px-8 py-3.5 text-base font-semibold text-gray-700 transition hover:bg-gray-100"
        >
          Volver al inicio
        </button>
      </div>
    </div>

    <!-- Estado: iFrame de Truora -->
    <template v-else-if="flowStatus === 'pending'">
      <div class="mb-6 text-center">
        <div class="mx-auto mb-3 flex h-10 w-10 items-center justify-center rounded-xl bg-primary-600 text-sm font-bold text-white">
          N
        </div>
        <h1 class="text-2xl font-bold tracking-tight text-gray-900">Verificación de identidad</h1>
        <p class="mt-2 text-sm text-gray-500">
          Completa el proceso para continuar con tu plan {{ planLabel }}
        </p>
      </div>

      <div class="w-full max-w-[470px]">
        <iframe
          :src="truoraURL"
          allow="camera"
          class="mx-auto block w-full rounded-xl shadow-lg"
          style="height: 700px; max-height: 80vh; border: none"
        ></iframe>
      </div>
    </template>

    <!-- Estado: Éxito — redirigiendo a Stripe -->
    <div v-else-if="flowStatus === 'success'" class="w-full max-w-md text-center">
      <div class="mx-auto mb-8 flex h-24 w-24 items-center justify-center rounded-full bg-green-100">
        <svg
          v-if="redirectingToPayment"
          class="h-10 w-10 animate-spin text-green-500"
          viewBox="0 0 24 24"
          fill="none"
        >
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
        </svg>
        <svg v-else class="h-12 w-12 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
        </svg>
      </div>
      <h1 class="mb-3 text-3xl font-bold tracking-tight text-gray-900">
        Verificación completada
      </h1>
      <p v-if="redirectingToPayment" class="mb-10 text-lg text-gray-600">
        Redirigiendo al pago de tu plan {{ planLabel }}...
      </p>
      <p v-else class="mb-10 text-lg text-gray-600">
        Tu identidad fue verificada exitosamente.
      </p>
      <button
        v-if="!redirectingToPayment"
        @click="router.push('/')"
        class="inline-flex items-center gap-2 rounded-xl bg-primary-600 px-8 py-3.5 text-base font-semibold text-white shadow-lg shadow-primary-600/25 transition hover:bg-primary-700"
      >
        Ir al inicio
      </button>
    </div>

    <!-- Estado: Fallo -->
    <div v-else-if="flowStatus === 'failed'" class="w-full max-w-md text-center">
      <div class="mx-auto mb-8 flex h-24 w-24 items-center justify-center rounded-full bg-red-100">
        <svg class="h-12 w-12 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </div>
      <h1 class="mb-3 text-3xl font-bold tracking-tight text-gray-900">
        Verificación no completada
      </h1>
      <p class="mb-10 text-lg text-gray-600">
        El proceso fue cancelado o no se pudo completar. Puedes intentarlo de nuevo.
      </p>
      <div class="flex flex-col items-center gap-3 sm:flex-row sm:justify-center">
        <button
          @click="fetchToken"
          class="inline-flex items-center gap-2 rounded-xl bg-primary-600 px-8 py-3.5 text-base font-semibold text-white shadow-lg shadow-primary-600/25 transition hover:bg-primary-700"
        >
          Reintentar
        </button>
        <button
          @click="router.push('/')"
          class="inline-flex items-center gap-2 rounded-xl border border-gray-200 px-8 py-3.5 text-base font-semibold text-gray-700 transition hover:bg-gray-100"
        >
          Volver al inicio
        </button>
      </div>
    </div>
  </div>
</template>
