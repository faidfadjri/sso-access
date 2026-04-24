"use client";

import React, { useState, useEffect, Suspense } from 'react';
import styles from '../LoginForm/LoginForm.module.css';
import { Eye, EyeOff, Loader, AlertTriangle } from 'react-feather';
import { useToast } from '@/app/context/ToastContext';
import { useRouter, useSearchParams } from 'next/navigation';
import { resetPasswordRequest } from '@/app/api/services/oauth/oauth.service';
import Modal from '@/app/components/common/Modal/Modal';

function CreateNewPasswordFormContent() {
  const { showToast } = useToast();
  const router = useRouter();
  const searchParams = useSearchParams();
  const token = searchParams.get('token');
  
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [isSuccessModalOpen, setIsSuccessModalOpen] = useState(false);
  
  useEffect(() => {
    if (!token) {
      router.push('/login');
    }
  }, [token, router]);

  const [formData, setFormData] = useState({
    password: '',
    confirmPassword: ''
  });

  const togglePasswordVisibility = () => {
    setShowPassword(!showPassword);
  };
  
  const toggleConfirmPasswordVisibility = () => {
    setShowConfirmPassword(!showConfirmPassword);
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token) {
      router.push('/login');
      return;
    }
    
    if (formData.password !== formData.confirmPassword) {
      showToast("Passwords do not match", "error");
      return;
    }
    
    // Add additional password strength validation if needed here
    if (formData.password.length < 8) {
      showToast("Password must be at least 8 characters long", "error");
      return;
    }
    
    setIsLoading(true);

    try {
      await resetPasswordRequest({
        token: token,
        password: formData.password,
        password_confirmation: formData.confirmPassword
      });
      setIsSuccessModalOpen(true);
    } catch (error: unknown) {
      const err = error as { response?: { data?: { message?: string } } };
      showToast(err.response?.data?.message || "An error occurred while resetting your password", "error");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="w-full">
        <div className="mb-6 text-center">
            <h2 className="text-xl font-bold text-gray-800">
                Create New Password
            </h2>
            <p className="text-sm text-gray-500 mt-2">
                Your new password must be different from previously used passwords.
            </p>
        </div>

        <form onSubmit={handleSubmit} className={styles.form}>
        <div className={styles.inputGroup}>
            <label htmlFor="password" className={styles.label}>
            New Password
            </label>
            <div className={styles.inputWrapper}>
            <input
                id="password"
                name="password"
                type={showPassword ? "text" : "password"}
                placeholder="Enter new password"
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

        <div className={styles.inputGroup}>
            <label htmlFor="confirmPassword" className={styles.label}>
            Confirm Password
            </label>
            <div className={styles.inputWrapper}>
            <input
                id="confirmPassword"
                name="confirmPassword"
                type={showConfirmPassword ? "text" : "password"}
                placeholder="Confirm new password"
                className={styles.input}
                value={formData.confirmPassword}
                onChange={handleInputChange}
                required
            />
            <button
                type="button"
                className={styles.eyeIcon}
                onClick={toggleConfirmPasswordVisibility}
                aria-label={showConfirmPassword ? "Hide password" : "Show password"}
            >
                {showConfirmPassword ? (
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
                <span>Resetting...</span>
            </div>
            ) : (
            "Reset Password"
            )}
        </button>
        </form>

      <Modal isOpen={isSuccessModalOpen} onClose={() => {}} title="Password Reset Successful">
        <div className="flex flex-col items-center justify-center py-4">
          <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center text-green-500 mb-4">
            <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7"></path>
            </svg>
          </div>
          <p className="text-center text-gray-600 mb-6">Your password has been successfully reset. You can now login with your new password.</p>
          <button 
            type="button" 
            onClick={() => router.push('/login')} 
            className="w-full bg-primary text-white py-2.5 rounded-lg text-sm font-semibold hover:bg-primary/90 transition-colors shadow-sm"
          >
            Go to Login
          </button>
        </div>
      </Modal>
    </div>
  );
}

export default function CreateNewPasswordForm() {
  return (
    <Suspense fallback={<div className="flex justify-center p-8 w-full"><Loader className="animate-spin text-primary" size={32} /></div>}>
      <CreateNewPasswordFormContent />
    </Suspense>
  );
}
