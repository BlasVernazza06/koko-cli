import { betterAuth } from "better-auth";
[[ if eq .ORM "drizzle" ]]
import { drizzleAdapter } from "better-auth/adapters/drizzle";
import { db } from "@[[ .ProjectName ]]/db";

export const auth = betterAuth({
  database: drizzleAdapter(db, {
    provider: "[[ if eq .Database "sqlite" ]]sqlite[[ else if eq .Database "mysql" ]]mysql[[ else ]]pg[[ end ]]",
  }),
  emailAndPassword: {
    enabled: true,
  },
});
[[ else if eq .ORM "prisma" ]]
import { prismaAdapter } from "better-auth/adapters/prisma";
import { prisma } from "@[[ .ProjectName ]]/db";

export const auth = betterAuth({
  database: prismaAdapter(prisma, {
    provider: "[[ if eq .Database "sqlite" ]]sqlite[[ else if eq .Database "mysql" ]]mysql[[ else ]]postgresql[[ end ]]",
  }),
  emailAndPassword: {
    enabled: true,
  },
});
[[ else ]]
export const auth = betterAuth({
  emailAndPassword: {
    enabled: true,
  },
});
[[ end ]]

export type Auth = typeof auth;
