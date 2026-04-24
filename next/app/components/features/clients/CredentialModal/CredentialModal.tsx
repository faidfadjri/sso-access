"use client";

import { Folder, X } from "react-feather";
import { useState, useRef, useEffect } from "react";
import Image from "next/image";
import { createServiceClient, updateServiceClient } from "@/app/api/services/clients/client.service";
import { Client } from "@/app/api/services/clients/client.type";
import { useToast } from "@/app/context/ToastContext";
import { CredentialField, Modal } from "@/app/components/common";
import { CredentialModalProps, formValues } from "./CredentialModal.type";
import { API_BASE_URL } from "@/app/api/config";

export default function CredentialModal({ isOpen, onClose, onSuccess, data, mode }: CredentialModalProps) {
  // const isDebug = process.env.NEXT_PUBLIC_DEBUG === "true";
  const defaultFormData: formValues = { 
    name: data?.service_name || "", 
    redirect_url: data?.redirect_url || "", 
    description: data?.description || "", 
    service_id: data?.service_id || null,
    is_active: data?.is_active ?? true
  };
  const [formData, setFormData] = useState<formValues>(defaultFormData);
  const [file, setFile] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [createdClient, setCreatedClient] = useState<Client | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const toast = useToast();

  useEffect(() => {
    setFormData({
      service_id: data?.service_id ?? null,
      name: data?.service_name ?? "",
      redirect_url: data?.redirect_url ?? "",
      description: data?.description ?? "",
      is_active: data?.is_active ?? true
    })

    if (mode === "update" && data?.logo) {
      const formattedLogoPath = data.logo.replace(/\\/g, '/');
      const cleanBaseUrl = (API_BASE_URL || '').replace(/\/$/, '');
      const cleanLogoPath = formattedLogoPath.replace(/^\//, '');
      setPreviewUrl(`${cleanBaseUrl}/${cleanLogoPath}`);
    } else {
      setPreviewUrl(null);
    }

  setFile(null)


  }, [data, isOpen, mode])

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0];
    if (selectedFile) {
      if (selectedFile.size > 2 * 1024 * 1024) {
        toast.showToast("Image size should be less than 2MB", "error");
        return;
      }
      setFile(selectedFile);
      setPreviewUrl(URL.createObjectURL(selectedFile));
    }
  };

  const handleRemoveImage = (e: React.MouseEvent) => {
    e.stopPropagation();
    setFile(null);
    if (previewUrl) { URL.revokeObjectURL(previewUrl); setPreviewUrl(null); }
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

  const handleSubmit = async (e: React.SubmitEvent) => {
    e.preventDefault();

    if (!formData.name || !formData.redirect_url) {
        toast.showToast("Please fill all required fields", "error")
        return
    }

    if (mode === "create" && !file) {
      toast.showToast("Logo is required", "error")
      return
    }

    try {
        setLoading(true)

        const payload = new FormData()
        payload.append("name", formData.name)
        payload.append("redirect_url", formData.redirect_url)
        payload.append("description", formData.description)

        if (file) {
          payload.append("logo", file)
        }

        if (mode === "create") {
          const response = await createServiceClient(payload)

          if (response?.data) setCreatedClient(response.data);
          else { onSuccess(); onClose(); }
        } else {
          payload.append("is_active", String(formData.is_active))
          await updateServiceClient(formData.service_id!, payload)
          onSuccess()
          onClose()
        }

        toast.showToast(
          mode === "create"
            ? "Service client created successfully"
            : "Service client updated successfully",
          "success"
        )

        

      } catch (err) {
        console.log(err)
        toast.showToast("Failed to process request", "error")
      } finally {
        setLoading(false)
      }
  };

  const handleCloseSuccess = () => { setCreatedClient(null); onSuccess(); onClose(); };

  if (createdClient) {
    return (
      <Modal isOpen={isOpen} onClose={handleCloseSuccess} title="Credentials Created">
        <div className="space-y-4">
          <div className="bg-rose-50 text-rose-800 p-3 rounded-lg text-sm border border-blue-100">
            Here are your Client ID and Client Secret. <strong>Please copy them now</strong> as the secret will not be shown again.
          </div>
          <CredentialField label="Client ID" value={createdClient.client_id} />
          <CredentialField label="Client Secret" value={createdClient.client_secret} isSecret />
          <div className="pt-4 flex justify-end">
            <button type="button" onClick={handleCloseSuccess} className="bg-secondary text-white px-6 py-2.5 rounded-lg text-sm font-semibold hover:bg-primary/90 transition-colors shadow-sm">
              Got it!
            </button>
          </div>
        </div>
      </Modal>
    );
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={<>{mode == "create"? "Create a New Credentials" : "Update Credentials"}<p className="text-gray-500 text-xs mt-1 font-normal">Please fill the required form below</p></>}>
      <form className="space-y-4" onSubmit={handleSubmit}>
        <div className="space-y-1.5">
          <label className="block text-xs font-bold text-gray-700">Service Icon <span className="text-red-500">*</span></label>
          <input type="file" ref={fileInputRef} className="hidden" accept="image/*" onChange={handleFileChange} />
          <div onClick={() => fileInputRef.current?.click()} className={`border border-gray-200 rounded-lg p-3 flex items-center gap-3 cursor-pointer transition-colors bg-white hover:border-primary/50 ${previewUrl ? "border-primary/50" : ""}`}>
            {previewUrl ? (
              <div className="relative w-10 h-10 shrink-0">
                <Image src={previewUrl} alt="Preview" fill className="object-cover rounded-md" unoptimized/>
                <button type="button" onClick={handleRemoveImage} className="absolute -top-2 -right-2 bg-red-500 text-white rounded-full p-0.5 hover:bg-red-600 transition-colors"><X size={12} /></button>
              </div>
            ) : (
              <div className="w-10 h-10 rounded-md bg-gray-50 flex items-center justify-center shrink-0 text-gray-400"><Folder size={20} /></div>
            )}
            <div className="flex flex-col">
              <span className={`text-sm ${previewUrl ? "text-gray-900 font-medium" : "text-gray-400"}`}>{file ? file.name : "Browse Service Icon Image"}</span>
              <span className="text-xs text-gray-400">Max size: 2MB (JPG, PNG)</span>
            </div>
          </div>
        </div>

        {[{ label: "Service Name", name: "name" as const, placeholder: "e.g: My App", type: "text" }, { label: "Redirect URL", name: "redirect_url" as const, placeholder: "e.g: https://myapp.com/redirect-url", type: "text" }].map((field) => (
          <div key={field.name} className="space-y-1.5">
            <label className="block text-xs font-bold text-gray-700">{field.label} <span className="text-red-500">*</span></label>
            <input type={field.type} name={field.name} value={formData[field.name]} onChange={handleInputChange} placeholder={field.placeholder} className="w-full px-3 py-2.5 rounded-lg border border-gray-200 text-sm focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary transition-all placeholder:text-gray-300" />
          </div>
        ))}

        <div className="space-y-1.5">
          <label className="block text-xs font-bold text-gray-700">Description</label>
          <textarea name="description" value={formData.description} onChange={handleInputChange} placeholder="Brief description of the service" rows={2} className="w-full px-3 py-2.5 rounded-lg border border-gray-200 text-sm focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary transition-all placeholder:text-gray-300 resize-none" />
        </div>

        {mode === "update" && (
          <div className="space-y-1.5">
            <label className="block text-xs font-bold text-gray-700">Status</label>
            <div className="flex bg-gray-100/80 p-1 rounded-md w-full sm:w-[240px]">
              <label className="flex-1 text-center cursor-pointer">
                <input
                  type="radio"
                  name="is_active"
                  className="hidden"
                  checked={formData.is_active === true}
                  onChange={() => setFormData((prev) => ({ ...prev, is_active: true }))}
                />
                <div className={`py-2 text-sm font-semibold rounded transition-all duration-200 ${
                  formData.is_active 
                    ? 'bg-white text-primary shadow-sm ring-1 ring-gray-200/50' 
                    : 'text-gray-500 hover:text-gray-700'
                }`}>
                  Active
                </div>
              </label>
              <label className="flex-1 text-center cursor-pointer">
                <input
                  type="radio"
                  name="is_active"
                  className="hidden"
                  checked={formData.is_active === false}
                  onChange={() => setFormData((prev) => ({ ...prev, is_active: false }))}
                />
                <div className={`py-2 text-sm font-semibold rounded transition-all duration-200 ${
                  formData.is_active === false
                    ? 'bg-white text-red-500 shadow-sm ring-1 ring-gray-200/50' 
                    : 'text-gray-500 hover:text-gray-700'
                }`}>
                  Inactive
                </div>
              </label>
            </div>
          </div>
        )}

        <div className="pt-4 flex justify-end">
          <button type="submit" disabled={loading} className="bg-secondary text-white px-6 py-2.5 rounded-lg text-sm font-semibold hover:bg-secondary/90 transition-colors shadow-lg shadow-secondary/20 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed">
            {loading ? "Generating..." : (mode === "create" ? "Generate Credentials" : "Update Credentials")}
          </button>
        </div>
      </form>
    </Modal>
  );
}
