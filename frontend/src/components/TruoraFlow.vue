<script setup lang="ts">
import { onMounted, onUnmounted, ref } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();
const flowStatus = ref<"pending" | "success" | "failed">("pending");

const TRUORA_URL =
  "https://identity.truora.com/?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoiIiwiYWRkaXRpb25hbF9kYXRhIjoie1wiY291bnRyeVwiOlwiQUxMXCIsXCJmbG93X2lkXCI6XCJJUEYyYWI0MjRiY2I1ZTUyNWM2NmM2ODkzOTJiNjViY2FjMlwiLFwicmVkaXJlY3RfdXJsXCI6XCJodHRwczovL3d3dy5nb29nbGUuY29tXCIsXCJwcm9jZXNzX2lkXCI6XCJJRFBmNzU1NTdiOGZlZWVkNjIwNmRhNDc1YTAwMWMzMTdkN1wifSIsImFwcGxpY2F0aW9uX2lkIjoiIiwiY2xpZW50X2lkIjoiVENJNjY5NmRkMDdhMTczZmM2ZGQzOGVmMDQ2OTJkMjg1YmUiLCJleHAiOjE3NzQxMTY0OTUsImdyYW50IjoiZGlnaXRhbC1pZGVudGl0eSIsImlhdCI6MTc3NDExNTU5NSwiaXNzIjoiaHR0cHM6Ly9jb2duaXRvLWlkcC51cy1lYXN0LTEuYW1hem9uYXdzLmNvbS91cy1lYXN0LTFfUmJvQ2lFd01nIiwianRpIjoiNjdiMzQxYzYtYmI4OC00MTU0LTg0NjItYTI3Mzc5Mjk0ZDYzIiwia2V5X25hbWUiOiJrZXlfbmFtZV8xNzc0MTE1NTk1MzQzIiwia2V5X3R5cGUiOiJ0ZXN0IiwidXNlcm5hbWUiOiJUQ0k2Njk2ZGQwN2ExNzNmYzZkZDM4ZWYwNDY5MmQyODViZS1rZXlfbmFtZV8xNzc0MTE1NTk1MzQzIn0.fuHvxhVnc4UBwzlnf3FDYwQ64JnCa_--dNnYGOpQSbk";

function handleMessage(event: MessageEvent) {
  if (event.data === "truora.process.succeeded") {
    console.log("¡Flujo completado con éxito!");
    flowStatus.value = "success";
  } else if (event.data === "truora.process.failed") {
    console.log("El flujo falló");
    flowStatus.value = "failed";
  }
}

onMounted(() => {
  window.addEventListener("message", handleMessage);
});

onUnmounted(() => {
  window.removeEventListener("message", handleMessage);
});
</script>

<template>
  <div class="flex min-h-screen flex-col items-center justify-center bg-gray-50 px-4 py-10">
    <!-- Estado: Pendiente (mostrando iFrame) -->
    <template v-if="flowStatus === 'pending'">
      <div class="mb-6 text-center">
        <div class="mx-auto mb-3 flex h-10 w-10 items-center justify-center rounded-xl bg-primary-600 text-sm font-bold text-white">
          N
        </div>
        <h1 class="text-2xl font-bold tracking-tight text-gray-900">Verificación de identidad</h1>
        <p class="mt-2 text-sm text-gray-500">Completa el proceso para activar tu cuenta</p>
      </div>

      <div class="w-full max-w-[470px]">
        <iframe
          :src="TRUORA_URL"
          allow="camera"
          class="mx-auto block w-full rounded-xl shadow-lg"
          style="height: 700px; max-height: 80vh; border: none"
        ></iframe>
      </div>
    </template>

    <!-- Estado: Éxito -->
    <div v-else-if="flowStatus === 'success'" class="w-full max-w-md text-center">
      <div class="mx-auto mb-8 flex h-24 w-24 items-center justify-center rounded-full bg-green-100">
        <svg class="h-12 w-12 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
        </svg>
      </div>
      <h1 class="mb-3 text-3xl font-bold tracking-tight text-gray-900">
        Verificación completada
      </h1>
      <p class="mb-10 text-lg text-gray-600">
        Tu identidad fue verificada exitosamente. Tu suscripción está activa.
      </p>
      <button
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
          @click="flowStatus = 'pending'"
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
