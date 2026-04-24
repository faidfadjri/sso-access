import { Client } from "@/app/api/services/clients/client.type";

export interface ResultCredentialModalProps {
  isOpen: boolean;
  onClose: () => void;
  ServiceClient?: Client | null;
}
