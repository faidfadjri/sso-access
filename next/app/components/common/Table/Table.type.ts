import { ReactNode } from "react";

export interface Column<T> {
  header: string | ReactNode;
  accessor?: keyof T; // Optional if using render
  render?: (row: T, index: number) => ReactNode; // Custom render function
  className?: string; // Additional classes for the cell (e.g., text alignment)
  headerClassName?: string;
}

export interface PaginationProps {
  currentPage: number;
  totalPages: number;
  itemsPerPage?: number;
  onPageChange: (page: number) => void;
  onRowsPerPageChange?: (rows: number) => void;
  rowsPerPageOptions?: number[];
}

export interface TableProps<T> {
  columns: Column<T>[];
  data: T[];
  keyField?: keyof T | ((row: T) => string | number);
  pagination?: PaginationProps;
}
