import { Resend } from 'resend';

export const resend = process.env.RESEND_API_KEY
  ? new Resend(process.env.RESEND_API_KEY)
  : null;

interface SendEmailParams {
  to: string | string[];
  subject: string;
  html: string;
  from?: string;
}

export async function sendEmail({
  to,
  subject,
  html,
  from = 'Koko SaaS <onboarding@resend.dev>',
}: SendEmailParams) {
  if (!resend) {
    console.warn('[EMAIL WARNING] RESEND_API_KEY not configured. Simulated sending email to:', to);
    return { success: true, simulated: true };
  }

  try {
    const data = await resend.emails.send({
      from,
      to,
      subject,
      html,
    });
    return { success: true, data };
  } catch (error) {
    console.error('[EMAIL ERROR] Failed to send email via Resend:', error);
    return { success: false, error };
  }
}
