import { createAuthClient } from "better-auth/client";

export const authClient = createAuthClient({
  [[ if ne .Backend "self" ]]
  baseURL: import.meta.env.PUBLIC_SERVER_URL || "http://localhost:3001",
  [[ end ]]
});

export const { signIn, signUp, signOut, useSession } = authClient;