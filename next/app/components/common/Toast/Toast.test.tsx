import { render, screen, fireEvent, act } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import Toast from './Toast';

describe('Toast Component', () => {
  const mockOnClose = vi.fn();

  it('renders the toast message', () => {
    render(<Toast id="test-id" message="Test Message" onClose={mockOnClose} />);
    expect(screen.getByText('Test Message')).toBeDefined();
  });

  it('calls onClose after duration', () => {
    vi.useFakeTimers();
    render(<Toast id="test-id" message="Test Message" duration={3000} onClose={mockOnClose} />);
    
    // Fast-forward time
    act(() => {
      vi.advanceTimersByTime(3000);
      // Wait for the exit animation (300ms)
      vi.advanceTimersByTime(300);
    });

    expect(mockOnClose).toHaveBeenCalledWith('test-id');

    vi.useRealTimers();
  });

  it('calls onClose when close button is clicked', () => {
    vi.useFakeTimers();
    render(<Toast id="test-id" message="Test Message" onClose={mockOnClose} />);
    
    const closeButton = screen.getByRole('button', { name: /close/i });
    fireEvent.click(closeButton);

    // Wait for exit animation
    act(() => {
      vi.advanceTimersByTime(300);
    });

    expect(mockOnClose).toHaveBeenCalledWith('test-id');
    vi.useRealTimers();
  });

  it('renders correct variant styles', () => {
    const { container } = render(<Toast id="test-id" message="Error" type="error" onClose={mockOnClose} />);
    // Check for a class specific to error variant (border-red-500)
    expect(container.firstChild).toHaveClass('border-red-500');
  });
});
