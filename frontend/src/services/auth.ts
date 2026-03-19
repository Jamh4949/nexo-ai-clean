import { ref } from "vue";
import {
  onAuthStateChanged as firebaseOnAuthStateChanged,
  signInWithPopup,
  signInWithEmailAndPassword,
  createUserWithEmailAndPassword,
  signOut as firebaseSignOut,
  GoogleAuthProvider,
  GithubAuthProvider,
  type User as FirebaseUser,
} from "firebase/auth";
import { auth } from "../firebase/config";

export interface User {
  id: string;
  email: string;
  displayName?: string;
  photoURL?: string;
}

function mapFirebaseUser(fbUser: FirebaseUser): User {
  return {
    id: fbUser.uid,
    email: fbUser.email ?? "",
    displayName: fbUser.displayName ?? undefined,
    photoURL: fbUser.photoURL ?? undefined,
  };
}

export const currentUser = ref<User | null>(null);
export const authLoading = ref(true);

firebaseOnAuthStateChanged(auth, (fbUser) => {
  currentUser.value = fbUser ? mapFirebaseUser(fbUser) : null;
  authLoading.value = false;
});

const googleProvider = new GoogleAuthProvider();
const githubProvider = new GithubAuthProvider();

export async function loginWithGoogle(): Promise<User> {
  const result = await signInWithPopup(auth, googleProvider);
  return mapFirebaseUser(result.user);
}

export async function loginWithGitHub(): Promise<User> {
  const result = await signInWithPopup(auth, githubProvider);
  return mapFirebaseUser(result.user);
}

export async function loginWithEmail(
  email: string,
  password: string
): Promise<User> {
  const result = await signInWithEmailAndPassword(auth, email, password);
  return mapFirebaseUser(result.user);
}

export async function registerWithEmail(
  email: string,
  password: string
): Promise<User> {
  const result = await createUserWithEmailAndPassword(auth, email, password);
  return mapFirebaseUser(result.user);
}

export async function logout(): Promise<void> {
  await firebaseSignOut(auth);
}
