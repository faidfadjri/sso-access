"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { Grid } from "react-feather";
import menuItems from "@/public/data/menu.json";
import { iconMap } from "@/app/libs/menu";
import { authSession } from "@/app/api";
import { ROLES } from "@/app/libs/roles";

const BottomNav = () => {
  const pathname = usePathname();
  const [role, setRole] = useState<string | null>(null);

  useEffect(() => {
    setRole(authSession.getRole() || null);
  }, []);

  return (
    <div className="fixed bottom-0 left-0 right-0 border-t rounded-t-2xl border-gray-200 px-2 md:hidden z-50 shadow-[0_-4px_6px_-1px_rgba(0,0,0,0.05)] bg-secondary">
      <nav className="flex justify-around items-center h-[60px]">
        {menuItems.map((item) => {
          if (item.title === "Clients" && role !== ROLES.SUPER_ADMIN) {
            return null;
          }
          if (item.title === "Users" && role !== ROLES.SUPER_ADMIN && role !== ROLES.ADMIN) {
            return null;
          }
          const Icon = iconMap[item.icon] || Grid;
          const isActive = pathname === item.path;
          
          return (
            <Link
              key={item.path}
              href={item.path}
              className={`flex flex-col items-center justify-center w-16 h-full gap-1 transition-colors relative ${
                isActive ? "text-white" : "text-gray-400 hover:text-gray-300"
              }`}
            >
              {isActive && (
                 <div className="absolute top-0 left-1/2 -translate-x-1/2 w-8 h-[3px] bg-white rounded-b-md" />
              )}
              <Icon size={20} />
              <span className="text-[10px] font-medium leading-none mt-0.5">{item.title}</span>
            </Link>
          );
        })}
      </nav>
    </div>
  );
}

export default BottomNav;