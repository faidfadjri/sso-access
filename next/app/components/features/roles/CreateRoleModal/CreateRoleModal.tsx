"use client";

import { Modal } from "@/app/components/common";
import { useState, useEffect } from "react";
import { useToast } from "@/app/context/ToastContext";
import { getServiceClients } from "@/app/api/services/clients";
import { Client } from "@/app/api/services/clients/client.type";
import { createRole, updateRole } from "@/app/api/services/roles/roles.service";
import { Role } from "@/app/api/services/roles/roles.service.type";

interface CreateRoleModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
  role?: Role | null;
}

export const CreateRoleModal = ({ isOpen, onClose, onSuccess, role }: CreateRoleModalProps) => {
  const [roleName, setRoleName] = useState("");
  const [selectedServiceId, setSelectedServiceId] = useState("");
  const [services, setServices] = useState<Client[]>([]);
  const [loadingServices, setLoadingServices] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const toast = useToast();
  const isEdit = !!role;

  useEffect(() => {
    if (isOpen) {
      fetchServices();
      if (role) { setRoleName(role.role_name); setSelectedServiceId(String(role.service_id)); }
      else { setRoleName(""); setSelectedServiceId(""); }
    }
  }, [isOpen, role]);

  const fetchServices = async () => {
    try {
      setLoadingServices(true);
      const response = await getServiceClients({ page: 1, show: 100, sort: "desc" });
      setServices(response.data.rows);
    } catch (error) {
      toast.showToast("Failed to load services", "error");
    } finally {
      setLoadingServices(false);
    }
  };

  const handleSubmit = async () => {
    if (!roleName.trim()) { toast.showToast("Role Name is required", "error"); return; }
    if (!selectedServiceId) { toast.showToast("Service is required", "error"); return; }
    try {
      setIsSubmitting(true);
      const data = { role_name: roleName, service_id: Number(selectedServiceId) };
      if (isEdit && role) {
        await updateRole(Number(role.service_role_id), data);
        toast.showToast("Role updated successfully", "success");
      } else {
        await createRole(data);
        toast.showToast("Role created successfully", "success");
      }
      onSuccess();
    } catch (error) {
      toast.showToast(`Failed to ${isEdit ? "update" : "create"} role`, "error");
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
          <h2 className="text-xl font-bold text-gray-900">{isEdit ? "Edit Role" : "Create Role"}</h2>
          <p className="text-sm text-gray-500 mt-1 font-normal">{isEdit ? "Update the role details below" : "Please fill the required form below"}</p>
        </div>
      }
      maxWidth="max-w-md"
    >
      <div className="mt-6">
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-semibold text-gray-900 mb-2">Choose Service <span className="text-red-500">*</span></label>
              <select className="w-full px-4 py-2 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary bg-white" value={selectedServiceId} onChange={(e) => setSelectedServiceId(e.target.value)} disabled={loadingServices}>
                <option value="">Select a service</option>
                {services.map((service) => (<option key={service.service_id} value={service.service_id}>{service.service_name}</option>))}
              </select>
              {loadingServices && <p className="text-xs text-gray-500 mt-1">Loading services...</p>}
            </div>
            <div>
              <label className="block text-sm font-semibold text-gray-900 mb-2">Role Name <span className="text-red-500">*</span></label>
              <input type="text" className="w-full px-4 py-2 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary" placeholder="Enter role name" value={roleName} onChange={(e) => setRoleName(e.target.value)} />
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
