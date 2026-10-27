import { z } from 'zod';
export declare const createTenantSchema: z.ZodObject<{
    name: z.ZodString;
}, z.core.$strip>;
export type CreateTenantDto = z.infer<typeof createTenantSchema>;
