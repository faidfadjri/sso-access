import { Client } from "@/app/api";


type Mode = "create" | "update"

export interface CredentialModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
  data?: Client|null;
  mode: Mode
}


export type formValues = {
  service_id?: number | null;
  name: string;
  redirect_url: string;
  description: string;
  is_active?: boolean;
}
