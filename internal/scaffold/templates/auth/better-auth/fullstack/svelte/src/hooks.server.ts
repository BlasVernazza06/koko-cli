import { building } from "/environment";
import { auth } from "@[[ .ProjectName ]]/auth";
import { svelteKitHandler } from "better-auth/svelte-kit";
import type { Handle } from "@sveltejs/kit";

export const handle: Handle = async ({ event, resolve }) => {
  return svelteKitHandler({
    event,
    resolve,
    auth,
    building,
  });
};