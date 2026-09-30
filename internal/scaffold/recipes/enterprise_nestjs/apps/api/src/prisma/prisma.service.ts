import { Injectable, OnModuleInit, OnModuleDestroy } from '@nestjs/common';
import { PrismaClient } from '@prisma/client';

@Injectable()
export class PrismaService extends PrismaClient implements OnModuleInit, OnModuleDestroy {
  async onModuleInit() {
    try {
      await this.$connect();
      console.log('🍃 [PrismaService] Connected to PostgreSQL database');
    } catch (err) {
      console.warn('⚠️ [PrismaService] Database offline, starting with fallback:', err);
    }
  }

  async onModuleDestroy() {
    await this.$disconnect();
  }
}
