"use client";

import { useState, useEffect, useCallback } from "react";
import { Search, Plus, Trash2, Edit } from "react-feather";
import { Column, CreateUserModal, Table, Alert, AccessTab, RolesTab } from "@/components";
import { useDebounce } from "@/app/hooks/useDebounce";
import { useToast } from "@/app/context/ToastContext";
import { getUsers, deleteUser } from "@/app/api/services/users/users.service";
import { User } from "@/app/api/services/users/users.service.type";

export default function UserPage() {
  const [activeTab, setActiveTab] = useState("Users");
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [totalPages, setTotalPages] = useState(1);
  const [search, setSearch] = useState("");
  
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);

  const toast = useToast();
  const tabs = ["Users", "Access", "Roles"];

  const debouncedSearch = useDebounce(search, 500);

  const fetchUsers = useCallback(async () => {
    try {
      setLoading(true);
      const response = await getUsers({
        page: page,
        show: rowsPerPage,
        search: debouncedSearch,
        sort: "desc"
      });
      
      setUsers(response.data.rows);
      setTotalPages(response.data.total_pages);
    } catch (error: unknown) {
      toast.showToast("Failed to load users", "error");
    } finally {
      setLoading(false);
    }
  }, [page, debouncedSearch, rowsPerPage, toast]);

  useEffect(() => {
    if (activeTab === "Users") {
      fetchUsers();
    }
  }, [fetchUsers, activeTab]);

  const handlePageChange = useCallback((newPage: number) => {
    setPage(newPage);
  }, []);

  const handleRowsPerPageChange = useCallback((newLimit: number) => {
    setRowsPerPage(newLimit);
    setPage(1);
  }, []);

  const handleSearch = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setSearch(e.target.value);
    setPage(1);
  }, []);

  const handleDelete = (user: User) => {
    setSelectedUser(user);
    setIsDeleteAlertOpen(true);
  };

  const confirmDelete = async () => {
    if (!selectedUser) return;

    try {
      await deleteUser(selectedUser.user_id.toString());
      
      toast.showToast("User deleted successfully", "success");
      setIsDeleteAlertOpen(false);
      setSelectedUser(null);
      fetchUsers();
    } catch (error) {
      toast.showToast("Failed to delete user", "error");
    }
  };

  const handleEdit = (user: User) => {
    setSelectedUser(user);
    setIsModalOpen(true);
  };

  const columns: Column<User>[] = [
    { 
      header: "No", 
      accessor: "user_id",
      className: "text-center", 
      headerClassName: "text-center w-16",
      render: (_, index) => <span>{(page - 1) * rowsPerPage + index + 1}</span>
    },
    { header: "Username", accessor: "username", className: "font-medium text-gray-800" },
    { header: "Full Name", accessor: "full_name" },
    { header: "Email", accessor: "email" },
    { header: "Phone", accessor: "phone" },
    { 
      header: "Role IP", 
      accessor: "role_name",
      className: "text-center", 
      headerClassName: "text-center",
      render: (user) => (
        <span className={`px-2 py-1 rounded text-xs font-semibold ${
          user.role_name ? "bg-blue-100 text-blue-800" : "bg-gray-100 text-gray-800"
        }`}>
          {user.role_name ? user.role_name : "End User"}
        </span>
      )
    },
    { 
      header: "Created At", 
      accessor: "created_at",
      className: "whitespace-nowrap",
      render: (user) => new Date(user.created_at).toLocaleDateString('en-GB', {
        day: '2-digit',
        month: 'short',
        year: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      })
    },
    { 
      header: "Actions", 
      headerClassName: "text-center",
      className: "text-center",
      render: (user) => (
        <div className="flex justify-center gap-2">
          <button 
            onClick={() => handleDelete(user)}
            className="p-1.5 text-white bg-red-400 rounded-full hover:bg-red-500 transition-colors shadow-sm cursor-pointer"
          >
            <Trash2 size={14} />
          </button>
          <button 
            onClick={() => handleEdit(user)}
            className="p-1.5 text-white bg-teal-400 rounded-full hover:bg-teal-500 transition-colors shadow-sm cursor-pointer"
          >
            <Edit size={14} />
          </button>
        </div>
      )
    },
  ];

  return (
    <div className="space-y-6">
      
      {/* Tabs */}
      <div className="flex gap-2 border-b border-gray-200 pb-1 bg-white w-fit rounded-full">
        {tabs.map((tab) => (
          <button
            key={tab}
            onClick={() => setActiveTab(tab)}
            className={`px-6 py-2 rounded-full text-sm font-semibold transition-colors cursor-pointer ${
              activeTab === tab
                ? "bg-secondary text-white"
                : "text-gray-500 hover:text-gray-700 hover:bg-gray-100"
            }`}
          >
            {tab}
          </button>
        ))}
      </div>

      {/* Toolbar */}
      <div className="flex flex-col sm:flex-row justify-between gap-4">
        <div className="relative w-full max-w-lg">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={16} />
          <input 
            type="text" 
            placeholder="Search username or email"
            value={search}
            onChange={handleSearch}
            className="w-full pl-10 pr-4 py-2 bg-white border border-gray-200 rounded-full text-sm focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary shadow-sm"
          />
        </div>
        {activeTab === "Users" && (
            <button 
            onClick={() => {
                setSelectedUser(null);
                setIsModalOpen(true);
            }}
            className="flex items-center cursor-pointer gap-2 bg-secondary text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-secondary/90 transition-colors shadow-sm whitespace-nowrap"
            >
            <Plus size={16} />
            Add New User
            </button>
        )}
      </div>

      {/* Table - Only for Users tab */}
      {activeTab === "Users" && (
        <>
          {loading && (
            <div className="flex justify-center items-center py-8">
              <div className="w-8 h-8 border-4 border-gray-300 border-t-blue-500 rounded-full animate-spin" />
            </div>
          )}
          {!loading && (
            <Table 
              columns={columns} 
              data={users}
              keyField="user_id"
              pagination={{
                currentPage: page,
                totalPages: totalPages,
                itemsPerPage: rowsPerPage,
                onPageChange: handlePageChange,
                onRowsPerPageChange: handleRowsPerPageChange
              }}
            />
          )}
        </>
      )}

      {activeTab === "Access" && (
          <AccessTab search={debouncedSearch} />
      )}

      {activeTab === "Roles" && (
          <RolesTab search={debouncedSearch} />
      )}

      <CreateUserModal 
        isOpen={isModalOpen}
        onClose={() => {
            setIsModalOpen(false);
            setSelectedUser(null);
        }}
        onSuccess={() => {
          setIsModalOpen(false);
          setSelectedUser(null);
          fetchUsers();
        }}
        user={selectedUser}
      />

      <Alert 
        isOpen={isDeleteAlertOpen}
        onClose={() => {
          setIsDeleteAlertOpen(false);
          setSelectedUser(null);
        }}
        type="danger"
        title="Delete User"
        message={`Are you sure you want to delete ${selectedUser?.full_name}? This action cannot be undone.`}
        onConfirm={confirmDelete}
        confirmText="Delete"
      />
    </div>
  );
}
