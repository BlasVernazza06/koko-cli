import { relations } from 'drizzle-orm';
import { user, session, account } from './schema';
import { workspace } from '../workspaces/schema';
import { userActivity, notification } from '../activity/schema';
import { chat } from '../chats/schema';
import { quizAttempt } from '../quizzes/schema';

export const userRelations = relations(user, ({ many }) => ({
  sessions: many(session),
  accounts: many(account),
  workspaces: many(workspace),
  activities: many(userActivity),
  notifications: many(notification),
  chats: many(chat),
  quizAttempts: many(quizAttempt),
}));

export const sessionRelations = relations(session, ({ one }) => ({
  user: one(user, { fields: [session.userId], references: [user.id] }),
}));

export const accountRelations = relations(account, ({ one }) => ({
  user: one(user, { fields: [account.userId], references: [user.id] }),
}));
