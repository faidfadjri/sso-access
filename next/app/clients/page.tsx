'use client'

import { useEffect, useState, useCallback } from "react";
import { useToast } from "@/app/context/ToastContext";
import { Client, deleteServiceClient, getServiceClients } from "../api";
import { Alert, CredentialModal, ResultCredentialModal, ServiceTable } from "../components";


export default function ServicesPage() {
  const [clients, setClients] = useState<Client[]>([]);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(5);
  const [totalPages, setTotalPages] = useState(1);
  const [search, setSearch] = useState("");

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isDetailModalOpen, setIsDetailModalOpen] = useState(false);
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState(false);
  const [selectedClient, setSelectedClient] = useState<Client | null>(null);

  const toast = useToast();

  const fetchClients = useCallback(async () => {
    try {
      setLoading(true);
      const response = await getServiceClients({
        page: page,
        show: rowsPerPage,
        search: search,
        sort: "desc"
      });
      
      setClients(response.data.rows);
      setTotalPages(response.data.total_pages);
    } catch (error: unknown) {
      console.log(error)
      toast.showToast("Failed to load clients", "error");
    } finally {
      setLoading(false);
    }
  }, [page, search, rowsPerPage, toast]);

  useEffect(() => {
    fetchClients();
  }, [fetchClients]);

  const handlePageChange = useCallback((newPage: number) => {
    setPage(newPage);
  }, []);

  const handleRowsPerPageChange = useCallback((newLimit: number) => {
    setRowsPerPage(newLimit);
    setPage(1);
  }, []);

  const handleSearch = useCallback((query: string) => {
    setSearch(query);
    setPage(1);
  }, []);

  const handleDelete = (client: Client) => {
    setSelectedClient(client);
    setIsDeleteAlertOpen(true);
  };

  const confirmDelete = async () => {
    if (!selectedClient) return;

    try {
      await deleteServiceClient(selectedClient.service_id);
      
      toast.showToast("Service client deleted successfully", "success");
      setIsDeleteAlertOpen(false);
      fetchClients();
    } catch (error: unknown) {
       console.error("Delete error:", error);
       toast.showToast("Failed to delete service client", "error");
    }
  };

  const handleDetail = (client: Client) => {
    setSelectedClient(client);
    setIsDetailModalOpen(true);
  };

  const handleEdit = (client: Client) => {
    setSelectedClient(client);
    setIsModalOpen(true);
  }

  const handleAdd = () => {
    setSelectedClient(null);
    setIsModalOpen(true);
  }

  return (
    <div className="max-w-7xl mx-auto">
      <CredentialModal 
        mode={selectedClient ? "update" : "create"}
        isOpen={isModalOpen} 
        data={selectedClient}
        onClose={() => setIsModalOpen(false)}
        onSuccess={() => {
            setIsModalOpen(false);
            fetchClients();
        }}
      />
      
      <ServiceTable 
        data={clients} 
        isLoading={loading}
        pagination={{
            currentPage: page,
            totalPages: totalPages,
            itemsPerPage: rowsPerPage,
            onPageChange: handlePageChange
        }}
        onSearch={handleSearch}
        onRefresh={fetchClients}
        onDelete={handleDelete}
        onDetail={handleDetail}
        rowsPerPage={rowsPerPage}
        onRowsPerPageChange={handleRowsPerPageChange}
        onAdd={handleAdd}
        onEdit={handleEdit}
      />

      <Alert
        isOpen={isDeleteAlertOpen}
        onClose={() => setIsDeleteAlertOpen(false)}
        type="danger"
        title="Delete Service Client"
        message={`Are you sure you want to delete ${selectedClient?.service_name}? This action cannot be undone.`}
        onConfirm={confirmDelete}
        confirmText="Delete"
      />
      
      <ResultCredentialModal
        isOpen={isDetailModalOpen}
        ServiceClient={selectedClient}
        onClose={() => {
          setIsDetailModalOpen(false)
        }}
      />
    </div>
  );
}

