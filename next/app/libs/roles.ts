export const ROLES = {
  SUPER_ADMIN: "Super Admin",
  ADMIN: "Admin",
  USER: "User",
} as const;

export type Role = typeof ROLES[keyof typeof ROLES];
