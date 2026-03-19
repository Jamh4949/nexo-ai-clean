<script setup lang="ts">
import { useRouter } from "vue-router";

const router = useRouter();
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-gray-50 px-6">
    <div class="w-full max-w-md text-center">
      <!-- Ícono check animado -->
      <div class="mx-auto mb-8 flex h-24 w-24 items-center justify-center rounded-full bg-green-100">
        <svg
          class="h-12 w-12 text-green-500"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          stroke-width="2"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M5 13l4 4L19 7"
            class="animate-[draw_0.5s_ease-out_0.3s_both]"
            style="stroke-dasharray: 24; stroke-dashoffset: 24; animation: draw 0.5s ease-out 0.3s forwards"
          />
        </svg>
      </div>

      <h1 class="mb-3 text-3xl font-bold tracking-tight text-gray-900">
        ¡Suscripción exitosa!
      </h1>

      <p class="mb-2 text-lg text-gray-600">
        Bienvenido a <span class="font-semibold text-primary-600">NexoAI</span>
      </p>

      <p class="mb-10 text-sm leading-relaxed text-gray-500">
        Tu pago se procesó correctamente. En breve recibirás un correo
        de confirmación con los detalles de tu suscripción.
      </p>

      <button
        @click="router.push('/')"
        class="inline-flex items-center gap-2 rounded-xl bg-primary-600 px-8 py-3.5 text-base font-semibold text-white shadow-lg shadow-primary-600/25 transition hover:bg-primary-700"
      >
        <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-4 0a1 1 0 01-1-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 01-1 1h-2z" />
        </svg>
        Volver al inicio
      </button>
    </div>
  </div>
</template>

<style scoped>
@keyframes draw {
  to {
    stroke-dashoffset: 0;
  }
}
</style>
