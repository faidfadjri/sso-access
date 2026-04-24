"use client"

import { Briefcase, Phone, Edit2 } from "react-feather";
import { JWTPayload } from "@/app/api/services/oauth/oauth.type";
import { useState, useRef, useEffect } from "react";

export type ProfileCardProps = {
  user: JWTPayload | null;
  onAvatarChange?: (file: File) => void;
  isEditing?: boolean;
}

export default function ProfileCard({ user, onAvatarChange, isEditing = false }: ProfileCardProps) {
  const [preview, setPreview] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const getInitials = (name: string) => {
    if (!name) return "U";
    return name
      .split(" ")
      .map((n) => n[0])
      .slice(0, 2)
      .join("")
      .toUpperCase();
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      const objectUrl = URL.createObjectURL(file);
      setPreview(objectUrl);
      if (onAvatarChange) {
        onAvatarChange(file);
      }
    }
  };

  const handleEditClick = () => {
    if (isEditing) {
        fileInputRef.current?.click();
    }
  };
  
  // Cleanup preview URL on unmount
  useEffect(() => {
    return () => {
        if (preview) URL.revokeObjectURL(preview);
    }
  }, [preview]);

  return (
    <div className="bg-white rounded-[10px] p-6 shadow-xs h-fit relative">
      <div className="w-full h-1/4 bg-gradient-to-r from-pink-200 to-blue-200 absolute top-0 left-0 rounded-t-[10px]"></div>    
      <div className="flex flex-col items-center">
        <div className="relative mb-4">
           <div className="relative z-10 p-1 bg-white rounded-full">
            <div className={`w-[120px] h-[120px] rounded-full overflow-hidden border-4 border-white shadow-lg relative bg-gray-200 group ${isEditing ? 'cursor-pointer' : ''}`} onClick={handleEditClick}>
            
               {preview ? (
                 <img src={preview} alt="Profile Preview" className="w-full h-full object-cover" />
               ) : user?.photo ? (
                 <img src={`${process.env.NEXT_PUBLIC_BASE_API_URL}${user?.photo}`} alt="Profile" className="w-full h-full object-cover" />
               ) : (
                 <div className="w-full h-full flex items-center justify-center bg-primary text-white text-4xl font-bold">
                    {getInitials(user?.full_name || "")}
                 </div>
               )}
               
               {isEditing && (
                <div className="absolute inset-0 bg-black/30 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                        <Edit2 className="text-white" size={24} />
                </div>
               )}
            </div>
            <input 
                type="file" 
                ref={fileInputRef} 
                className="hidden" 
                accept="image/*"
                onChange={handleFileChange}
                disabled={!isEditing}
            />
           </div>
        </div>

        <h3 className="text-xl font-bold text-gray-800 mb-1">{user?.full_name || ""}</h3>
        <p className="text-gray-500 font-medium text-sm mb-6">{user?.role_name}</p>

        <div className="w-full space-y-4 border-t border-gray-100 pt-6">
            <div className="flex items-start gap-4">
                <div className="p-2 bg-yellow-100 text-yellow-600 rounded-full shrink-0">
                    <Briefcase size={20} />
                </div>
                <div>
                    <p className="text-xs md:text-sm text-gray-400 font-semibold uppercase tracking-wide">Company</p>
                    <p className="text-sm md:text-md font-semibold text-gray-700">Akastra Toyota</p>
                </div>
            </div>

            {user?.phone && (
                <div className="flex items-start gap-4">
                <div className="p-2 bg-teal-100 text-teal-600 rounded-full shrink-0">
                    <Phone size={20} />
                </div>
                <div>
                    <p className="text-xs md:text-sm text-gray-400 font-semibold uppercase tracking-wide">Phone Number</p>
                    <p className="text-sm md:text-md font-semibold text-gray-700">{user?.phone || ""}</p>
                </div>
            </div>
            )}
        </div>
      </div>
    </div>
  );
}
