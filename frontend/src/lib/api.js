import { getAccessToken } from "./auth";

const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL || "http://localhost:8080/api";

if (!API_BASE_URL) {
  throw new Error(
    "Missing VITE_API_BASE_URL. Set it in frontend/.env to your Go backend URL."
  );
}

async function request(path, options = {}) {
  const config = {
    headers: {
      "Content-Type": "application/json",
      ...(options.headers || {}),
    },
    credentials: "include",
    ...options,
  };

  const token = getAccessToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE_URL}${path}`, config);
  let payload = null;

  try {
    payload = await response.json();
  } catch (err) {
    payload = null;
  }

  if (!response.ok) {
    const message = payload?.message || "Request failed";
    throw new Error(message);
  }

  return payload;
}

export const authApi = {
  login: (data) =>
    request("/login", {
      method: "POST",
      body: JSON.stringify(data),
    }),
  register: (data) =>
    request("/register", {
      method: "POST",
      body: JSON.stringify(data),
    }),
};

export const doctorApi = {
  getAll: async () => {
    const response = await request("/admin/doctors");
    return response?.data ?? [];
  },
  getById: async (id) => {
    const response = await request(`/admin/doctors/${id}`);
    return response?.data ?? null;
  },
  create: async (payload) => {
    const response = await request("/admin/doctors", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    return response?.data ?? null;
  },
  update: async (id, payload) => {
    const response = await request(`/admin/doctors/${id}`, {
      method: "PUT",
      body: JSON.stringify(payload),
    });
    return response?.data ?? null;
  },
  remove: async (id) => {
    await request(`/admin/doctors/${id}`, { method: "DELETE" });
  },
};

export const doctorScheduleApi = {
  getMine: async () => {
    const response = await request("/doctor/schedules");
    return response?.data ?? [];
  },
  create: async (payload) => {
    const response = await request("/doctor/schedules", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    return response?.data ?? null;
  },
  update: async (id, payload) => {
    const response = await request(`/doctor/schedules/${id}`, {
      method: "PUT",
      body: JSON.stringify(payload),
    });
    return response?.data ?? null;
  },
  remove: async (id) => {
    await request(`/doctor/schedules/${id}`, { method: "DELETE" });
  },
};
