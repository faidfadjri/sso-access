"use client";

import { useState, useEffect, useCallback } from "react";
import { Trash2, Edit, Plus, Navigation, Users } from "react-feather";
import { useDebounce } from "@/app/hooks/useDebounce";
import { useToast } from "@/app/context/ToastContext";
import { getRoles, deleteRole } from "@/app/api/services/roles/roles.service";
import { Role } from "@/app/api/services/roles/roles.service.type";
import {Table, Alert, Column, CreateRoleModal} from "@/components";

import { AssignRoleModal } from "../Modal/AssignRoleModal/AssignRoleModal";
import { authSession } from "@/app/api";
import { ROLES } from "@/app/libs/roles";

interface RolesTabProps {
  search: string;
}

export const RolesTab = ({ search }: RolesTabProps) => {
  const [roles, setRoles] = useState<Role[]>([]);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [totalPages, setTotalPages] = useState(1);
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState(false);
  const [selectedRole, setSelectedRole] = useState<Role | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isAssignRoleModalOpen, setIsAssignRoleModalOpen] = useState(false);
  const toast = useToast();
  const debouncedSearch = useDebounce(search, 500);

  const fetchRoles = useCallback(async () => {
    try {
      setLoading(true);
      const response = await getRoles({ page, show: rowsPerPage, search: debouncedSearch, sort: "desc" });
      setRoles(response.data.rows || []);
      setTotalPages(response.data.total_pages || 1);
    } catch (error) {
      console.error(error);
    } finally {
      setLoading(false);
    }
  }, [page, rowsPerPage, debouncedSearch]);

  useEffect(() => { fetchRoles(); }, [fetchRoles]);

  const handleDelete = (role: Role) => { setSelectedRole(role); setIsDeleteAlertOpen(true); };
  const handleEdit = (role: Role) => { setSelectedRole(role); setIsModalOpen(true); };
  const handleManageUsers = (role: Role) => { setSelectedRole(role); setIsAssignRoleModalOpen(true); };

  const confirmDelete = async () => {
    if (!selectedRole) return;
    try {
      await deleteRole(selectedRole.service_role_id);
      toast.showToast("Role deleted successfully", "success");
      setIsDeleteAlertOpen(false);
      setSelectedRole(null);
      fetchRoles();
    } catch { toast.showToast("Failed to delete role", "error"); }
  };

  const columns: Column<Role>[] = [
    { header: "No", accessor: "service_role_id", className: "text-center", headerClassName: "text-center w-16", render: (_, index) => <span>{(page - 1) * rowsPerPage + index + 1}</span> },
    { header: "Role Name", accessor: "role_name", className: "font-medium text-gray-800" },
    { header: "Service", accessor: "service_id", render: (role) => <span>{role.service_name}</span> },
    { header: "Created At", accessor: "created_at", className: "whitespace-nowrap", render: (role) => new Date(role.created_at).toLocaleDateString("en-GB", { day: "2-digit", month: "short", year: "2-digit", hour: "2-digit", minute: "2-digit" }) },
    {
      header: "Actions", headerClassName: "text-center", className: "text-center",
      render: (role) => (
        <div className="flex justify-center gap-2">
          <button onClick={() => handleManageUsers(role)} className="p-1.5 text-white bg-blue-500 rounded-full hover:bg-blue-600 transition-colors shadow-sm cursor-pointer" title="Manage Users"><Users size={14} /></button>
          {authSession.getRole() == ROLES.SUPER_ADMIN && (
            <>
            <button onClick={() => handleEdit(role)} className="p-1.5 text-white bg-teal-400 rounded-full hover:bg-teal-500 transition-colors shadow-sm cursor-pointer" title="Edit Role"><Edit size={14} /></button>
            <button onClick={() => handleDelete(role)} className="p-1.5 text-white bg-red-400 rounded-full hover:bg-red-500 transition-colors shadow-sm cursor-pointer" title="Delete Role"><Trash2 size={14} /></button>
            </>
          )}
        </div>
      )
    },
  ];

  return (
    <>
      <div className="flex justify-end mb-4 gap-2">
        {authSession.getRole() == ROLES.SUPER_ADMIN && (
          <button onClick={() => { setSelectedRole(null); setIsModalOpen(true); }} className="flex items-center cursor-pointer gap-2 bg-secondary text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-secondary/90 transition-colors shadow-sm whitespace-nowrap">
            <Plus size={16} />Add New Role
          </button>
        )}
        <button onClick={() => { setSelectedRole(null); setIsAssignRoleModalOpen(true); }} className="flex items-center cursor-pointer gap-2 bg-primary text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-primary/90 transition-colors shadow-sm whitespace-nowrap">
          <Navigation size={16} />Assign Role
        </button>
      </div>
      {loading && (<div className="flex justify-center items-center py-8"><div className="w-8 h-8 border-4 border-gray-300 border-t-blue-500 rounded-full animate-spin" /></div>)}
      {!loading && (<Table columns={columns} data={roles} keyField="service_role_id" pagination={{ currentPage: page, totalPages, itemsPerPage: rowsPerPage, onPageChange: setPage, onRowsPerPageChange: (v) => { setRowsPerPage(v); setPage(1); } }} />)}
      <CreateRoleModal isOpen={isModalOpen} onClose={() => { setIsModalOpen(false); setSelectedRole(null); }} onSuccess={() => { setIsModalOpen(false); setSelectedRole(null); fetchRoles(); }} role={selectedRole} />
      <AssignRoleModal isOpen={isAssignRoleModalOpen} onClose={() => { setIsAssignRoleModalOpen(false); setSelectedRole(null); }} onSuccess={() => { setIsAssignRoleModalOpen(false); setSelectedRole(null); fetchRoles(); }} role={selectedRole} />
      <Alert isOpen={isDeleteAlertOpen} onClose={() => { setIsDeleteAlertOpen(false); setSelectedRole(null); }} type="danger" title="Delete Role" message={`Are you sure you want to delete ${selectedRole?.role_name}? This action cannot be undone.`} onConfirm={confirmDelete} confirmText="Delete" />
    </>
  );
};
