import { z } from "zod";

export const apiResponseSchema = <T extends z.ZodType>(dataSchema: T) =>
  z.object({
    success: z.boolean(),
    data: dataSchema.optional(),
    error: z.string().optional(),
  });

export type ApiResponse<T> = {
  readonly success: boolean;
  readonly data?: T;
  readonly error?: string;
};

export type PaginationParams = {
  readonly limit?: number;
  readonly offset?: number;
};
