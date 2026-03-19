<script setup lang="ts">
import { useRouter } from "vue-router";

const router = useRouter();
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-gray-50 px-6">
    <div class="w-full max-w-md text-center">
      <!-- Ícono X -->
      <div class="mx-auto mb-8 flex h-24 w-24 items-center justify-center rounded-full bg-amber-100">
        <svg
          class="h-12 w-12 text-amber-500"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          stroke-width="2"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M12 9v3m0 4h.01M5.07 19h13.86a2 2 0 001.74-2.97L13.74 4.03a2 2 0 00-3.48 0L3.33 16.03A2 2 0 005.07 19z"
          />
        </svg>
      </div>

      <h1 class="mb-3 text-3xl font-bold tracking-tight text-gray-900">
        Pago no completado
      </h1>

      <p class="mb-2 text-lg text-gray-600">
        El proceso de pago fue cancelado o no se pudo completar.
      </p>

      <p class="mb-10 text-sm leading-relaxed text-gray-500">
        No se realizó ningún cargo a tu cuenta. Puedes volver
        a intentarlo cuando quieras.
      </p>

      <div class="flex flex-col items-center gap-3 sm:flex-row sm:justify-center">
        <button
          @click="router.push({ path: '/', hash: '#pricing' })"
          class="inline-flex items-center gap-2 rounded-xl bg-primary-600 px-8 py-3.5 text-base font-semibold text-white shadow-lg shadow-primary-600/25 transition hover:bg-primary-700"
        >
          <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Volver a intentar
        </button>

        <button
          @click="router.push('/')"
          class="inline-flex items-center gap-2 rounded-xl border border-gray-200 px-8 py-3.5 text-base font-semibold text-gray-700 transition hover:bg-gray-100"
        >
          Ir al inicio
        </button>
      </div>
    </div>
  </div>
</template>
