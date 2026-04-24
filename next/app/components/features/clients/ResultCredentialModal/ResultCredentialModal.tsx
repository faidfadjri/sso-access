"use client";

import { Check, Key } from "react-feather";
import { ResultCredentialModalProps } from "./ResultCredentialModal.type";
import { CredentialField, Modal } from "@/app/components/common";

export default function ResultCredentialModal({ isOpen, onClose, ServiceClient }: ResultCredentialModalProps) {
  if (!ServiceClient) return null;

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={
        <div className="text-center pt-4">
          <div className="flex justify-center mb-4">
            <div className="w-16 h-16 bg-sky-50 rounded-full flex items-center justify-center relative">
              <div className="absolute inset-0 bg-sky-100 rounded-full opacity-50 scale-125 animate-pulse-slow" />
              <Key className="text-primary relative z-10" size={32} />
              <div className="absolute -bottom-1 -right-1 bg-white p-0.5 rounded-full z-20 shadow-sm">
                <div className="bg-emerald-500 rounded-full p-1">
                  <Check size={10} className="text-white" />
                </div>
              </div>
            </div>
          </div>
          <h2 className="font-h5 font-bold text-gray-900 mb-1">Credentials</h2>
          <p className="text-gray-500 text-xs px-6 mt-2">Please store these credentials securely.<br /></p>
        </div>
      }
    >
      <div className="space-y-4 mt-2">
        <CredentialField label="Service Name" value={ServiceClient.service_name} />
        <CredentialField label="Redirect URL" value={ServiceClient.redirect_url} />
        <CredentialField label="Client ID" value={ServiceClient.client_id} />
        <div className="pt-4">
          <button onClick={onClose} className="w-full bg-secondary text-white px-6 py-2.5 rounded-lg text-sm font-semibold hover:bg-secondary/90 transition-all shadow-lg shadow-secondary/20 cursor-pointer">
            I have saved these credentials
          </button>
        </div>
      </div>
    </Modal>
  );
}
