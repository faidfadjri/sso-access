"use client";

import { ChevronLeft, ChevronRight } from "react-feather";
import { TableProps } from "./Table.type";

export default function Table<T extends Record<string, any>>({
  columns,
  data = [],
  keyField = "id",
  pagination,
}: TableProps<T>) {
  
  const getKey = (row: T, index: number) => {
    if (typeof keyField === "function") {
      return keyField(row);
    }
    return row[keyField] || index;
  };

  return (
    <div className="bg-white rounded-xl shadow-sm overflow-hidden border border-gray-100">
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead className="bg-secondary text-white text-xs uppercase font-medium">
            <tr>
              {columns.map((col, index) => (
                <th 
                  key={index} 
                  className={`px-6 py-4 tracking-wider ${col.headerClassName || "text-left"}`}
                >
                  {col.header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-50">
            {data.map((row, rowIndex) => (
              <tr 
                key={getKey(row, rowIndex)} 
                className="hover:bg-gray-50/50 transition-colors text-sm text-gray-600"
              >
                {columns.map((col, colIndex) => (
                  <td 
                    key={colIndex} 
                    className={`px-6 py-4 ${col.className || ""}`}
                  >
                    {col.render 
                      ? col.render(row, rowIndex) 
                      : (col.accessor ? row[col.accessor] : "")
                    }
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      
      {/* Pagination */}
      {pagination && (
        <div className="flex justify-end p-4 border-t border-gray-100">
          <div className="flex items-center gap-4">
            {pagination.onRowsPerPageChange && (
              <div className="flex items-center gap-2 text-sm text-gray-500">
                <span>Show</span>
                <select
                  value={pagination.itemsPerPage}
                  onChange={(e) => pagination.onRowsPerPageChange?.(Number(e.target.value))}
                  className="bg-white border border-gray-200 rounded px-2 py-1 focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary cursor-pointer"
                >
                  {(pagination.rowsPerPageOptions || [5, 10, 20]).map((option) => (
                    <option key={option} value={option}>
                      {option}
                    </option>
                  ))}
                </select>
                <span>entries</span>
              </div>
            )}
            <div className="flex items-center gap-1">
            <button 
                onClick={() => pagination.onPageChange(pagination.currentPage - 1)}
                disabled={pagination.currentPage === 1}
                className="w-8 h-8 flex items-center justify-center rounded border border-gray-200 text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
            >
              <ChevronLeft size={16} />
            </button>
            
            {Array.from({ length: Math.min(5, pagination.totalPages) }, (_, i) => {
                // Simple pagination logic: show current page and neighbors, or just first 5 for now
                // Improving logic to show window around current page could be better but sticking to simple for now
                // Actually, let's just show all pages if small, or complex logic. 
                // For simplicity, let's just show simple range or use a pagination library if available. 
                // implementing a simple window logic: 
                let pageNum = i + 1;
                if (pagination.totalPages > 5) {
                    if (pagination.currentPage > 3) {
                        pageNum = pagination.currentPage - 2 + i;
                    }
                    if (pageNum > pagination.totalPages) {
                        pageNum = pagination.totalPages - (4 - i);
                    }
                }
                
                return (
                    <button 
                        key={pageNum}
                        onClick={() => pagination.onPageChange(pageNum)}
                        className={`w-8 h-8 flex items-center justify-center rounded border text-xs font-medium cursor-pointer ${
                            pagination.currentPage === pageNum 
                                ? "border-secondary bg-secondary text-white" 
                                : "border-gray-200 text-gray-500 hover:bg-gray-50"
                        }`}
                    >
                        {pageNum}
                    </button>
                )
            })}

            <button 
                onClick={() => pagination.onPageChange(pagination.currentPage + 1)}
                disabled={pagination.currentPage === pagination.totalPages}
                className="w-8 h-8 flex items-center justify-center rounded border border-gray-200 text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
            >
              <ChevronRight size={16} />
            </button>
          </div>
          </div>
        </div>
      )}
    </div>
  );
}
