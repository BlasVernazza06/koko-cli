import LoginForm from '@/components/auth/forms/signIn-form';

export default function LoginPage() {
  return (
    <>
      <LoginForm />
      <div className="flex items-center justify-center gap-2 pt-3">
        <p className="text-sm">¿No tiene una cuenta?</p>
        <a
          href="/auth/register"
          className="text-primary font-medium hover:underline"
        >
          Registrarse
        </a>
      </div>
    </>
  );
}
