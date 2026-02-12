import { z } from "zod";

export const userProfileSchema = z.object({
  id: z.string(),
  username: z.string(),
  email: z.string().email(),
  created_at: z.string(),
});

export type UserProfile = z.infer<typeof userProfileSchema>;
