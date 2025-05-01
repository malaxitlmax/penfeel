import React, { useState } from 'react';
import {
  Box,
  Button,
  Container,
  Heading,
  Input,
  Stack,
  Text,
} from '@chakra-ui/react';

const LoginPage: React.FC = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // TODO: Добавить логику авторизации
    console.log('Login attempt:', { email, password });
  };

  return (
    <Container maxW="container.sm" py={10}>
      <Box
        p={8}
        borderWidth={1}
        borderRadius="lg"
        boxShadow="lg"
        bg="white"
      >
        <Stack gap={6} as="form" onSubmit={handleSubmit}>
          <Heading size="lg">Вход в систему</Heading>
          
          <Box>
            <Text mb={2}>Email</Text>
            <Input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="Введите ваш email"
              required
            />
          </Box>

          <Box>
            <Text mb={2}>Пароль</Text>
            <Input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Введите ваш пароль"
              required
            />
          </Box>

          <Button
            type="submit"
            colorScheme="blue"
            width="full"
            size="lg"
          >
            Войти
          </Button>
        </Stack>
      </Box>
    </Container>
  );
};

export default LoginPage; 