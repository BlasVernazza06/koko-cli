import SignUpForm from '@/components/auth/forms/signUp-form';

export default function RegisterPage() {
  return (
    <>
      <SignUpForm />
      <div className="flex items-center justify-center gap-2 pt-3">
        <p className="text-sm">¿Ya tiene una cuenta?</p>
        <a
          href="/auth/login"
          className="text-primary font-medium hover:underline"
        >
          Iniciar Sesión
        </a>
      </div>
    </>
  );
}
