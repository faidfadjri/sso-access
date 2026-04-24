export interface JWT {
  token: string;
  expiresAt?: string;
}

export const setToken = (token: string) => {
  localStorage.setItem('token', token);
};

export const getToken = (): string | null => {
  return localStorage.getItem('token');
};

export const removeToken = () => {
  localStorage.removeItem('token');
};

export const parseJwt = (token: string): any | null => {
  try {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split('')
        .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
        .join('')
    );
    return JSON.parse(jsonPayload);
  } catch (err) {
    return null;
  }
};

export const getUsername = (): string | null => {
  const token = getToken();
  if (!token) return null;
  const payload = parseJwt(token);
  return payload?.username || null;
};

export const getRole = (): string | null => {
  const token = getToken();
  if (!token) return null;
  const payload = parseJwt(token);
  return payload?.role || null;
};

export const isTokenValid = (): boolean => {
  const token = getToken();
  if (!token) return false;
  const payload = parseJwt(token);
  if (!payload || !payload.exp) return false;
  const now = Math.floor(Date.now() / 1000);
  return payload.exp > now;
};

export const logout = () => {
  removeToken();
};

export async function generateCodeChallenge(verifier: string) {
  const data = new TextEncoder().encode(verifier)
  const digest = await crypto.subtle.digest("SHA-256", data)

  return btoa(String.fromCharCode(...new Uint8Array(digest)))
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=+$/, "")
}


export const generateRandomState = () => {
    return Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);
} 