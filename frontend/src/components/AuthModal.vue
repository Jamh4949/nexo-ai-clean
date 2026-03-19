<script setup lang="ts">
import { ref } from "vue";
import {
  loginWithGoogle,
  loginWithGitHub,
  loginWithEmail,
  registerWithEmail,
} from "../services/auth";

const emit = defineEmits<{ close: [] }>();

const isRegister = ref(false);
const email = ref("");
const password = ref("");
const loading = ref(false);
const error = ref("");

function clearError() {
  error.value = "";
}

async function handleGoogle() {
  loading.value = true;
  clearError();
  try {
    await loginWithGoogle();
    emit("close");
  } catch (e: any) {
    error.value = firebaseErrorMsg(e.code);
  } finally {
    loading.value = false;
  }
}

async function handleGitHub() {
  loading.value = true;
  clearError();
  try {
    await loginWithGitHub();
    emit("close");
  } catch (e: any) {
    error.value = firebaseErrorMsg(e.code);
  } finally {
    loading.value = false;
  }
}

async function handleEmailSubmit() {
  if (!email.value || !password.value) return;
  loading.value = true;
  clearError();
  try {
    if (isRegister.value) {
      await registerWithEmail(email.value, password.value);
    } else {
      await loginWithEmail(email.value, password.value);
    }
    emit("close");
  } catch (e: any) {
    error.value = firebaseErrorMsg(e.code);
  } finally {
    loading.value = false;
  }
}

function firebaseErrorMsg(code: string): string {
  const messages: Record<string, string> = {
    "auth/email-already-in-use": "Este correo ya está registrado.",
    "auth/invalid-email": "El correo electrónico no es válido.",
    "auth/user-not-found": "No existe una cuenta con este correo.",
    "auth/wrong-password": "Contraseña incorrecta.",
    "auth/invalid-credential": "Credenciales inválidas. Verifica tu correo y contraseña.",
    "auth/weak-password": "La contraseña debe tener al menos 6 caracteres.",
    "auth/too-many-requests": "Demasiados intentos. Espera un momento e intenta de nuevo.",
    "auth/popup-closed-by-user": "Se cerró la ventana de autenticación.",
    "auth/account-exists-with-different-credential":
      "Ya existe una cuenta con este correo usando otro método de inicio de sesión.",
  };
  return messages[code] || "Ocurrió un error. Intenta de nuevo.";
}
</script>

<template>
  <Teleport to="body">
    <div
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
      @click.self="emit('close')"
    >
      <div
        class="relative w-full max-w-sm rounded-2xl bg-white p-8 shadow-2xl"
        @click.stop
      >
        <!-- Botón cerrar -->
        <button
          @click="emit('close')"
          class="absolute right-4 top-4 rounded-lg p-1 text-gray-400 transition hover:bg-gray-100 hover:text-gray-600"
        >
          <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>

        <!-- Header -->
        <div class="mb-6 text-center">
          <div class="mx-auto mb-3 flex h-10 w-10 items-center justify-center rounded-xl bg-primary-600 text-sm font-bold text-white">
            N
          </div>
          <h2 class="text-xl font-bold text-gray-900">
            {{ isRegister ? "Crear cuenta" : "Iniciar sesión" }}
          </h2>
          <p class="mt-1 text-sm text-gray-500">
            {{ isRegister ? "Regístrate para comenzar" : "Accede a tu cuenta de NexoAI" }}
          </p>
        </div>

        <!-- Error -->
        <div
          v-if="error"
          class="mb-4 rounded-lg bg-red-50 px-4 py-3 text-sm text-red-600"
        >
          {{ error }}
        </div>

        <!-- Botones OAuth -->
        <div class="space-y-3">
          <button
            @click="handleGoogle"
            :disabled="loading"
            class="flex w-full items-center justify-center gap-3 rounded-xl border border-gray-200 bg-white px-4 py-3 text-sm font-medium text-gray-700 transition hover:bg-gray-50 disabled:opacity-60"
          >
            <svg class="h-5 w-5" viewBox="0 0 24 24">
              <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 01-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z" fill="#4285F4"/>
              <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
              <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18A10.96 10.96 0 001 12c0 1.77.42 3.45 1.18 4.93l3.66-2.84z" fill="#FBBC05"/>
              <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
            </svg>
            Continuar con Google
          </button>

          <button
            @click="handleGitHub"
            :disabled="loading"
            class="flex w-full items-center justify-center gap-3 rounded-xl border border-gray-200 bg-gray-900 px-4 py-3 text-sm font-medium text-white transition hover:bg-gray-800 disabled:opacity-60"
          >
            <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12"/>
            </svg>
            Continuar con GitHub
          </button>
        </div>

        <!-- Separador -->
        <div class="my-6 flex items-center gap-3">
          <div class="h-px flex-1 bg-gray-200"></div>
          <span class="text-xs font-medium text-gray-400">O con correo</span>
          <div class="h-px flex-1 bg-gray-200"></div>
        </div>

        <!-- Formulario Email/Password -->
        <form @submit.prevent="handleEmailSubmit" class="space-y-4">
          <div>
            <label class="mb-1.5 block text-sm font-medium text-gray-700">
              Correo electrónico
            </label>
            <input
              v-model="email"
              type="email"
              required
              autocomplete="email"
              placeholder="tu@email.com"
              @input="clearError"
              class="w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm outline-none transition focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20"
            />
          </div>

          <div>
            <label class="mb-1.5 block text-sm font-medium text-gray-700">
              Contraseña
            </label>
            <input
              v-model="password"
              type="password"
              required
              :autocomplete="isRegister ? 'new-password' : 'current-password'"
              placeholder="••••••••"
              @input="clearError"
              class="w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm outline-none transition focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20"
            />
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="w-full rounded-xl bg-primary-600 py-3 text-sm font-semibold text-white transition hover:bg-primary-700 disabled:opacity-60"
          >
            <span v-if="loading" class="inline-flex items-center gap-2">
              <svg class="h-4 w-4 animate-spin" viewBox="0 0 24 24" fill="none">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
              </svg>
              Procesando...
            </span>
            <span v-else>
              {{ isRegister ? "Crear cuenta" : "Iniciar sesión" }}
            </span>
          </button>
        </form>

        <!-- Toggle Login/Register -->
        <p class="mt-6 text-center text-sm text-gray-500">
          {{ isRegister ? "¿Ya tienes cuenta?" : "¿No tienes cuenta?" }}
          <button
            @click="isRegister = !isRegister; clearError()"
            class="ml-1 font-semibold text-primary-600 hover:text-primary-700"
          >
            {{ isRegister ? "Iniciar sesión" : "Crear cuenta" }}
          </button>
        </p>
      </div>
    </div>
  </Teleport>
</template>
