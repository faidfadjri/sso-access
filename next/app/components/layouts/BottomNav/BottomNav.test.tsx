import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import BottomNav from './BottomNav';

// Mock next/navigation
vi.mock('next/navigation', () => ({
  usePathname: () => '/',
}));

describe('BottomNav Component', () => {
  it('renders navigation links', () => {
    render(<BottomNav />);
    expect(screen.getByText('Apps')).toBeDefined();
    expect(screen.getByText('Account')).toBeDefined();
  });

  it('highlights active link', () => {
    render(<BottomNav />);
    const params = screen.getByText('Apps').closest('a');
    expect(params).toHaveClass('text-white');
  });
});
