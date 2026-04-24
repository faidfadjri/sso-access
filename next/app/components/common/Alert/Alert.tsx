"use client";

import { useEffect, useState, ReactNode } from "react";
import { AlertProps } from "./Alert.type";

export default function Alert({
  isOpen,
  onClose,
  type = "info",
  title,
  message,
  onConfirm,
  confirmText = "Confirm",
  cancelText = "Cancel",
  showCancel = true,
}: AlertProps): ReactNode {
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
      }, 300);
      return () => clearTimeout(timer);
    }
  }, [isOpen]);

  if (!isRendered) return null;

  const colorMap = {
    success: "bg-emerald-400",
    danger: "bg-red-500",
    warning: "bg-amber-400",
    info: "bg-sky-400",
  };

  const buttonColorMap = {
    success: "bg-emerald-500 hover:bg-emerald-600 shadow-emerald-200",
    danger: "bg-red-500 hover:bg-red-600 shadow-red-200",
    warning: "bg-amber-500 hover:bg-amber-600 shadow-amber-200",
    info: "bg-sky-500 hover:bg-sky-600 shadow-sky-200",
  };

  return (
    <div
      className={`fixed inset-0 z-[60] flex items-start justify-center pt-20 p-4 transition-all duration-150
        ${
          isVisible
            ? "bg-black/20 backdrop-blur-xs opacity-100 pointer-events-auto"
            : "bg-black/0 backdrop-blur-none opacity-0 pointer-events-none"
        }
      `}
    >
      <div
        className={`bg-white rounded-xl shadow-xl w-full max-w-md overflow-hidden relative
          transition-all duration-150 ease-out transform
          ${
            isVisible
              ? "opacity-100 scale-100 translate-y-0"
              : "opacity-0 scale-95 -translate-y-4"
          }
        `}
      >
        {/* Top Color Bar */}
        <div className={`h-2 w-full ${colorMap[type]}`} />

        <div className="p-6">
          <h3 className="text-lg font-bold text-gray-900 mb-2 leading-tight">
            {title}
          </h3>
          
          {message && (
            <div className="text-gray-500 text-sm mb-6">
              {message}
            </div>
          )}

          <div className="flex justify-end gap-3">
            {showCancel && (
              <button
                onClick={onClose}
                className="px-4 py-2 cursor-pointer rounded-lg text-sm font-semibold text-gray-600 hover:bg-gray-100 transition-colors"
              >
                {cancelText}
              </button>
            )}
            
            <button
              onClick={() => {
                if (onConfirm) onConfirm();
                else onClose();
              }}
              className={`px-4 py-2 cursor-pointer rounded-lg text-sm font-semibold text-white shadow-lg transition-colors ${buttonColorMap[type]}`}
            >
              {confirmText}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
