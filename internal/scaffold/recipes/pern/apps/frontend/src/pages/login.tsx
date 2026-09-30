import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { Card, Flex, Box, Text, Heading, TextField, Button } from '@radix-ui/themes';
import { Mail, Lock, LogIn, AlertCircle } from 'lucide-react';
import api from '../lib/api';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setErrorMsg(null);

    try {
      const response = await api.post('/auth/login', { email, password });
      localStorage.setItem('token', response.data.token);
      localStorage.setItem('user', JSON.stringify(response.data.user));
      navigate('/dashboard');
    } catch (err: any) {
      setErrorMsg(err.response?.data?.error || 'Ocurrió un error al iniciar sesión');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Flex align="center" justify="center" style={{ minHeight: '100vh', background: 'var(--gray-2)' }}>
      <Box style={{ width: '100%', maxWidth: '400px', padding: '16px' }}>
        <Card size="3" style={{ borderRadius: '16px', border: '1px solid var(--gray-5)' }}>
          <Flex direction="column" gap="4">
            <Box style={{ textAlign: 'center' }}>
              <Heading size="6" weight="bold" style={{ marginBottom: '4px' }}>
                ¡Bienvenido de nuevo!
              </Heading>
              <Text size="2" color="gray">
                Ingresa tus credenciales para acceder
              </Text>
            </Box>

            {errorMsg && (
              <Flex gap="2" align="center" style={{ padding: '8px 12px', background: 'var(--red-3)', borderRadius: '8px', color: 'var(--red-11)' }}>
                <AlertCircle size={16} />
                <Text size="2" weight="medium">{errorMsg}</Text>
              </Flex>
            )}

            <form onSubmit={handleSubmit}>
              <Flex direction="column" gap="3">
                <Flex direction="column" gap="1">
                  <Text as="label" size="2" weight="bold" color="gray" htmlFor="email">
                    Correo Electrónico
                  </Text>
                  <TextField.Root
                    id="email"
                    type="email"
                    placeholder="correo@ejemplo.com"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                  >
                    <TextField.Slot>
                      <Mail size={16} />
                    </TextField.Slot>
                  </TextField.Root>
                </Flex>

                <Flex direction="column" gap="1">
                  <Text as="label" size="2" weight="bold" color="gray" htmlFor="password">
                    Contraseña
                  </Text>
                  <TextField.Root
                    id="password"
                    type="password"
                    placeholder="••••••••"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                  >
                    <TextField.Slot>
                      <Lock size={16} />
                    </TextField.Slot>
                  </TextField.Root>
                </Flex>

                <Button type="submit" size="3" style={{ cursor: 'pointer', marginTop: '8px' }} disabled={isLoading}>
                  {isLoading ? 'Cargando...' : (
                    <Flex gap="2" align="center">
                      <LogIn size={18} />
                      Iniciar Sesión
                    </Flex>
                  )}
                </Button>
              </Flex>
            </form>

            <Text size="2" align="center" color="gray">
              ¿No tienes una cuenta?{' '}
              <Link to="/register" style={{ color: 'var(--accent-11)', fontWeight: 'bold', textDecoration: 'none' }}>
                Registrarse
              </Link>
            </Text>
          </Flex>
        </Card>
      </Box>
    </Flex>
  );
}
