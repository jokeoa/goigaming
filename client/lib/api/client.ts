import { API_BASE_URL } from "@/lib/constants";

type RequestOptions = {
  readonly method?: string;
  readonly body?: unknown;
  readonly headers?: Record<string, string>;
  readonly params?: Record<string, string | number | undefined>;
};

class ApiError extends Error {
  constructor(
    message: string,
    public readonly status: number,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

function getToken(): string | null {
  if (typeof window === "undefined") return null;
  try {
    const stored = localStorage.getItem("goi-auth-storage");
    if (!stored) return null;
    const parsed = JSON.parse(stored);
    return parsed?.state?.token ?? null;
  } catch {
    return null;
  }
}

async function request<T>(
  endpoint: string,
  options: RequestOptions = {},
): Promise<T> {
  const { method = "GET", body, headers = {}, params } = options;

  let url = `${API_BASE_URL}${endpoint}`;
  if (params) {
    const searchParams = new URLSearchParams();
    for (const [key, value] of Object.entries(params)) {
      if (value !== undefined) {
        searchParams.set(key, String(value));
      }
    }
    const qs = searchParams.toString();
    if (qs) {
      url = `${url}?${qs}`;
    }
  }

  const token = getToken();
  const requestHeaders: Record<string, string> = {
    "Content-Type": "application/json",
    ...headers,
  };
  if (token) {
    requestHeaders.Authorization = `Bearer ${token}`;
  }

  const response = await fetch(url, {
    method,
    headers: requestHeaders,
    body: body ? JSON.stringify(body) : undefined,
  });

  const json = await response.json();

  if (!response.ok || !json.success) {
    throw new ApiError(
      json.error ?? `Request failed with status ${response.status}`,
      response.status,
    );
  }

  return json.data as T;
}

export const api = {
  get: <T>(
    endpoint: string,
    params?: Record<string, string | number | undefined>,
  ) => request<T>(endpoint, { params }),

  post: <T>(endpoint: string, body?: unknown) =>
    request<T>(endpoint, { method: "POST", body }),

  put: <T>(endpoint: string, body?: unknown) =>
    request<T>(endpoint, { method: "PUT", body }),

  delete: <T>(endpoint: string) => request<T>(endpoint, { method: "DELETE" }),
} as const;

export { ApiError };
