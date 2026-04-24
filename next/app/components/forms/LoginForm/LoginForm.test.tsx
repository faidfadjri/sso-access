import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import LoginForm from './LoginForm';
import { oauthLogin, generateAuthorizationURL } from '@/app/api';

// Mock the ToastContext
const mockShowToast = vi.fn();
vi.mock('@/app/context/ToastContext', () => ({
  useToast: () => ({
    showToast: mockShowToast,
  }),
}));

vi.mock('@/app/api', async (importOriginal) => {
  const actual: any = await importOriginal();
  return {
    ...actual,
    oauthLogin: vi.fn(),
    generateAuthorizationURL: vi.fn(() => 'http://test.url'),
  };
});

// Mock next/link to behave like a regular anchor
vi.mock('next/link', () => {
    return {
      __esModule: true,
      default: ({ children, href }: { children: React.ReactNode; href: string }) => {
        return <a href={href}>{children}</a>;
      },
    };
  });

describe('LoginForm Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.stubEnv('NEXT_PUBLIC_CLIENT_ID', 'test-client-id');
    vi.stubEnv('NEXT_PUBLIC_REDIRECT_URI', 'http://localhost');
  });

  it('renders login form inputs', async () => {
    render(<LoginForm />);
    expect(await screen.findByLabelText(/email or username/i)).toBeDefined();
    expect(screen.getByLabelText(/^Password$/i)).toBeDefined();
    expect(screen.getByRole('button', { name: /sign in/i })).toBeDefined();
  });

  it('updates input values', async () => {
    render(<LoginForm />);
    const usernameInput = await screen.findByLabelText(/email or username/i) as HTMLInputElement;
    fireEvent.change(usernameInput, { target: { value: 'testuser' } });
    expect(usernameInput.value).toBe('testuser');
  });

  it('toggles password visibility', async () => {
    render(<LoginForm />);
    const passwordInput = await screen.findByLabelText(/^Password$/i) as HTMLInputElement;
    const toggleButton = screen.getByLabelText(/show password/i);

    expect(passwordInput.type).toBe('password');
    
    fireEvent.click(toggleButton);
    expect(passwordInput.type).toBe('text');
    
    fireEvent.click(toggleButton);
    expect(passwordInput.type).toBe('password');
  });

  it('submits form and calls oauthLogin', async () => {
    render(<LoginForm />);
    const usernameInput = await screen.findByLabelText(/email or username/i);
    const passwordInput = screen.getByLabelText(/^Password$/i);
    const submitButton = screen.getByRole('button', { name: /sign in/i });

    fireEvent.change(usernameInput, { target: { value: 'user@example.com' } });
    fireEvent.change(passwordInput, { target: { value: 'password123' } });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(oauthLogin).toHaveBeenCalled();
    });
  });
});
