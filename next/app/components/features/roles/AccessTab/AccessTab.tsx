"use client";

import { useState, useEffect, useCallback } from "react";
import { Trash2, Plus } from "react-feather";
import { useDebounce } from "@/app/hooks/useDebounce";
import { useToast } from "@/app/context/ToastContext";
import { getUserAccess, deleteUserAccess } from "@/app/api/services/user-access/user-access.service";
import { UserAccess } from "@/app/api/services/user-access/user-access.service.type";
import { Table, Column, Alert } from "@/app/components";
import { CreateAccessModal } from "@/app/components/features/roles/CreateAccessModal/CreateAccessModal";

interface AccessTabProps {
  search: string;
}

export const AccessTab = ({ search }: AccessTabProps) => {
  const [accessList, setAccessList] = useState<UserAccess[]>([]);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [totalPages, setTotalPages] = useState(1);
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState(false);
  const [selectedAccess, setSelectedAccess] = useState<UserAccess | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const toast = useToast();
  const debouncedSearch = useDebounce(search, 500);

  const fetchAccess = useCallback(async () => {
    try {
      setLoading(true);
      const response = await getUserAccess({ page, show: rowsPerPage, search: debouncedSearch, sort: "desc" });
      setAccessList(response.data.rows || []);
      setTotalPages(response.data.total_pages || 1);
    } catch {
      toast.showToast("Failed to load access list", "error");
    } finally {
      setLoading(false);
    }
  }, [page, rowsPerPage, debouncedSearch, toast]);

  useEffect(() => { fetchAccess(); }, [fetchAccess]);

  const handleDelete = (access: UserAccess) => { setSelectedAccess(access); setIsDeleteAlertOpen(true); };
  const confirmDelete = async () => {
    if (!selectedAccess) return;
    try {
      await deleteUserAccess(selectedAccess.access_id);
      toast.showToast("Access revoked successfully", "success");
      setIsDeleteAlertOpen(false);
      setSelectedAccess(null);
      fetchAccess();
    } catch { toast.showToast("Failed to revoke access", "error"); }
  };

  const columns: Column<UserAccess>[] = [
    { header: "No", accessor: "access_id", className: "text-center", headerClassName: "text-center w-16", render: (_, index) => <span>{(page - 1) * rowsPerPage + index + 1}</span> },
    { header: "User Details", accessor: "full_name", render: (access) => (<div className="flex flex-col"><span className="font-medium text-gray-900">{access.full_name}</span><span className="text-xs text-gray-500">{access.email}</span><span className="text-[10px] text-gray-400">@{access.username}</span></div>) },
    { header: "Service", accessor: "service_name", render: (access) => (<div className="flex flex-col"><span className="font-medium text-gray-800">{access.service_name}</span><span className="text-[10px] text-gray-400 font-mono truncate max-w-[150px]" title={access.redirect_url}>{access.redirect_url}</span></div>) },
    { header: "Status", accessor: "status", className: "text-center", headerClassName: "text-center", render: (access) => (<span className={`px-2 py-1 rounded-full text-[10px] font-bold uppercase tracking-wider ${access.status === "active" ? "bg-green-100 text-green-700 border border-green-200" : "bg-red-100 text-red-700 border border-red-200"}`}>{access.status}</span>) },
    { header: "Created At", accessor: "created_at", className: "whitespace-nowrap text-sm text-gray-500", render: (access) => new Date(access.created_at).toLocaleDateString("en-GB", { day: "2-digit", month: "short", year: "2-digit", hour: "2-digit", minute: "2-digit" }) },
    { header: "Actions", headerClassName: "text-center", className: "text-center", render: (access) => (<div className="flex justify-center gap-2"><button onClick={() => handleDelete(access)} className="p-1.5 text-white bg-red-400 rounded-full hover:bg-red-500 transition-colors shadow-sm cursor-pointer" title="Revoke Access"><Trash2 size={14} /></button></div>) },
  ];

  return (
    <>
      <div className="flex justify-end mb-4">
        <button onClick={() => setIsModalOpen(true)} className="flex items-center cursor-pointer gap-2 bg-secondary text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-secondary/90 transition-colors shadow-sm whitespace-nowrap">
          <Plus size={16} />Add New Access
        </button>
      </div>
      {loading ? (
        <div className="flex justify-center items-center py-20">
          <div className="relative">
            <div className="w-12 h-12 border-4 border-gray-100 rounded-full" />
            <div className="w-12 h-12 border-4 border-t-secondary border-transparent rounded-full animate-spin absolute top-0 left-0" />
          </div>
        </div>
      ) : (
        <Table columns={columns} data={accessList} keyField="access_id" pagination={{ currentPage: page, totalPages, itemsPerPage: rowsPerPage, onPageChange: setPage, onRowsPerPageChange: (v) => { setRowsPerPage(v); setPage(1); } }} />
      )}
      <CreateAccessModal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} onSuccess={() => { setIsModalOpen(false); fetchAccess(); }} />
      <Alert isOpen={isDeleteAlertOpen} onClose={() => { setIsDeleteAlertOpen(false); setSelectedAccess(null); }} type="danger" title="Revoke Access" message="Are you sure you want to revoke access? This action cannot be undone." onConfirm={confirmDelete} confirmText="Revoke" />
    </>
  );
};
