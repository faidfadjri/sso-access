import { Client } from "@/app/api/services/clients/client.type";
import { PaginationProps } from "@/app/components/common/Table/Table.type";

export type ServiceTableProps = {
  data: Client[];
  isLoading?: boolean;
  pagination?: PaginationProps;
  onSearch?: (query: string) => void;
  onRefresh?: () => void;
  onDelete?: (client: Client) => void;
  onDetail?: (client: Client) => void;
  onAdd?: () => void;
  onEdit?: (client: Client) => void;
  rowsPerPage?: number;
  onRowsPerPageChange?: (rows: number) => void;
}
