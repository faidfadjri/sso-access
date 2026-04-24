"use client";

import React, { useState, useEffect } from 'react';
import styles from './LoginForm.module.css';
import Link from 'next/link';
import { Eye, EyeOff, Loader, AlertTriangle } from 'react-feather';
import { useToast } from '@/app/context/ToastContext';
import { AxiosError } from 'axios';
import { ErrorResponse } from '@/app/api/services/response.type';
import { generateCodeChallenge, generateRandomState } from '@/app/libs/auth';
import { localStorageLib } from '@/app/libs/local-storage';
import { generateAuthorizationURL, OauthAuthorizationParams, oauthLogin } from '@/app/api';

export default function LoginForm() {
  const { showToast } = useToast();
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [isValidating, setIsValidating] = useState(true);
  const [appError, setAppError] = useState<string | null>(null);
  
  const [formData, setFormData] = useState({
    email_or_username: '',
    password: '',
    rememberMe: false
  });

  useEffect(() => {
    const checkApp = async () => {
      const urlParams = new URLSearchParams(window.location.search);
      const client_id = urlParams.get("client_id") || process.env.NEXT_PUBLIC_CLIENT_ID;
      const redirect_uri = urlParams.get("redirect_uri") || process.env.NEXT_PUBLIC_REDIRECT_URI;

      if (client_id && redirect_uri) {
        try {
          setIsValidating(true);
        } catch (error) {
           setAppError("Invalid application or configuration. Please contact the administrator.");
        } finally {
          setIsValidating(false);
        }
      } else {
        setAppError("Invalid application or configuration. Please contact the administrator.");
        setIsValidating(false);
      }
    };

    checkApp();
  }, []);

  const togglePasswordVisibility = () => {
    setShowPassword(!showPassword);
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
  };

const handleSubmit = async (e: React.FormEvent) => {
  e.preventDefault();
  setIsLoading(true);

  try {
    const urlParams = new URLSearchParams(window.location.search);

    const code_verifier = generateRandomState();
    localStorageLib.set("code_verifier", code_verifier);
    const code_challenge = await generateCodeChallenge(code_verifier);

    const params: OauthAuthorizationParams = {
        client_id: urlParams.get("client_id") || process.env.NEXT_PUBLIC_CLIENT_ID || "",
        redirect_uri: urlParams.get("redirect_uri") || process.env.NEXT_PUBLIC_REDIRECT_URI || "",
        response_type: urlParams.get("response_type") || "code",
        scope: urlParams.get("scope") || "all",
        state: urlParams.get("state") || code_verifier,
    }

    if(!urlParams.get("client_id")) {
        params.code_challenge = code_challenge;
        params.code_challenge_method = "S256";
    }

    await oauthLogin(formData, params);

    await new Promise(resolve => setTimeout(resolve, 150));

    window.location.href = generateAuthorizationURL(params);

  } catch (error: unknown) {
    const err = error as AxiosError<ErrorResponse>;
    showToast(
      err.response?.data?.message || err.message || "An error occurred",
      "error"
    );
  } finally {
    setIsLoading(false);
  }
};


  if (isValidating) {
      return (
          <div className="flex flex-col items-center justify-center p-8 space-y-4">
              <Loader className="animate-spin text-primary" size={32} />
              <p className="text-gray-500 text-sm">Validating application...</p>
          </div>
      )
  }

  if (appError) {
      return (
          <div className="flex flex-col items-center justify-center p-8 space-y-4 text-center">
              <div className="w-16 h-16 bg-red-50 rounded-full flex items-center justify-center text-red-500 mb-2">
                  <AlertTriangle size={32} />
              </div>
              <h3 className="text-lg font-bold text-gray-900">Access Denied</h3>
              <p className="text-gray-500 text-sm max-w-xs">{appError}</p>
          </div>
      )
  }

  return (
    <div className="w-full">
        <form onSubmit={handleSubmit} className={styles.form}>
        <div className={styles.inputGroup}>
            <label htmlFor="email_or_username" className={styles.label}>
            Email or Username
            </label>
            <div className={styles.inputWrapper}>
            <input
                id="email_or_username"
                name="email_or_username"
                type="text"
                placeholder="example@gmail.com"
                className={styles.input}
                value={formData.email_or_username}
                onChange={handleInputChange}
                required
            />
            </div>
        </div>

        <div className={styles.inputGroup}>
            <label htmlFor="password" className={styles.label}>
            Password
            </label>
            <div className={styles.inputWrapper}>
            <input
                id="password"
                name="password"
                type={showPassword ? "text" : "password"}
                placeholder="Enter password"
                className={styles.input}
                value={formData.password}
                onChange={handleInputChange}
                required
            />
            <button
                type="button"
                className={styles.eyeIcon}
                onClick={togglePasswordVisibility}
                aria-label={showPassword ? "Hide password" : "Show password"}
            >
                {showPassword ? (
                <EyeOff size={20} />
                ) : (
                <Eye size={20} />
                )}
            </button>
            </div>
        </div>

        <button type="submit" className={styles.submitButton} disabled={isLoading}>
            {isLoading ? (
            <div className={styles.loadingWrapper}>
                <Loader className={styles.loader} size={20} />
                <span>Signing in...</span>
            </div>
            ) : (
            "Sign in"
            )}
        </button>

        <div className={styles.footer}>
            Forgot your <Link href="/forgot-username" className={styles.link}>Username</Link> or <Link href="/forgot-password" className={styles.link}>Password</Link>?
            <p className='text-xs mt-2 opacity-70 italic'>version: stable@1.0.0</p>
        </div>
        </form>
    </div>
  );
}