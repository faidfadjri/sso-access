"use client";

import { useEffect, useState } from "react";
import { Edit2, Save, Eye, EyeOff } from "react-feather";
import { AccountInfo } from "./PersonalInfoForm.type";
import { LoginResponse, UpdateAccountRequest, JWTPayload } from "@/app/api/services/oauth/oauth.type";

type PersonalInfoFormProps = {
  user: JWTPayload | null;
  onSave: (data: UpdateAccountRequest) => Promise<void>;
  isLoading?: boolean;
  isEditing: boolean;
  onToggleEdit: () => void;
}

export default function PersonalInfoForm({ 
  user, 
  onSave, 
  isLoading = false, 
  isEditing, 
  onToggleEdit 
}: PersonalInfoFormProps) {
  const [formData, setFormData] = useState<AccountInfo>({
    fullName: "",
    username: "",
    phone: "",
    photo: "",
    email: "",
    password: "",
    passwordConfirm: ""
  });

  useEffect(() => {
    if (user) {
      setFormData(prev => ({
        ...prev,
        fullName: user.full_name || "",
        username: user.username || "",
        phone: user.phone || "",
        photo: user.photo || "", 
        email: user.email || "",
      }));
    }
  }, [user]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const handleSave = async () => {
    if (!isEditing) {
      onToggleEdit();
      return;
    }

    const payload: UpdateAccountRequest = {};
    if (formData.fullName) payload.full_name = formData.fullName;
    if (formData.username) payload.username = formData.username;
    if (formData.phone) payload.phone = formData.phone;
    if (formData.email) payload.email = formData.email;
    if (formData.password) payload.password = formData.password;
    if (formData.passwordConfirm) payload.password_confirmation = formData.passwordConfirm;

    await onSave(payload);
  };

  return (
    <div className="bg-white rounded-[5px] p-8 shadow-xs">
      <div className="flex flex-row justify-between items-center mb-8 border-b border-gray-200 pb-3">
       <div className="flex flex-col">
         <h3 className="text-lg md:text-xl font-bold text-gray-800">Personal Info</h3>
          <p className="text-sm text-gray-500">Update your personal information down here</p>
       </div>
        <button 
          type="button"
          onClick={handleSave}
          disabled={isLoading}
          className={`
            flex items-center gap-2 px-4 py-2 rounded-lg transition-colors font-medium text-sm cursor-pointer
            ${isEditing 
              ? "bg-primary text-white hover:bg-primary-dark shadow-md" 
              : "text-primary border border-primary/20 hover:bg-primary/5"} 
            ${isLoading ? "opacity-70 cursor-not-allowed" : ""}
          `}
        >
            {isLoading ? (
                <span className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
            ) : (
                isEditing ? <Save size={16} /> : <Edit2 size={16} />
            )}
            {isEditing ? "Save" : "Edit"}
        </button>
      </div>

      <div className="grid grid-cols-1 gap-6">
        <FormInput 
          label="Full Name" 
          name="fullName" 
          value={formData.fullName} 
          onChange={handleChange} 
          disabled={!isEditing} 
          required 
        />
        
        <FormInput 
          label="Email" 
          name="email" 
          value={formData.email} 
          onChange={handleChange} 
          disabled={!isEditing} 
          required 
        />

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <FormInput 
            label="Username" 
            name="username" 
            value={formData.username} 
            onChange={handleChange} 
            disabled={!isEditing} 
            required 
          />
          <FormInput 
            label="Phone Number" 
            name="phone" 
            value={formData.phone} 
            onChange={handleChange} 
            disabled={!isEditing} 
            required 
          />
        </div>

        {isEditing && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <PasswordInput
              label="Password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              disabled={!isEditing}
            />
            <PasswordInput
              label="Password Confirmation"
              name="passwordConfirm"
              value={formData.passwordConfirm}
              onChange={handleChange}
              disabled={!isEditing}
            />
          </div>
        )}
      </div>
    </div>
  );
}

type FormInputProps = {
  label: string;
  name: string;
  value: string | undefined;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  disabled: boolean;
  required?: boolean;
}

const FormInput = ({ label, name, value, onChange, disabled, required }: FormInputProps) => (
  <div className="space-y-2">
    <label className="block text-sm font-semibold text-gray-700">
      {label} {required && <span className="text-red-800">*</span>}
    </label>
    <input
      type="text"
      name={name}
      value={value || ""}
      onChange={onChange}
      disabled={disabled}
      className={`w-full px-4 py-3 rounded-lg text-sm font-medium border transition-colors focus:outline-none focus:ring-1 focus:ring-primary/0 ${
        disabled 
          ? "bg-slate-50 border-gray-100 text-gray-600 cursor-default" 
          : "bg-white border-blue-500/50 text-gray-800"
      }`}
      required={required}
    />
  </div>
);

const PasswordInput = ({ label, name, value, onChange, disabled }: Omit<FormInputProps, 'required'>) => {
  const [show, setShow] = useState(false);
  
  return (
    <div className="space-y-2">
      <label className="block text-sm font-semibold text-gray-700">{label}</label>
      <div className="relative">
        <input
          type={show ? "text" : "password"}
          name={name}
          value={value || ""}
          onChange={onChange}
          disabled={disabled}
          className={`w-full px-4 py-3 rounded-lg text-sm font-medium border transition-colors focus:outline-none focus:ring-1 focus:ring-primary/0 ${
            disabled 
              ? "bg-slate-50 border-gray-100 text-gray-400 cursor-default" 
              : "bg-white border-blue-500/50 text-gray-800"
          }`}
        />
        <button
          type="button"
          onClick={() => setShow(!show)}
          className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
        >
          {show ? <EyeOff size={18} /> : <Eye size={18} />}
        </button>
      </div>
    </div>
  );
};