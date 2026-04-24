"use client"

import { useState } from "react";
import { CredentialFieldProps } from "./CredentialField.type";
import { Check, Copy, Eye, EyeOff } from "react-feather";

export default function CredentialField({ label, value, isSecret }: CredentialFieldProps) {
  const [copied, setCopied] = useState(false);
  const [isVisible, setIsVisible] = useState(!isSecret);

  const handleCopy = () => {
    if (!value) return;
    navigator.clipboard.writeText(value);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const toggleVisibility = () => {
    setIsVisible(!isVisible);
  };

  return (
    <div className="space-y-1.5">
      <label className="block text-xs font-bold text-gray-700">
        {label}
      </label>
      <div className="relative group">
        <input
          type={isVisible ? "text" : "password"}
          value={value}
          readOnly
          className="w-full pl-3 pr-20 py-2.5 rounded-lg border border-gray-200 bg-gray-50 text-sm font-medium text-gray-600 focus:outline-none cursor-default transition-colors group-hover:bg-white group-hover:border-gray-300"
        />
        <div className="absolute right-2 top-1/2 -translate-y-1/2 flex items-center gap-1">
          {isSecret && (
            <button
              onClick={toggleVisibility}
              className="cursor-pointer p-1.5 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-md transition-all"
              title={isVisible ? "Hide secret" : "Show secret"}
              type="button"
            >
              {isVisible ? (
                <EyeOff size={16} />
              ) : (
                <Eye size={16} />
              )}
            </button>
          )}
          <button
            onClick={handleCopy}
            className="cursor-pointer p-1.5 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-md transition-all"
            title="Copy to clipboard"
            type="button"
          >
            {copied ? (
              <Check size={16} className="text-green-500" />
            ) : (
              <Copy size={16} />
            )}
          </button>
        </div>
      </div>
    </div>
  );
};