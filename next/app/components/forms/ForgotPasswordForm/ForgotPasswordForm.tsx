"use client";

import React, { useState } from 'react';
import styles from '../LoginForm/LoginForm.module.css';
import Link from 'next/link';
import { Loader } from 'react-feather';
import { useToast } from '@/app/context/ToastContext';
import { forgotUsernameorPasswordRequest } from '@/app/api/services/oauth/oauth.service';

interface ForgotPasswordFormProps {
  type: "username" | "password";
}

export default function ForgotPasswordForm({ type: initialType }: ForgotPasswordFormProps) {
  const { showToast } = useToast();
  const [type, setType] = useState<"username" | "password">(initialType);
  const [isLoading, setIsLoading] = useState(false);
  const [email, setEmail] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      await forgotUsernameorPasswordRequest({
        email: email,
        forgot_type: type
      });
      showToast(`We've sent recovery instructions to ${email}`, "success");
      setEmail('');
    } catch (error: unknown) {
      const err = error as { response?: { data?: { message?: string } } };
      showToast(err.response?.data?.message || "An error occurred while processing your request", "error");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="w-full">
        <div className="mb-6 text-center">
            <h2 className="text-xl font-bold text-gray-800">
                Recover {type === "username" ? "Username" : "Password"}
            </h2>
            <p className="text-sm text-gray-500 mt-2">
                Enter your email address to receive recovery instructions.
            </p>
        </div>

        <div className="flex justify-center border-b border-gray-200 mb-6">
            <button 
                className={`flex-1 py-2 text-sm font-medium border-b-2 transition-colors ${type === 'username' ? 'border-primary text-primary' : 'border-transparent text-gray-500 hover:text-gray-700'}`}
                onClick={() => setType('username')}
            >
                Forgot Username
            </button>
            <button 
                className={`flex-1 py-2 text-sm font-medium border-b-2 transition-colors ${type === 'password' ? 'border-primary text-primary' : 'border-transparent text-gray-500 hover:text-gray-700'}`}
                onClick={() => setType('password')}
            >
                Forgot Password
            </button>
        </div>

        <form onSubmit={handleSubmit} className={styles.form}>
        <div className={styles.inputGroup}>
            <label htmlFor="email" className={styles.label}>
                Email Address
            </label>
            <div className={styles.inputWrapper}>
            <input
                id="email"
                name="email"
                type="email"
                placeholder="example@gmail.com"
                className={styles.input}
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
            />
            </div>
        </div>

        <button type="submit" className={styles.submitButton} disabled={isLoading}>
            {isLoading ? (
            <div className={styles.loadingWrapper}>
                <Loader className={styles.loader} size={20} />
                <span>Processing...</span>
            </div>
            ) : (
            "Send Recovery Email"
            )}
        </button>

        <div className={styles.footer}>
            Remembered your details? <Link href="/login" className={styles.link}>Back to Login</Link>
        </div>
        </form>
    </div>
  );
}
