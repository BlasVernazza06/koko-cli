import { NextResponse } from 'next/server';
import { db, userProfile, tenant } from '@repo/db';
import { getSession } from '@repo/auth/session';
import { eq } from 'drizzle-orm';

export async function GET() {
  const sessionData = await getSession();
  if (!sessionData || !sessionData.user) {
    return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
  }

  // Get user profile listings associated with tenants
  const profiles = await db
    .select({
      id: tenant.id,
      name: tenant.name,
    })
    .from(userProfile)
    .innerJoin(tenant, eq(userProfile.tenantId, tenant.id))
    .where(eq(userProfile.userId, sessionData.user.id));

  return NextResponse.json(profiles);
}
