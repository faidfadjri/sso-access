"use client";

import { useState, useEffect, useRef } from "react";
import { Search, X } from "react-feather";
import { useDebounce } from "@/app/hooks/useDebounce";
import { getUsers } from "@/app/api/services/users/users.service";
import { User } from "@/app/api/services/users/users.service.type";

interface UserAutosuggestProps {
  onSelect: (user: User | null) => void;
  selectedUser: User | null;
  className?: string;
}

export const UserAutosuggest = ({ onSelect, selectedUser, className = "" }: UserAutosuggestProps) => {
  const [query, setQuery] = useState("");
  const [options, setOptions] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const [isOpen, setIsOpen] = useState(false);
  const wrapperRef = useRef<HTMLDivElement>(null);
  const debouncedQuery = useDebounce(query, 500);

  useEffect(() => {
    const fetchUsers = async () => {
      if (!debouncedQuery) { setOptions([]); return; }
      try {
        setLoading(true);
        const response = await getUsers({ page: 1, show: 5, search: debouncedQuery, sort: "desc" });
        setOptions(response.data.rows || []);
        setIsOpen(true);
      } catch (error) {
        console.error("Failed to fetch users", error);
        setOptions([]);
      } finally {
        setLoading(false);
      }
    };
    fetchUsers();
  }, [debouncedQuery]);

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (wrapperRef.current && !wrapperRef.current.contains(event.target as Node)) setIsOpen(false);
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const handleSelect = (user: User) => { onSelect(user); setQuery(""); setIsOpen(false); };
  const clearSelection = () => { onSelect(null); setQuery(""); };

  return (
    <div className={`relative ${className}`} ref={wrapperRef}>
      {selectedUser ? (
        <div className="flex items-center justify-between p-2 border border-gray-200 rounded-lg bg-gray-50">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center text-gray-600 font-semibold text-sm">
              {selectedUser.full_name?.substring(0, 2).toUpperCase()}
            </div>
            <div>
              <p className="text-sm font-medium text-gray-900">{selectedUser.full_name}</p>
              <p className="text-xs text-gray-500">{selectedUser.email}</p>
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
            placeholder="Search user by name..."
            value={query}
            onChange={(e) => { setQuery(e.target.value); setIsOpen(true); }}
            onFocus={() => { if (options && options.length > 0) setIsOpen(true); }}
          />
        </div>
      )}

      {isOpen && options && options.length > 0 && !selectedUser && (
        <div className="absolute z-10 w-full mt-1 bg-white border border-gray-200 rounded-lg shadow-lg max-h-60 overflow-y-auto">
          {options.map((user) => (
            <button key={user.user_id} onClick={() => handleSelect(user)} className="w-full flex items-center gap-3 p-2 hover:bg-gray-50 transition-colors text-left">
              <div className="w-8 h-8 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center font-semibold text-xs shrink-0">
                {user.full_name?.substring(0, 2).toUpperCase()}
              </div>
              <div className="overflow-hidden">
                <p className="text-sm font-medium text-gray-900 truncate">{user.full_name}</p>
                <p className="text-xs text-gray-500 truncate">{user.email}</p>
              </div>
            </button>
          ))}
        </div>
      )}

      {isOpen && !loading && (!options || options.length === 0) && query && (
        <div className="absolute z-10 w-full mt-1 bg-white border border-gray-200 rounded-lg shadow-lg p-3 text-center text-sm text-gray-500">
          No users found.
        </div>
      )}
    </div>
  );
};
