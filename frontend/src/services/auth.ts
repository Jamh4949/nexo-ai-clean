/**
 * Servicio de autenticación preparado para integrar Firebase Auth o Supabase.
 *
 * Para conectar Firebase:
 *   1. npm install firebase
 *   2. Inicializar firebase app con tu config
 *   3. Usar getAuth(), signInWithPopup(), etc.
 *
 * Para conectar Supabase:
 *   1. npm install @supabase/supabase-js
 *   2. Crear cliente supabase con createClient(url, anonKey)
 *   3. Usar supabase.auth.signInWithOAuth(), etc.
 */

export interface User {
  id: string;
  email: string;
  displayName?: string;
  photoURL?: string;
}

type AuthCallback = (user: User | null) => void;

class AuthService {
  private currentUser: User | null = null;
  private listeners: AuthCallback[] = [];

  getUser(): User | null {
    return this.currentUser;
  }

  onAuthStateChanged(callback: AuthCallback): () => void {
    this.listeners.push(callback);
    callback(this.currentUser);
    return () => {
      this.listeners = this.listeners.filter((cb) => cb !== callback);
    };
  }

  async signIn(): Promise<User> {
    // Simulación — reemplazar con Firebase/Supabase
    const user: User = {
      id: "sim_user_001",
      email: "demo@nexoai.com",
      displayName: "Usuario Demo",
    };

    this.currentUser = user;
    this.notifyListeners();
    return user;
  }

  async signOut(): Promise<void> {
    this.currentUser = null;
    this.notifyListeners();
  }

  private notifyListeners() {
    this.listeners.forEach((cb) => cb(this.currentUser));
  }
}

export const authService = new AuthService();
