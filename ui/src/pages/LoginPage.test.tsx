import React from 'react';
import { describe, expect, it, vi } from 'vitest';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { ToastProvider } from '../components/ToastProvider';
import { LoginPage } from './LoginPage';

vi.mock('../lib/api', () => {
  return {
    login: vi.fn(async (username: string, password: string) => {
      if (username === 'admin' && password === 'good-password') return 'mock-token';
      throw new Error('invalid credentials');
    }),
  };
});

describe('LoginPage', () => {
  it('stores token and navigates on success', async () => {
    const user = userEvent.setup();
    render(
      <ToastProvider>
        <MemoryRouter initialEntries={['/login']}>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/" element={<div>DASH</div>} />
          </Routes>
        </MemoryRouter>
      </ToastProvider>
    );

    await user.clear(screen.getByLabelText(/Password/i));
    await user.type(screen.getByLabelText(/Password/i), 'good-password');
    await user.click(screen.getByRole('button', { name: /login/i }));

    expect(window.localStorage.getItem('redforge_token')).toBe('mock-token');
    expect(await screen.findByText('DASH')).toBeInTheDocument();
  });

  it('shows a toast on failure', async () => {
    const user = userEvent.setup();
    render(
      <ToastProvider>
        <MemoryRouter initialEntries={['/login']}>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
          </Routes>
        </MemoryRouter>
      </ToastProvider>
    );

    await user.clear(screen.getByLabelText(/Password/i));
    await user.type(screen.getByLabelText(/Password/i), 'bad-password');
    await user.click(screen.getByRole('button', { name: /login/i }));

    expect(await screen.findByText(/Login failed/i)).toBeInTheDocument();
  });
});

