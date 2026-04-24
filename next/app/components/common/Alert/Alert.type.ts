import { ReactNode } from "react";

export type AlertType = "success" | "danger" | "warning" | "info";

export interface AlertProps {
  isOpen: boolean;
  onClose: () => void;
  type?: AlertType;
  title: string;
  message?: ReactNode;
  onConfirm?: () => void;
  confirmText?: string;
  cancelText?: string;
  showCancel?: boolean;
}
