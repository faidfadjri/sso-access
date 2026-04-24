"use client";


import { useEffect, useState } from "react";
import { authSession } from "../api/services/oauth/oauth-session.service";
import { oauthGetMe, updateAccountRequest } from "../api/services/oauth/oauth.service";
import { JWTPayload, UpdateAccountRequest } from "../api/services/oauth/oauth.type";
import { useToast } from "../context/ToastContext";
import { PersonalInfoForm } from "../components/forms";
import { ProfileCard } from "../components/common";

export default function AccountPage() {
  const [user, setUser] = useState<JWTPayload | null>(null);
  const [loading, setLoading] = useState(false);
  const [avatarFile, setAvatarFile] = useState<File | null>(null);
  const [isEditing, setIsEditing] = useState(false);
  const toast = useToast();

  useEffect(() => {
    const sessionUser = authSession.get();
    if (sessionUser) {
      setUser(sessionUser);
    }

    oauthGetMe().then(response => {
       if(response.data) {
         setUser(prev => {
            const baseUser = prev || { session_id: "" }; 
            return {
                ...baseUser,
                full_name: response.data.full_name,
                email: response.data.email,
                username: response.data.username,
                photo: response.data.photo,
                phone: response.data.phone,
                role_name: response.data.role_name,
            } as JWTPayload;
         });
       }
    }).catch(err => console.error("Failed to refresh user data", err));
  }, []);


  const handleUpdateAccount = async (data: UpdateAccountRequest) => {
    try {
      setLoading(true);
      
      const payload = { ...data };
      if (avatarFile) {
        payload.photo = avatarFile;
      }

      const response = await updateAccountRequest(payload);

      // Optimistic update using payload to circumvent stale JWT cookies on /me endpoint
      const optimisticUpdate: Partial<JWTPayload> = {
         full_name: typeof payload.full_name === 'string' ? payload.full_name : undefined,
         username: typeof payload.username === 'string' ? payload.username : undefined,
         email: typeof payload.email === 'string' ? payload.email : undefined,
         phone: typeof payload.phone === 'string' ? payload.phone : undefined,
      };

      // Clean undefined values
      Object.keys(optimisticUpdate).forEach(key => optimisticUpdate[key as keyof Partial<JWTPayload>] === undefined && delete optimisticUpdate[key as keyof Partial<JWTPayload>]);

      let backendData = {};
      if (response && response.data) {
         backendData = response.data;
      }

      setUser(prev => {
          const updated = { ...prev, ...optimisticUpdate, ...backendData } as JWTPayload;
          authSession.updateUser(updated);
          return updated;
      });

      toast.showToast("Profile updated successfully", "success");
      setAvatarFile(null);
      setIsEditing(false); 

    } catch (error: any) {
      console.error("Failed to update account", error);

      const message = error.response?.data?.message || "Failed to update profile";
      toast.showToast(message, "error");
      throw error; 
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="grid grid-cols-1 xl:grid-cols-3 gap-4 items-start max-w-7xl mx-auto">
        <div className="bg-transparent col-span-1 xl:col-span-1">
          <ProfileCard user={user} onAvatarChange={setAvatarFile} isEditing={isEditing} />
        </div>
        
        <div className="bg-transparent col-span-1 xl:col-span-2">
          <PersonalInfoForm 
            user={user} 
            onSave={handleUpdateAccount} 
            isLoading={loading}
            isEditing={isEditing}
            onToggleEdit={() => setIsEditing(!isEditing)}
          />
        </div>
    </div>
  );
}
