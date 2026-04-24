"use client";

import { useState, useEffect } from "react";
import { Grid, Menu } from "react-feather";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useSidebar } from "./SidebarContext";
import Image from "next/image"
import styles from "./sidebar.module.css";
import menuItems from "@/public/data/menu.json";
import { iconMap } from "@/app/libs/menu";
import { authSession } from "@/app/api";
import { ROLES } from "@/app/libs/roles";

export default function Sidebar() {
  const pathname = usePathname();
  const { isCollapsed, toggleSidebar } = useSidebar();
  const [role, setRole] = useState<string | null>(null);

  useEffect(() => {
    setRole(authSession.getRole() || null);
  }, []);

  return (
    <>
      <aside
        className={`${styles.sidebar} ${isCollapsed ? styles.sidebarCollapsed : styles.sidebarExpanded
          } hidden md:flex flex-col`}
      >
        <div
          className={`${styles.header} ${isCollapsed ? styles.headerCollapsed : styles.headerExpanded
            }`}
        >
          <button
            onClick={toggleSidebar}
            className={styles.menuButton}
          >
            <Menu size={20} className={styles.menuIcon} />
          </button>
          {!isCollapsed && (
            <div className={styles.logoContainer}>
              <Image alt="app-logo" src="/icons/logo.svg" width={150} height={100} loading="eager"/>
            </div>
          )}
        </div>

        <nav className={styles.nav}>
          {menuItems.map((item) => {
            if (item.title === "Clients" && role !== ROLES.SUPER_ADMIN) {
              return null;
            }
            if (item.title === "Users" && role !== ROLES.SUPER_ADMIN && role !== ROLES.ADMIN) {
              return null;
            }
            const Icon = iconMap[item.icon] || Grid;
            return (
              <Link
                key={item.path}
                href={item.path}
                className={`${styles.navLink} ${pathname === item.path ? styles.activeLink : styles.inactiveLink
                  } ${isCollapsed ? styles.navLinkCollapsed : ""}`}
                title={isCollapsed ? item.title : ""}
              >
                <Icon size={20} className={styles.icon} />
                {!isCollapsed && <span>{item.title}</span>}
              </Link>
            );
          })}
        </nav>
      </aside>
    </>
  );
}
