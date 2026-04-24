"use client";

import { Search, Plus, Trash2, Info, ChevronUp, ChevronDown, Edit } from "react-feather";
import { useEffect, useState } from "react";
import { useDebounce } from "@/app/hooks/useDebounce";
import { Client } from "@/app/api/services/clients/client.type";
import { ServiceTableProps } from "./ServiceTable.type";
import { Column } from "@/app/components/common/Table/Table.type";
import { Table } from "@/app/components/common";

export default function ServiceTable({ data, pagination, onSearch, onDelete, onDetail, rowsPerPage, onRowsPerPageChange, isLoading, onAdd, onEdit }: ServiceTableProps) {
  const [searchQuery, setSearchQuery] = useState("");
  const [sortConfig, setSortConfig] = useState<{ key: string; direction: "asc" | "desc" } | null>(null);

  const debouncedSearchQuery = useDebounce(searchQuery, 500);

  useEffect(() => {
    onSearch?.(debouncedSearchQuery);
  }, [debouncedSearchQuery, onSearch]);

  const handleSort = (key: string) => {
    const direction: "asc" | "desc" = sortConfig?.key === key && sortConfig.direction === "asc" ? "desc" : "asc";
    setSortConfig({ key, direction });
  };

  const sortedData = [...data].sort((a, b) => {
    if (!sortConfig) return 0;
    const { key, direction } = sortConfig;
    const aValue = a[key as keyof Client];
    const bValue = b[key as keyof Client];
    if (aValue < bValue) return direction === "asc" ? -1 : 1;
    if (aValue > bValue) return direction === "asc" ? 1 : -1;
    return 0;
  });

  const columns: Column<Client>[] = [
    {
      header: "No.",
      render: (_, index) => {
        const page = pagination?.currentPage || 1;
        const limit = pagination?.itemsPerPage || 10;
        return (page - 1) * limit + index + 1;
      },
      className: "text-center w-16 text-gray-500",
      headerClassName: "text-center w-16",
    },
    { header: "Service Name", accessor: "service_name", className: "font-medium text-gray-800" },
    {
      header: "Client ID",
      accessor: "client_id",
      render: (row) => <div title={row.client_id}>{row.client_id?.length > 30 ? `${row.client_id.substring(0, 30)}...` : row.client_id}</div>,
    },
    {
      header: "Redirect URL",
      accessor: "redirect_url",
      render: (row) => <div className="truncate max-w-[200px]" title={row.redirect_url}>{row.redirect_url}</div>,
    },
    {
      header: "Status",
      accessor: "is_active",
      headerClassName: "text-center",
      className: "text-center",
      render: (row) => (
        <span className={`inline-flex items-center px-2.5 py-0.5 rounded text-xs font-medium ${row.is_active ? "bg-emerald-100 text-emerald-800" : "bg-red-100 text-red-800"}`}>
          {row.is_active ? "Active" : "Inactive"}
        </span>
      ),
    },
    {
      header: (
        <div className="flex items-center gap-1 cursor-pointer select-none group" onClick={() => handleSort("created_at")}>
          Created At
          <div className="flex flex-col -space-y-1">
            <ChevronUp size={12} className={`${sortConfig?.key === "created_at" && sortConfig.direction === "asc" ? "text-gray-800 stroke-[3px]" : "text-gray-300 group-hover:text-gray-400"}`} />
            <ChevronDown size={12} className={`${sortConfig?.key === "created_at" && sortConfig.direction === "desc" ? "text-gray-800 stroke-[3px]" : "text-gray-300 group-hover:text-gray-400"}`} />
          </div>
        </div>
      ),
      accessor: "created_at",
      className: "whitespace-nowrap",
      render: (row) => new Date(row.created_at).toLocaleString("en-GB", { day: "2-digit", month: "short", year: "numeric", hour: "2-digit", minute: "2-digit", hour12: false }).replace(",", ""),
    },
    {
      header: "Actions",
      headerClassName: "text-center",
      className: "text-center",
      render: (row) => (
        <div className="flex justify-center gap-2">
          <button onClick={() => onDelete?.(row)} className="p-1.5 text-white bg-red-400 rounded-full hover:bg-red-500 transition-colors shadow-sm cursor-pointer" title="Delete Service">
            <Trash2 size={14} />
          </button>
          <button onClick={() => onDetail?.(row)} className="p-1.5 text-white bg-sky-500 rounded-full hover:bg-sky-600 transition-colors shadow-sm cursor-pointer" title="View Details">
            <Info size={14} />
          </button>
          <button onClick={() => onEdit?.(row)} className="p-1.5 text-white bg-teal-500 rounded-full hover:bg-teal-600 transition-colors shadow-sm cursor-pointer" title="Edit Service">
            <Edit size={14} />
          </button>
        </div>
      ),
    },
  ];

  return (
    <div className="space-y-6">
      <h2 className="text-xl font-bold text-gray-800">Service Credentials</h2>
      <div className="flex flex-col sm:flex-row justify-between gap-4">
        <div className="flex items-center gap-2">
          <div className="relative w-full max-w-xs">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={16} />
            <input
              type="text"
              placeholder="Search Services"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-10 pr-4 py-2 bg-white border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
            />
          </div>
          {isLoading && <div className="w-5 h-5 flex-none shrink-0 border-2 border-gray-300 border-t-blue-500 rounded-full animate-spin" />}
        </div>
        <button onClick={onAdd} className="flex items-center cursor-pointer gap-2 bg-secondary text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-secondary/90 transition-colors shadow-sm whitespace-nowrap">
          <Plus size={16} />
          Create New Credentials
        </button>
      </div>
      <Table columns={columns} data={sortedData} keyField="service_id" pagination={{ ...pagination!, onRowsPerPageChange, itemsPerPage: rowsPerPage }} />
    </div>
  );
}
