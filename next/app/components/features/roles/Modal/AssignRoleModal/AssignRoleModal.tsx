"use client";

import { useState, useEffect, useCallback } from "react";
import { X, Trash2, User as UserIcon } from "react-feather";
import { Modal } from "@/app/components/common";
import { User } from "@/app/api/services/users/users.service.type";
import { Role } from "@/app/api/services/roles/roles.service.type";
import { Assignation } from "@/app/api/services/roles/assignation.service.type";
import { getAssignations, createAssignation, deleteAssignation } from "@/app/api/services/roles/assignation.service";
import { RoleAutosuggest } from "../../RoleAutosuggest/RoleAutosuggest";
import { UserAutosuggest } from "../../../users/UserAutosuggest/UserAutosuggest";
import { useToast } from "@/app/context/ToastContext";

interface AssignRoleModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
  role?: Role | null;
}

export const AssignRoleModal = ({ isOpen, onClose, onSuccess, role }: AssignRoleModalProps) => {
  const [selectedRole, setSelectedRole] = useState<Role | null>(null);
  const [selectedUsers, setSelectedUsers] = useState<User[]>([]);
  const [existingAssignations, setExistingAssignations] = useState<Assignation[]>([]);
  const [loading, setLoading] = useState(false);
  const [loadingAssignations, setLoadingAssignations] = useState(false);
  const toast = useToast();

  const fetchAssignations = useCallback(async (roleId: number, serviceId: number) => {
    try {
        setLoadingAssignations(true);
        const data = await getAssignations({ service_id: serviceId, role_id: roleId, show: -1 });
        setExistingAssignations(data.data.rows || []);
    } catch (error) {
        console.error("Failed to fetch assignations", error);
        toast.showToast("Failed to fetch current assignments", "error");
    } finally {
        setLoadingAssignations(false);
    }
  }, [toast]);

  useEffect(() => {
    if (isOpen) {
        if (role) {
            setSelectedRole(role);
            fetchAssignations(role.service_role_id, role.service_id);
        } else {
            setSelectedRole(null);
            setExistingAssignations([]);
        }
        setSelectedUsers([]);
        setLoading(false);
    }
  }, [isOpen, role, fetchAssignations]);

  // When selectedRole changes manually (if not locked by prop)
  useEffect(() => {
    if (selectedRole && !role) {
        fetchAssignations(selectedRole.service_role_id, selectedRole.service_id);
    } else if (!selectedRole) {
        setExistingAssignations([]);
    }
  }, [selectedRole, role, fetchAssignations]);

  const handleAddUser = (user: User | null) => {
    if (!user) return;
    // Check if pending in selectedUsers
    if (selectedUsers.some(u => u.user_id === user.user_id)) {
        toast.showToast("User already selected", "warning");
        return;
    }
    // Check if already assigned
    if (existingAssignations.some(a => a.user_id === user.user_id)) {
        toast.showToast(`User is already assigned to this role`, "warning");
        return;
    }
    setSelectedUsers([...selectedUsers, user]);
  };

  const handleRemoveUser = (userId: number) => {
    setSelectedUsers(selectedUsers.filter(u => u.user_id !== userId));
  };

  const handleDeleteAssignation = async (assignation: Assignation) => {
    if (!confirm(`Are you sure you want to unassign ${assignation.full_name}?`)) return;
    
    try {
        await deleteAssignation({
            user_id: assignation.user_id,
            role_id: assignation.service_role_id,
            service_id: assignation.service_id,
            // @ts-ignore - The API definition might include these wrapper fields in type but payload is flat, or vice versa. 
            // Based on usage in service: const response = await fetch(\`/api/services/${assignation.service_id}/roles/${assignation.role_id}/assignations/${assignation.user_id}\`
            // It seems 'deleteAssignation' function might not even use the body for the URL parameters, but let's pass the object matching the type.
        } as any); 

        toast.showToast("User unassigned successfully", "success");
        // Remove from list
        setExistingAssignations(prev => prev.filter(a => a.user_id !== assignation.user_id));
    } catch (error) {
        console.error("Failed to delete assignation", error);
        toast.showToast("Failed to unassign user", "error");
    }
  };

  const handleSave = async () => {
    if (!selectedRole) {
        toast.showToast("Please select a role", "error");
        return;
    }
    if (selectedUsers.length === 0) {
        toast.showToast("Please select at least one user to assign", "error");
        return;
    }

    try {
        setLoading(true);
        // Using createAssignation which accepts multiple user_ids
        await createAssignation({
            role_id: selectedRole.service_role_id,
            user_ids: selectedUsers.map(u => u.user_id),
            service_id: selectedRole.service_id
        } as any); // Type assertion if the definition is strictly BaseResponse wrapped
        
        toast.showToast("Users assigned successfully", "success");
        // Refresh list and clear selection
        setSelectedUsers([]);
        fetchAssignations(selectedRole.service_role_id, selectedRole.service_id);
    } catch (error) {
        console.error("Failed to assign roles", error);
        toast.showToast("Failed to assign roles", "error");
    } finally {
        setLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={
        <div className="mb-0">
            <h2 className="text-xl font-bold text-gray-900">{role ? `Manage Users for ${role.role_name}` : "Assign Role"}</h2>
            <p className="text-sm text-gray-500 mt-1 font-normal">Manage user assignments for this role</p>
        </div>
      }
      maxWidth="max-w-2xl"
    >
        <div className="space-y-6 mt-6">
            {/* Role Selection (ReadOnly if role prop provided) */}
            <div className="space-y-2">
                <label className="text-sm font-bold text-gray-700">Role <span className="text-red-500">*</span></label>
                {role ? (
                     <div className="p-3 bg-gray-50 border border-gray-200 rounded-lg text-gray-700 font-medium">
                        {role.role_name} <span className="text-gray-400 text-sm ml-2">({role.service_name})</span>
                     </div>
                ) : (
                    <RoleAutosuggest 
                        className="mt-2"
                        selectedRole={selectedRole} 
                        onSelect={setSelectedRole} 
                        placeholder="Select a role..."
                    />
                )}
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {/* Existing Assignments Column */}
                <div className="space-y-4">
                    <label className="text-sm font-bold text-gray-700 flex items-center gap-2">
                        Existing Assignments
                        <span className="bg-gray-100 text-gray-600 px-2 py-0.5 rounded-full text-xs">{existingAssignations.length}</span>
                    </label>
                    <div className="bg-gray-50 border border-gray-200 rounded-lg h-64 overflow-y-auto p-2 space-y-2">
                        {loadingAssignations ? (
                            <div className="flex justify-center items-center h-full">
                                <div className="w-5 h-5 border-2 border-gray-300 border-t-blue-500 rounded-full animate-spin" />
                            </div>
                        ) : existingAssignations.length > 0 ? (
                            existingAssignations.map(assignation => (
                                <div key={assignation.user_id} className="flex items-center justify-between p-2 pl-3 bg-white border border-gray-100 rounded-md shadow-sm group">
                                    <div className="flex items-center gap-2 overflow-hidden">
                                        <div className="bg-blue-100 text-blue-600 p-1 rounded-full"><UserIcon size={12} /></div>
                                        <div className="truncate">
                                            <p className="text-sm font-medium text-gray-700 truncate">{assignation.full_name}</p>
                                        </div>
                                    </div>
                                    <button 
                                        onClick={() => handleDeleteAssignation(assignation)}
                                        className="text-gray-400 hover:text-red-500 hover:bg-red-50 p-1.5 rounded-md transition-all opacity-0 group-hover:opacity-100"
                                        title="Unassign User"
                                    >
                                        <Trash2 size={14} />
                                    </button>
                                </div>
                            ))
                        ) : (
                            <div className="h-full flex flex-col items-center justify-center text-gray-400 text-sm">
                                <p>No users assigned yet</p>
                            </div>
                        )}
                    </div>
                </div>

                {/* Add New Users Column */}
                <div className="space-y-4">
                    <div className="space-y-2">
                        <label className="text-sm font-bold text-gray-700">Add New Users</label>
                        <UserAutosuggest 
                            className="mt-2"
                            selectedUser={null}
                            onSelect={handleAddUser}
                        />
                    </div>

                    <div className="space-y-2">
                        <label className="text-sm font-semibold text-gray-600">Selected to Assign</label>
                        {selectedUsers.length > 0 ? (
                            <div className="bg-white border border-gray-200 rounded-lg max-h-40 overflow-y-auto p-1 space-y-1">
                                {selectedUsers.map(user => (
                                    <div key={user.user_id} className="flex items-center justify-between p-2 bg-green-50 border border-green-100 rounded-md">
                                        <span className="text-sm font-medium text-gray-700">{user.full_name}</span>
                                        <button 
                                            onClick={() => handleRemoveUser(user.user_id)}
                                            className="text-gray-400 hover:text-red-500 transition-colors duration-100 p-1"
                                        >
                                            <X size={14} />
                                        </button>
                                    </div>
                                ))}
                            </div>
                        ) : (
                            <div className="p-4 border-2 border-dashed border-gray-200 rounded-lg text-center text-gray-400 text-sm">
                                Search and select users to add them here
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </div>

        <div className="flex justify-end gap-3 mt-8 pt-4 border-t border-gray-100">
            <button 
                onClick={onClose}
                disabled={loading}
                className="px-4 py-2 text-sm font-bold text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors disabled:opacity-50"
            >
                Close
            </button>
            <button 
                onClick={handleSave}
                disabled={loading || selectedUsers.length === 0}
                className="px-4 py-2 text-sm font-bold text-white bg-[#1e293b] rounded-lg hover:bg-[#334155] transition-colors disabled:opacity-50 flex items-center gap-2"
            >
                {loading && <span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />}
                Assign {selectedUsers.length > 0 ? `(${selectedUsers.length})` : ''}
            </button>
        </div>
    </Modal>
  );
};
