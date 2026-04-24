import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import Sidebar from './Sidebar';

// Mock next/navigation
vi.mock('next/navigation', () => ({
  usePathname: () => '/',
}));

// Mock SidebarContext
vi.mock('./SidebarContext', () => ({
  useSidebar: () => ({
    isCollapsed: false,
    toggleSidebar: vi.fn(),
  }),
}));

// Mock next/image
vi.mock("next/image", () => ({
    default: (props: any) => <img {...props} />,
}));

describe('Sidebar Component', () => {
  it('renders sidebar links', () => {
    render(<Sidebar />);
    expect(screen.getByText('Apps')).toBeDefined();
    expect(screen.getByText('Account')).toBeDefined();
  });

  it('highlights active link', () => {
    render(<Sidebar />);
    const activeLink = screen.getByText('Apps').closest('a');
    expect(activeLink?.className).toMatch(/activeLink/);
  });
});
