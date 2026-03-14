import React from 'react';
import { describe, expect, it, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import App from './App';

vi.mock('./lib/api', () => {
  return {
    getMe: vi.fn(async () => ({ username: 'alice', role: 'admin', display_name: 'Alice' })),
    updateProfile: vi.fn(async () => undefined),
    changePassword: vi.fn(async () => undefined),
    adminListUsers: vi.fn(async () => []),
    adminCreateUser: vi.fn(async () => undefined),
    // Keep dashboard from crashing if a test ever hits it; no tasking coverage here.
    apiBase: vi.fn(() => 'http://localhost:9080'),
    listAgents: vi.fn(async () => []),
    listTasks: vi.fn(async () => []),
    listResults: vi.fn(async () => []),
  };
});

describe('App navigation', () => {
  it('renders shell for authenticated users and navigates to Profile', async () => {
    window.localStorage.setItem('redforge_token', 'test-token');
    window.history.pushState({}, '', '/profile');

    const user = userEvent.setup();
    render(<App />);

    expect(await screen.findByText('alice (admin)', { selector: '.page-sub' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: /Profile/i })).toBeInTheDocument();
    expect(screen.getByText('Change Password')).toBeInTheDocument();

    await user.click(screen.getByRole('link', { name: /Users/i }));
    expect(await screen.findByText(/User Management/i)).toBeInTheDocument();
  });
});
