import { createAuthClient } from "better-auth/react";

export const authClient = createAuthClient({
  [[ if ne .Backend "self" ]]
  baseURL: process.env.NEXT_PUBLIC_SERVER_URL || "http://localhost:3001",
  [[ end ]]
});

export const { signIn, signUp, signOut, useSession } = authClient;
