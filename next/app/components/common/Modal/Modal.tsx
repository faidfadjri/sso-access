"use client";

import { X } from "react-feather";
import { useState, useEffect, ReactNode } from "react";
import { ModalProps } from "./Modal.type";


export default function Modal({
  isOpen,
  onClose,
  title,
  children,
  maxWidth = "max-w-md",
}: ModalProps) {
  const [isRendered, setIsRendered] = useState(false);
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    if (isOpen) {
      setIsRendered(true);
      const timer = setTimeout(() => {
        setIsVisible(true);
      }, 10);
      return () => clearTimeout(timer);
    } else {
      setIsVisible(false);
      const timer = setTimeout(() => {
        setIsRendered(false);
      }, 150);
      return () => clearTimeout(timer);
    }
  }, [isOpen]);

  if (!isRendered) return null;

  return (
    <div
      className={`fixed inset-0 z-50 overflow-y-auto p-4 transition-all duration-150
        ${
          isVisible
            ? "bg-black/20 backdrop-blur-xs opacity-100 pointer-events-auto"
            : "bg-black/0 backdrop-blur-none opacity-0 pointer-events-none"
        }
      `}
    >
      <div className="flex min-h-full items-center justify-center">
        <div
          className={`bg-white rounded-2xl shadow-xl w-full ${maxWidth} p-8 relative
            transition-all duration-200 ease-out transform
            ${
              isVisible
                ? "opacity-100 scale-100 translate-y-0"
                : "opacity-0 scale-95 translate-y-4"
            }
          `}
        >
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-gray-400 hover:text-gray-600 transition-colors"
        >
          <X size={20} className="cursor-pointer" />
        </button>

        {title && (
          <div className="mb-6">
            <h2 className="text-xl font-bold text-gray-800 mb-1">{title}</h2>
          </div>
        )}

        {children}
        </div>
      </div>
    </div>
  );
}
