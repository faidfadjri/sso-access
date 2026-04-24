"use client";

import { useEffect, useState } from "react";
import { CheckCircle, Info, AlertCircle, XCircle, Loader, X } from "react-feather";

export type ToastType = "success" | "info" | "warning" | "error" | "loading";

export interface ToastProps {
  id: string;
  message: string;
  type?: ToastType;
  duration?: number;
  onClose: (id: string) => void;
}

const variants = {
  success: {
    icon: CheckCircle,
    styles: "border-green-500 text-green-700 bg-white shadow-green-100",
    iconColor: "text-green-500",
  },
  info: {
    icon: Info,
    styles: "border-blue-500 text-blue-700 bg-white shadow-blue-100",
    iconColor: "text-blue-500",
  },
  warning: {
    icon: AlertCircle,
    styles: "border-yellow-500 text-yellow-700 bg-white shadow-yellow-100",
    iconColor: "text-yellow-500",
  },
  error: {
    icon: XCircle,
    styles: "border-red-500 text-red-700 bg-white shadow-red-100",
    iconColor: "text-red-500",
  },
  loading: {
    icon: Loader,
    styles: "border-blue-400 text-blue-600 bg-white shadow-blue-50",
    iconColor: "text-blue-400 animate-spin",
  },
};

export default function Toast({ id, message, type = "info", duration = 3000, onClose }: ToastProps) {
  const [isVisible, setIsVisible] = useState(false);
  const variant = variants[type];
  const Icon = variant.icon;

  useEffect(() => {
    // Trigger enter animation
    requestAnimationFrame(() => setIsVisible(true));

    if (duration && duration > 0 && type !== 'loading') {
      const timer = setTimeout(() => {
        handleClose();
      }, duration);
      return () => clearTimeout(timer);
    }
  }, [duration, type]);

  const handleClose = () => {
    setIsVisible(false);
    // Wait for exit animation to finish before removing from DOM
    setTimeout(() => {
      onClose(id);
    }, 300); // Match transition duration
  };

  return (
    <div
      className={`
        flex items-center w-full max-w-sm p-4 gap-2 mb-3 text-sm font-medium rounded-lg shadow-lg border-l-4 transition-all duration-300 ease-in-out transform
        ${variant.styles}
        ${isVisible ? "translate-x-0 opacity-100" : "translate-x-full opacity-0"}
      `}
      role="alert"
    >
      <div className={`inline-flex items-center justify-center flex-shrink-0 w-8 h-8 ${variant.iconColor}`}>
        <Icon size={20} />
      </div>
      <div className="ml-3 text-sm font-normal flex-1">{message}</div>
      <button
        type="button"
        className={`ml-auto -mx-1.5 -my-1.5 rounded-lg cursor-pointer p-1.5 inline-flex h-8 w-8 text-gray-400 hover:text-gray-900 hover:bg-gray-100 transition-colors hover:bg-transparent`}
        aria-label="Close"
        onClick={handleClose}
      >
        <span className="sr-only">Close</span>
        <X size={16} />
      </button>
    </div>
  );
}
