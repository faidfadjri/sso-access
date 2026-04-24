"use client";

import { useState, useRef, useEffect } from "react";
import { Upload, Check } from "react-feather";
import Modal from "@/app/components/common/Modal/Modal";
import { createUser, updateUser } from "@/app/api/services/users/users.service";
import { useToast } from "@/app/context/ToastContext";
import { User } from "@/app/api/services/users/users.service.type";

interface UserModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
  user?: User | null;
}

interface FormData {
  full_name: string;
  email: string;
  username: string;
  phone: string;
  password: string;
  passwordConfirmation: string;
  photo: File | null;
}

const INITIAL_FORM_DATA: FormData = {
  full_name: "",
  email: "",
  username: "",
  phone: "",
  password: "",
  passwordConfirmation: "",
  photo: null,
};

export default function CreateUserModal({
  isOpen,
  onClose,
  onSuccess,
  user
}: UserModalProps) {
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState<FormData>(INITIAL_FORM_DATA);
  const [previewUrl, setPreviewUrl] = useState<string>("");
  const fileInputRef = useRef<HTMLInputElement>(null);

  const toast = useToast();
  const isEditMode = !!user;

  useEffect(() => {
    if (isOpen) {
        if (user) {
            setFormData({
                full_name: user.full_name,
                email: user.email,
                username: user.username,
                phone: user.phone,
                password: "",
                passwordConfirmation: "",
                photo: null,
            });
            setPreviewUrl(typeof user.photo === 'string' ? process.env.NEXT_PUBLIC_BASE_API_URL + user.photo : "");
        } else {
            resetForm();
        }
    }
  }, [isOpen, user]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
        const file = e.target.files[0];
      setFormData((prev) => ({
        ...prev,
        photo: file,
      }));
      setPreviewUrl(URL.createObjectURL(file));
    }
  };

  const validateForm = (): boolean => {
    const { full_name, email, username, phone, password, passwordConfirmation } = formData;

    if (!full_name.trim()) {
      toast.showToast("Full name is required", "error");
      return false;
    }
    if (!email.trim() || !/^\S+@\S+\.\S+$/.test(email)) {
      toast.showToast("Valid email is required", "error");
      return false;
    }
    if (!username.trim()) {
      toast.showToast("Username is required", "error");
      return false;
    }
    if (!phone.trim()) {
      toast.showToast("Phone is required", "error");
      return false;
    }
    
    if (!isEditMode || password) {
        if (!password) {
            toast.showToast("Password is required", "error");
            return false;
        }
        if (password !== passwordConfirmation) {
            toast.showToast("Passwords do not match", "error");
            return false;
        }
    }
    
    return true;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) return;

    try {
      setLoading(true);
      
      const payload: any = {
        full_name: formData.full_name,
        email: formData.email,
        username: formData.username,
        phone: formData.phone,
        photo: formData.photo || "", // Handle optional photo
        // Include password only if provided
        ...(formData.password && formData.password.trim() !== "" && { 
            password: formData.password,
            password_confirmation: formData.passwordConfirmation
        }),
      };

      if (isEditMode && user) {
        const updateData: any = {
            ...user, // Spread existing user to keep ID, created_at etc
            ...payload
        };
        
        // Ensure we don't send password if it wasn't intended to be updated
        if (!payload.password && updateData.password) {
            delete updateData.password;
        }

        await updateUser(user.user_id.toString(), updateData);
        toast.showToast("User updated successfully", "success");
      } else {
         await createUser({
            ...payload,
            user_id: 0,
            created_at: new Date(),
            updated_at: new Date(),
            password: formData.password, // Required for create
            password_confirmation: formData.passwordConfirmation
         });
         toast.showToast("User created successfully", "success");
      }
      
      resetForm();
      onSuccess();
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : `Failed to ${isEditMode ? 'update' : 'create'} user`;
      toast.showToast(errorMessage, "error");
    } finally {
      setLoading(false);
    }
  };

  const resetForm = () => {
    setFormData(INITIAL_FORM_DATA);
    setPreviewUrl("");
    if (fileInputRef.current) {
        fileInputRef.current.value = "";
    }
  };

  const handleClose = () => {
    resetForm();
    onClose();
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={handleClose}
      title={
        <div>
          <h2 className="text-xl font-bold text-gray-900">
            {isEditMode ? "Edit User" : "Create a New User"}
          </h2>
          <p className="text-sm text-gray-500 mt-1 font-normal">
            Please fill the required form below
          </p>
        </div>
      }
    >
      <form className="space-y-5" onSubmit={handleSubmit}>
        {/* Photo Upload */}
        <div className="flex justify-center">
            <div className="relative group cursor-pointer" onClick={() => fileInputRef.current?.click()}>
                <div className="w-24 h-24 rounded-full bg-gray-100 border-2 border-dashed border-gray-300 flex items-center justify-center overflow-hidden hover:border-primary transition-colors">
                    {previewUrl ? (
                        <img 
                            src={previewUrl} 
                            alt="Preview" 
                            className="w-full h-full object-cover"
                        />
                    ) : (
                        <div className="text-center text-gray-400">
                            <Upload size={24} className="mx-auto mb-1" />
                            <span className="text-[10px]">Upload</span>
                        </div>
                    )}
                </div>
                <input 
                    type="file" 
                    ref={fileInputRef} 
                    onChange={handleFileChange} 
                    className="hidden" 
                    accept="image/*"
                />
            </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
             {/* Full Name */}
            <div className="space-y-1.5">
            <label className="block text-xs font-bold text-gray-700">
                Full Name <span className="text-red-500">*</span>
            </label>
            <input
                type="text"
                name="full_name"
                value={formData.full_name}
                onChange={handleInputChange}
                placeholder="e.g: Ahmad Zakaria"
                disabled={loading}
                className="w-full px-3 py-2.5 rounded-lg border border-gray-200 text-sm
                focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary
                transition-all placeholder:text-gray-300 disabled:bg-gray-50 disabled:cursor-not-allowed"
            />
            </div>

            {/* Email */}
            <div className="space-y-1.5">
            <label className="block text-xs font-bold text-gray-700">
                Email <span className="text-red-500">*</span>
            </label>
            <input
                type="email"
                name="email"
                value={formData.email}
                onChange={handleInputChange}
                placeholder="e.g: user@example.com"
                disabled={loading}
                className="w-full px-3 py-2.5 rounded-lg border border-gray-200 text-sm
                focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary
                transition-all placeholder:text-gray-300 disabled:bg-gray-50 disabled:cursor-not-allowed"
            />
            </div>

            {/* Username */}
            <div className="space-y-1.5">
            <label className="block text-xs font-bold text-gray-700">
                Username <span className="text-red-500">*</span>
            </label>
            <input
                type="text"
                name="username"
                value={formData.username}
                onChange={handleInputChange}
                placeholder="e.g: ahmadzakaria"
                disabled={loading}
                className="w-full px-3 py-2.5 rounded-lg border border-gray-200 text-sm
                focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary
                transition-all placeholder:text-gray-300 disabled:bg-gray-50 disabled:cursor-not-allowed"
            />
            </div>

            {/* Phone */}
            <div className="space-y-1.5">
            <label className="block text-xs font-bold text-gray-700">
                Phone <span className="text-red-500">*</span>
            </label>
            <input
                type="text"
                name="phone"
                value={formData.phone}
                onChange={handleInputChange}
                placeholder="e.g: 0858xxxxxx"
                disabled={loading}
                className="w-full px-3 py-2.5 rounded-lg border border-gray-200 text-sm
                focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary
                transition-all placeholder:text-gray-300 disabled:bg-gray-50 disabled:cursor-not-allowed"
            />
            </div>

            {/* Password */}
            <div className="space-y-1.5">
            <label className="block text-xs font-bold text-gray-700">
                Password {!isEditMode && <span className="text-red-500">*</span>}
            </label>
            <input
                type="password"
                name="password"
                value={formData.password}
                onChange={handleInputChange}
                placeholder={isEditMode ? "Leave blank to keep current" : "******"}
                disabled={loading}
                autoComplete="new-password"
                className="w-full px-3 py-2.5 rounded-lg border border-gray-200 text-sm
                focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary
                transition-all placeholder:text-gray-300 disabled:bg-gray-50 disabled:cursor-not-allowed"
            />
            </div>

            {/* Password Confirmation */}
            <div className="space-y-1.5">
            <label className="block text-xs font-bold text-gray-700">
                Confirm Password {!isEditMode && <span className="text-red-500">*</span>}
            </label>
            <input
                type="password"
                name="passwordConfirmation"
                value={formData.passwordConfirmation}
                onChange={handleInputChange}
                placeholder={isEditMode ? "Leave blank to keep current" : "******"}
                disabled={loading}
                autoComplete="new-password"
                className="w-full px-3 py-2.5 rounded-lg border border-gray-200 text-sm
                focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary
                transition-all placeholder:text-gray-300 disabled:bg-gray-50 disabled:cursor-not-allowed"
            />
            </div>
        </div>

        {/* Actions */}
        <div className="pt-4 flex justify-end gap-3 border-t border-gray-100">
          <button
            type="button"
            onClick={handleClose}
            disabled={loading}
            className="px-6 py-2.5 rounded-lg text-sm font-bold text-gray-700 bg-slate-100 hover:bg-slate-200 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={loading}
            className="bg-secondary text-white px-8 py-2.5 rounded-lg text-sm font-bold
            hover:bg-secondary/90 transition-colors shadow-lg shadow-secondary/20 cursor-pointer
            disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
          >
            {loading && (
              <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
            )}
            {loading ? "Saving..." : (isEditMode ? "Update User" : "Save User")}
          </button>
        </div>
      </form>
    </Modal>
  );
}
