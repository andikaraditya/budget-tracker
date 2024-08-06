import { cookies } from "next/headers";

export type apiResponse = {
  status: number;
  data: any;
  error: any;
};

export function getToken() {
  return cookies().get("token")?.value;
}
