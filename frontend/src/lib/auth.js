const STORAGE_KEY = "auth";

function isBrowser() {
  return typeof window !== "undefined" && typeof localStorage !== "undefined";
}

export function getStoredAuth() {
  if (!isBrowser()) return null;
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    return raw ? JSON.parse(raw) : null;
  } catch (err) {
    console.error("Failed to parse auth data", err);
    return null;
  }
}

export function getStoredUser() {
  const auth = getStoredAuth();
  return auth?.user ?? null;
}

export function getAccessToken() {
  const auth = getStoredAuth();
  return auth?.token ?? null;
}

export function saveAuth(authPayload) {
  if (!isBrowser()) return;
  if (!authPayload) {
    clearAuth();
    return;
  }
  localStorage.setItem(STORAGE_KEY, JSON.stringify(authPayload));
}

export function clearAuth() {
  if (!isBrowser()) return;
  localStorage.removeItem(STORAGE_KEY);
}
