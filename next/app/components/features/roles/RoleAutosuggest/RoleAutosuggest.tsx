"use client";

import { useState, useEffect, useRef } from "react";
import { Search, X, Briefcase } from "react-feather";
import { useDebounce } from "@/app/hooks/useDebounce";
import { getRoles } from "@/app/api/services/roles/roles.service";
import { Role } from "@/app/api/services/roles/roles.service.type";

interface RoleAutosuggestProps {
  onSelect: (role: Role | null) => void;
  selectedRole: Role | null;
  className?: string;
  placeholder?: string;
}

export const RoleAutosuggest = ({ onSelect, selectedRole, className = "", placeholder = "Search role..." }: RoleAutosuggestProps) => {
  const [query, setQuery] = useState("");
  const [options, setOptions] = useState<Role[]>([]);
  const [loading, setLoading] = useState(false);
  const [isOpen, setIsOpen] = useState(false);
  const wrapperRef = useRef<HTMLDivElement>(null);
  const debouncedQuery = useDebounce(query, 500);

  useEffect(() => {
    const fetchRoles = async () => {
      try {
        setLoading(true); 
        const response = await getRoles({ page: 1, show: 5, search: debouncedQuery, sort: "desc" });
        setOptions(response.data.rows || []);
        if (debouncedQuery || isOpen) setIsOpen(true);
      } catch (error) {
        console.error("Failed to fetch roles", error);
      } finally {
        setLoading(false);
      }
    };
    
    if (isOpen || debouncedQuery) {
        fetchRoles();
    }
  }, [debouncedQuery, isOpen]);

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (wrapperRef.current && !wrapperRef.current.contains(event.target as Node)) setIsOpen(false);
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const handleSelect = (role: Role) => { onSelect(role); setQuery(""); setIsOpen(false); };
  const clearSelection = () => { onSelect(null); setQuery(""); };

  return (
    <div className={`relative ${className}`} ref={wrapperRef}>
      {selectedRole ? (
        <div className="flex items-center justify-between p-2 border border-gray-200 rounded-lg bg-gray-50">
          <div className="flex items-center gap-3">
             <div className="p-2 bg-purple-100 text-purple-600 rounded-full shrink-0">
                <Briefcase size={16} />
             </div>
            <div>
              <div className="flex flex-col gap-1">
                <p className="text-xs text-gray-500 truncate">{selectedRole.service_name}</p>
                <p className="text-sm font-medium text-gray-900 truncate">{selectedRole.role_name}</p>
              </div>
            </div>
          </div>
          <button type="button" onClick={clearSelection} className="p-1 text-gray-400 hover:text-gray-600 rounded-full hover:bg-gray-200"><X size={16} /></button>
        </div>
      ) : (
        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={16} />
          <input
            type="text"
            className="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
            placeholder={placeholder}
            value={query}
            onChange={(e) => { setQuery(e.target.value); setIsOpen(true); }}
            onFocus={() => setIsOpen(true)}
          />
        </div>
      )}

      {isOpen && options.length > 0 && !selectedRole && (
        <div className="absolute z-10 w-full mt-1 bg-white border border-gray-200 rounded-lg shadow-lg max-h-60 overflow-y-auto">
          {options.map((role) => (
            <button key={role.service_role_id} onClick={() => handleSelect(role)} className="w-full flex items-center gap-3 p-2 hover:bg-gray-50 transition-colors text-left">
              <div className="p-2 bg-purple-100 text-purple-600 rounded-full shrink-0">
                 <Briefcase size={14} />
              </div>
              <div className="overflow-hidden">
                <div className="flex flex-col gap-1">
                   <p className="text-xs text-gray-500 truncate">{role.service_name}</p>
                   <p className="text-sm font-medium text-gray-900 truncate">{role.role_name}</p>
                </div>
              </div>
            </button>
          ))}
        </div>
      )}
      
       {isOpen && !loading && options.length === 0 && query && (
        <div className="absolute z-10 w-full mt-1 bg-white border border-gray-200 rounded-lg shadow-lg p-3 text-center text-sm text-gray-500">
          No roles found.
        </div>
      )}
    </div>
  );
};
