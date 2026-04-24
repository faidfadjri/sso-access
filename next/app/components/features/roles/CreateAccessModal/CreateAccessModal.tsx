"use client";

import { useState, useEffect } from "react";
import { Modal } from "@/app/components/common";
import { UserAutosuggest } from "@/app/components/features/users/UserAutosuggest/UserAutosuggest";
import { User } from "@/app/api/services/users/users.service.type";
import { getServiceClients } from "@/app/api/services/clients";
import { Client } from "@/app/api/services/clients/client.type";
import { createUserAccess } from "@/app/api/services/user-access/user-access.service";
import { useToast } from "@/app/context/ToastContext";

interface CreateAccessModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

export const CreateAccessModal = ({ isOpen, onClose, onSuccess }: CreateAccessModalProps) => {
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [services, setServices] = useState<Client[]>([]);
  const [selectedServices, setSelectedServices] = useState<number[]>([]);
  const [loading, setLoading] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const toast = useToast();

  useEffect(() => {
    if (isOpen) {
      fetchServices();
      setSelectedUser(null);
      setSelectedServices([]);
    }
  }, [isOpen]);

  const fetchServices = async () => {
    try {
      setLoading(true);
      const response = await getServiceClients({ page: 1, show: 100, sort: "desc" });
      setServices(response.data.rows);
    } catch (error) {
      toast.showToast("Failed to load services", "error");
    } finally {
      setLoading(false);
    }
  };

  const handleServiceToggle = (serviceId: number) => {
    setSelectedServices((prev) => prev.includes(serviceId) ? prev.filter((id) => id !== serviceId) : [...prev, serviceId]);
  };

  const handleSelectAll = () => {
    setSelectedServices(selectedServices.length === services.length ? [] : services.map((s) => s.service_id));
  };

  const handleSubmit = async () => {
    if (!selectedUser) { toast.showToast("Please select a user", "error"); return; }
    if (selectedServices.length === 0) { toast.showToast("Please select at least one service", "error"); return; }
    try {
      setIsSubmitting(true);
      await createUserAccess({ user_id: selectedUser.user_id, service_ids: selectedServices, status: "active" });
      toast.showToast("Access created successfully", "success");
      onSuccess();
    } catch (error) {
      toast.showToast("Failed to create access", "error");
    } finally {
      setIsSubmitting(false);
    }
  };

  if (!isOpen) return null;

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={
        <div className="mb-0">
          <h2 className="text-xl font-bold text-gray-900">Create Access</h2>
          <p className="text-sm text-gray-500 mt-1 font-normal">Please fill the required form below</p>
        </div>
      }
      maxWidth="max-w-md"
    >
      <div className="mt-6">
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-semibold text-gray-900 mb-2">Choose User <span className="text-red-500">*</span></label>
              <UserAutosuggest onSelect={setSelectedUser} selectedUser={selectedUser} />
            </div>
            <div>
              <div className="flex items-center justify-between mb-2">
                <label className="block text-sm font-semibold text-gray-900">Choose Service <span className="text-red-500">*</span></label>
                <button onClick={handleSelectAll} className="text-sm text-blue-600 hover:text-blue-800 underline font-medium">
                  Select All ({selectedServices.length}/{services.length})
                </button>
              </div>
              <div className="grid grid-cols-2 gap-3 max-h-48 overflow-y-auto pr-1">
                {loading ? (
                  <div className="col-span-2 text-center py-4 text-gray-500 text-sm">Loading services...</div>
                ) : (
                  services.map((service) => (
                    <label key={service.service_id} className="flex items-center gap-2 p-2 rounded-lg hover:bg-gray-50 cursor-pointer border border-transparent hover:border-gray-200 transition-colors">
                      <input type="checkbox" className="w-4 h-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500" checked={selectedServices.includes(service.service_id)} onChange={() => handleServiceToggle(service.service_id)} />
                      <span className="text-sm text-gray-700 truncate">{service.service_name}</span>
                    </label>
                  ))
                )}
              </div>
            </div>
          </div>
          <div className="flex gap-3 justify-end mt-8">
            <button onClick={onClose} disabled={isSubmitting} className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors">Cancel</button>
            <button onClick={handleSubmit} disabled={isSubmitting} className="px-4 py-2 text-sm font-medium text-white bg-primary rounded-lg hover:bg-primary/90 transition-colors disabled:opacity-70 disabled:cursor-not-allowed">
              {isSubmitting ? "Saving..." : "Save"}
            </button>
          </div>
      </div>
    </Modal>
  );
};
