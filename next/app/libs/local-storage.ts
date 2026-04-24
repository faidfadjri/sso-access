const isBrowser = typeof window !== "undefined";

export const localStorageLib = {
  get<T>(key: string): T | null {
    if (!isBrowser) return null;

    try {
      const value = localStorage.getItem(key);
      return value ? (JSON.parse(value) as T) : null;
    } catch {
      return null;
    }
  },

  set<T>(key: string, value: T) {
    if (!isBrowser) return;
    localStorage.setItem(key, JSON.stringify(value));
  },

  remove(key: string) {
    if (!isBrowser) return;
    localStorage.removeItem(key);
  },
};