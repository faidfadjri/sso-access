"use client";

import { Search, ChevronDown, User, LogOut, File, FileText } from "react-feather";
import { useState, useRef, useEffect } from "react";
import { useTopbar } from "@/app/context/TopbarContext";
import { TopbarProps } from "./Topbar.type";
import { authSession } from "@/app/api/services/oauth/oauth-session.service";
import { oauthLogout } from "@/app/api/services/oauth/oauth.service";
import { useDebounce } from "@/app/hooks/useDebounce";
import { getUserAccess } from "@/app/api/services/user-access/user-access.service";
import { UserAccess } from "@/app/api/services/user-access/user-access.service.type";
import { localStorageLib } from "@/app/libs/local-storage";
import Link from "next/link";

const getInitials = (name?: string | null) => {
  if (!name) return "";
  const parts = name.trim().split(/\s+/);
  if (parts.length === 1) return parts[0].substring(0, 2).toUpperCase();
  return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
};

export default function Topbar({ title: propTitle }: TopbarProps) {
  const { title: contextTitle } = useTopbar();
  const title = propTitle || contextTitle;
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const searchRef = useRef<HTMLDivElement>(null);
  const [mounted, setMounted] = useState(false);
  const [imgError, setImgError] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const debouncedSearch = useDebounce(searchQuery, 300);
  const [searchResults, setSearchResults] = useState<UserAccess[]>([]);
  const [isSearchOpen, setIsSearchOpen] = useState(false);

  useEffect(() => {
    if (!debouncedSearch) {
      setSearchResults([]);
      setIsSearchOpen(false);
      return;
    }
    const fetchResults = async () => {
      try {
        const res = await getUserAccess({ search: debouncedSearch, page: 1, show: 5, user_id: authSession.getUserId() });
        setSearchResults(res.data.rows || []);
        setIsSearchOpen(true);
      } catch (err) {
        console.error(err);
      }
    };
    fetchResults();
  }, [debouncedSearch]);

  const handleAppClick = (app: UserAccess) => {
    const storedRecent = localStorageLib.get<UserAccess[]>("recentApps") || [];
    const filtered = storedRecent.filter((a) => a.access_id !== app.access_id);
    const updatedRecent = [app, ...filtered].slice(0, 5);
    localStorageLib.set("recentApps", updatedRecent);
    
    if (app.redirect_url) {
      window.location.href = app.redirect_url;
    }
    setIsSearchOpen(false);
    setSearchQuery("");
  };

  const getIconPath = (name: string) => {
    if (!name) return "/apps/website.png";
    const lowerName = name.toLowerCase();
    if (lowerName.includes("sivalet")) return "/apps/sivalet.png";
    if (lowerName.includes("promotive")) return "/apps/promotive.png";
    if (lowerName.includes("attendify")) return "/apps/attendify.png";
    if (lowerName.includes("parking")) return "/apps/parking.png";
    return "/apps/website.png";
  };
  useEffect(() => {
    setMounted(true);
    function handleClickOutside(event: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsDropdownOpen(false);
      }
      if (searchRef.current && !searchRef.current.contains(event.target as Node)) {
        setIsSearchOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  const logout = async () => {
    try {
      await oauthLogout();
    } catch(err) {
      console.error(err);
    } finally {
      authSession.clear();
      window.location.href = "/login";
    }
  }

  return (
    <header className="flex flex-col md:flex-row md:items-center justify-between gap-4 py-6 px-8 bg-white/50 backdrop-blur-sm sticky top-0 z-40 w-full">
      <div className="relative w-full max-w-md" ref={searchRef}>
        {title ? (
           <div className="font-p text-gray-400 font-medium">
             {title}
           </div>
        ) : (
          <>
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <Search size={18} className="text-gray-400" />
            </div>
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              onFocus={() => {
                if (searchResults.length > 0) setIsSearchOpen(true);
              }}
              suppressHydrationWarning
              className="block w-full pl-10 pr-3 py-2.5 border border-gray-200 rounded-full leading-5 bg-white text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500 sm:text-sm shadow-sm transition-shadow"
              placeholder="Search Your App"
            />
            {isSearchOpen && searchResults.length > 0 && (
              <div className="absolute top-12 left-0 w-full bg-white rounded-lg shadow-lg border border-gray-100 z-50 py-2 max-h-60 overflow-y-auto">
                {searchResults.map((app) => (
                  <div 
                    key={app.access_id}
                    onClick={() => handleAppClick(app)}
                    className="flex items-center gap-3 px-4 py-2 hover:bg-gray-50 cursor-pointer transition-colors"
                  >
                    <img src={getIconPath(app.service_name)} alt={app.service_name} className="w-8 h-8 rounded-full object-cover" />
                    <span className="text-sm font-medium text-gray-700">{app.service_name}</span>
                  </div>
                ))}
              </div>
            )}
          </>
        )}
      </div>

      <div className="flex items-center justify-between md:justify-end gap-6 w-full md:w-auto">

        <div className="text-sm font-medium text-gray-600 flex items-center gap-2">
          Have a problem? <a href="https://wa.me/+6285848238397" target="_blank" className="text-primary-dark hover:underline">Ask for helps</a> | <Link href="/docs" className="text-primary-dark hover:underline flex items-center gap-1"><FileText size={16}/> Docs</Link>
        </div>

        <div className="flex items-center gap-3 pl-6 border-l border-gray-200 relative" ref={dropdownRef}>
          <div className="hidden md:flex flex-col items-end">
            {mounted && authSession.get() && (
              <>
                <span className="text-sm font-bold text-gray-800 leading-tight">{authSession.getFullName()}</span>
                <span className="text-xs font-semibold text-white bg-[#0ea5e9] px-1.5 py-0.5 rounded leading-tight">{authSession.getRole()}</span>
              </>
            )}
          </div>

          <div
            className="relative group cursor-pointer flex items-center gap-1"
            onClick={() => setIsDropdownOpen(!isDropdownOpen)}
          >
            <div className="h-10 w-10 relative rounded-full overflow-hidden border-2 border-white shadow-sm bg-yellow-400 flex items-center justify-center text-yellow-800 font-bold select-none">
              {mounted && authSession.getPhoto() && !imgError ? (
                <img 
                  src={`${process.env.NEXT_PUBLIC_BASE_API_URL}${authSession.getPhoto()}`} 
                  alt={authSession.getFullName()} 
                  className="h-full w-full object-cover"
                  onError={() => setImgError(true)}
                />
              ) : (
                <span>{mounted ? getInitials(authSession.getFullName()) : ""}</span>
              )}
            </div>
            <ChevronDown size={16} className={`text-gray-400 transition-transform duration-200 ${isDropdownOpen ? "rotate-180" : ""}`} />
          </div>

          {isDropdownOpen && (
            <div className="absolute top-12 right-0 w-48 bg-white rounded-lg shadow-lg py-1 border border-gray-100 z-50 animate-in fade-in zoom-in-95 duration-100">
              <a href="/account" className="flex items-center gap-2 px-4 py-2 text-sm font-primary hover:bg-gray-50 transition-colors">
                <User size={16} />
                Account
              </a>
              <div className="border-t border-gray-100 my-1"></div>
              <button className="w-full flex items-center gap-2 px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors text-left" onClick={logout}>
                <LogOut size={16} />
                Logout
              </button>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}
