import { clearAuth, getAccessToken } from "./auth";

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
  const hadToken = Boolean(token);
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

    if (response.status === 401 && hadToken) {
      const normalized = String(message || "").toLowerCase();
      const isTokenIssue =
        normalized.includes("token") ||
        normalized.includes("expired") ||
        normalized.includes("credential");
      if (isTokenIssue) {
        clearAuth();
        if (typeof window !== "undefined") {
          window.setTimeout(() => {
            window.location.replace("/login");
          }, 0);
        }
      }
    }

    throw new Error(message);
  }

  return payload;
}

const ensureArray = (value) => (Array.isArray(value) ? value : []);
const unwrap = (payload) => payload?.data ?? payload ?? null;

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

export const doctorProfileApi = {
  getMine: async () => {
    const response = await request("/doctor/profile");
    return response?.data ?? null;
  },
  update: async (payload) => {
    const response = await request("/doctor/profile", {
      method: "PUT",
      body: JSON.stringify(payload),
    });
    return response?.data ?? null;
  },
};

export const patientApi = {
  searchDoctors: async (query = "") => {
    const params = new URLSearchParams();
    if (query.trim()) {
      params.set("q", query.trim());
    }
    const response = await request(
      params.size
        ? `/patient/doctors/search?${params.toString()}`
        : "/patient/doctors/search"
    );
    return ensureArray(unwrap(response));
  },
  getAppointments: async () => {
    const response = await request("/patient/appointments");
    return ensureArray(unwrap(response));
  },
  createAppointment: async (payload) => {
    return request("/patient/appointments", {
      method: "POST",
      body: JSON.stringify(payload),
    });
  },
  cancelAppointment: async (id) => {
    await request(`/patient/appointments/${id}/cancel`, { method: "PATCH" });
  },
};

export const doctorAppointmentApi = {
  getMine: async () => {
    const response = await request("/doctor/appointments");
    return ensureArray(unwrap(response));
  },
  updateStatus: async (id, status) => {
    await request(`/doctor/appointments/${id}`, {
      method: "PATCH",
      body: JSON.stringify({ status }),
    });
  },
};
