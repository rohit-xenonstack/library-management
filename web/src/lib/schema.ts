import { z } from "zod";

export const LoginSchema = z.object({
  email: z.string().email(),
});

export const RegisterSchema = z.object({
  name: z.string().max(25),
  email: z.string().email(),
  contact: z.string().min(10).max(10),
});
