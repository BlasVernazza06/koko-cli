'use server';

import { db, tenant, userProfile } from '@repo/db';
import { getSession } from '@repo/auth/session';
import { CreateTenantDto } from '@repo/validators';
import { revalidatePath } from 'next/cache';

export const CreateTenant = async (data: CreateTenantDto) => {
  const sessionData = await getSession();
  if (!sessionData || !sessionData.user) {
    throw new Error('Unauthorized');
  }

  const tenantId = crypto.randomUUID();
  const now = new Date();

  // Insert tenant
  await db.insert(tenant).values({
    id: tenantId,
    name: data.name,
    createdAt: now,
    updatedAt: now,
  });

  // Associate user with tenant as ADMIN
  await db.insert(userProfile).values({
    id: crypto.randomUUID(),
    userId: sessionData.user.id,
    tenantId: tenantId,
    role: 'ADMIN',
    createdAt: now,
  });

  revalidatePath('/dashboard');
  return { success: true };
};
