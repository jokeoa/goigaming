import type { UserProfile } from "@/types/user";
import { api } from "./client";

export function getMe(): Promise<UserProfile> {
  return api.get<UserProfile>("/users/me");
}
