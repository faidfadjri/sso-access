"use client";

import { usePathname } from "next/navigation";
import { Sidebar, Topbar, BottomNav } from "@/components";
import { ReactNode } from "react";

export default function AppShell({ children }: { children: ReactNode }) {
  const pathname = usePathname();
  
  const isPublicRoute = pathname === "/login" || pathname === "/forgot-password" || pathname === "/forgot-username" || pathname === "/reset-password";

  if (isPublicRoute) {
    return <>{children}</>;
  }

  return (
    <div className="flex min-h-screen bg-smoke">
      <Sidebar />
      <BottomNav />

      <div className="flex-1 flex flex-col min-w-0 pb-[70px] md:pb-0">
        <Topbar />

        <main className="flex-1 px-5 lg:p-8">
          {children}
          <div className="block md:hidden">
            <br /><br />
          </div>
        </main>
      </div>
    </div>
  );
}
