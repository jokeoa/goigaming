import type { TokenResponse } from "@/types/auth";
import type { UserProfile } from "@/types/user";
import { api } from "./client";

export function registerUser(data: {
  readonly username: string;
  readonly email: string;
  readonly password: string;
}): Promise<UserProfile> {
  return api.post<UserProfile>("/auth/register", data);
}

export function loginUser(data: {
  readonly email: string;
  readonly password: string;
}): Promise<TokenResponse> {
  return api.post<TokenResponse>("/auth/login", data);
}
