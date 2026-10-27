import { createAuthClient } from "better-auth/svelte";

export const authClient = createAuthClient({
  [[ if ne .Backend "self" ]]
  baseURL: import.meta.env.VITE_SERVER_URL || "http://localhost:3001",
  [[ end ]]
});

export const { signIn, signUp, signOut, useSession } = authClient;
